#!/bin/bash
source ./common.sh
trap "cleanup" INT QUIT TERM SIGHUP SIGINT SIGTERM

## FUNCTIONS

cleanup() {
        pr_error "Interrupt received"
        local pids=$(jobs -pr)
        [ -n "$pids" ] && kill $pids
        exit 1
}

prepare_workdir() {
    mkdir $WORKDIR
    echo $RANDOM_SUFFIX > $RANDOM_SUFFIX_FILE
    rm -rf $WORKDIR_PATH/latest
    ln -s $(pwd)/$WORKDIR $(pwd)/$WORKDIR_PATH/latest
    pr_info "preparing working directory"
}


prepare_cluster_setup() {
    pr_info "compiling the remote setup script"
    if [[ $IIP != '' && $EIP != '' && $RANDOM_SUFFIX != '' && $PULL_SECRET_PATH != '' ]]
    then
        PULL_SECRET="$(base64 -w 0 $PULL_SECRET_PATH)"
        $SED "s#_IIP_#$IIP#" $TEMPLATES/cluster_setup.sh > $WORKDIR/cluster_setup.sh
        $SED -i "s#_EIP_#$EIP#g" $WORKDIR/cluster_setup.sh
        $SED -i "s#_RANDOM_SUFFIX_#$RANDOM_SUFFIX#g" $WORKDIR/cluster_setup.sh
        $SED -i "s#_PULL_SECRET_#$PULL_SECRET#g" $WORKDIR/cluster_setup.sh
        $SED -i "s#_PASS_DEVELOPER_#$PASS_DEVELOPER#g" $WORKDIR/cluster_setup.sh
        $SED -i "s#_PASS_KUBEADMIN_#$PASS_KUBEADMIN#g" $WORKDIR/cluster_setup.sh
        $SED -i "s#_PASS_REDHAT_#$PASS_REDHAT#g" $WORKDIR/cluster_setup.sh
    else
        stop_if_failed 1 "internal IP, external IP, random suffix, pull secret path not set, are you calling ${FUNCNAME[0]} correctly?"
    fi
}

create_ec2_resources(){
    pr_info "creating EC2 resources"
    RESOURCES_NAME="openspot-ng-$RANDOM_SUFFIX"
    GROUPID=`aws ec2 create-security-group --group-name $RESOURCES_NAME --description "openspot-ng security group run timestamp: $RUN_TIMESTAMP" --no-paginate | $JQ -r '.GroupId'`
    stop_if_failed $? "failed to create EC2 security group"
    #KEYPAIR (Created just because mandatory, will be swapped manually fore core user later on)
    $AWS ec2 create-key-pair --key-name $RESOURCES_NAME --no-paginate
    stop_if_failed $? "failed to create EC2 keypair"
    #SG SETUP
    $AWS ec2 authorize-security-group-ingress --group-name $RESOURCES_NAME --protocol tcp --port 22 --cidr 0.0.0.0/0 --no-paginate
    stop_if_failed $? "failed to create SSH rule for security group"
    $AWS ec2 authorize-security-group-ingress --group-name $RESOURCES_NAME --protocol tcp --port 6443 --cidr 0.0.0.0/0 --no-paginate
    stop_if_failed $? "failed to create API rule for security group"
    $AWS ec2 authorize-security-group-ingress --group-name $RESOURCES_NAME --protocol tcp --port 443 --cidr 0.0.0.0/0 --no-paginate
    stop_if_failed $? "failed to create HTTPS rule for security group"
    #CREATE INSTANCE
    $AWS ec2 run-instances --no-paginate --tag-specifications "ResourceType=instance,Tags=[{Key=Name,Value=$RESOURCES_NAME}]" --image-id $AMI_ID --instance-type $INSTANCE_TYPE --security-group-ids $GROUPID --key-name $RESOURCES_NAME > $WORKDIR/$INSTANCE_DESCRIPTION 
    stop_if_failed $? "failed to launch EC2 instance"
}

destroy_ec2_resources() {
  ID=$WORKDIR/$INSTANCE_DESCRIPTION
  [ ! -f  $ID ] && stop_if_failed 1 "Missing EC2 resource descriptor $INSTANCE_DESCRIPTION in $WORKDIR"
  INSTANCE_ID=`$JQ -r '.Instances[0].InstanceId' $ID`
  stop_if_failed $? "Failed to parse InstanceId from $ID"
  KEY_NAME=`$JQ -r '.Instances[0].KeyName' $ID`
  stop_if_failed $? "Failed to parse KeyName from $ID"
  SG_ID=`$JQ -r '.Instances[0].SecurityGroups[0].GroupId' $ID`
  stop_if_failed $? "Failed to parse SecurityGroupName from $ID"
  AZ=`$JQ -r '.Instances[0].Placement.AvailabilityZone' $ID`
  stop_if_failed $? "Failed to parse AvailabilityZone from $ID"

  #check if instance is found 
  $AWS ec2 describe-instance-status --instance-id $INSTANCE_ID  > /dev/null 2>&1
  stop_if_failed $? "instance $INSTANCE_ID not found, is your AWS_PROFILE exported and working?"
  #KILL INSTANCE
  pr_info "terminating instance $INSTANCE_ID"
  $AWS ec2 terminate-instances --instance-ids $INSTANCE_ID
  stop_if_failed $? "failed to terminate instance $INSTANCE_ID"
  #WAIT FOR INSTANCE TO TERMINATE
  while [[ `$AWS ec2 describe-instance-status --instance-id $INSTANCE_ID | $JQ -r ".InstanceStatuses[0].SystemStatus.Status"` != null ]]
  do 
    pr_info "waiting instance $INSTANCE_ID to shutdown, hang on...."
    sleep 1
  done

  pr_info "removing security group $SG_ID"
  TRY=0
  until `$AWS ec2 delete-security-group --group-id $SG_ID`
  do
    pr_info "waiting the security group to be removable try $TRY"
    ((TRY++))
    [[ $TRY == $TEARDOWN_MAX_RETRIES ]] && stop_if_failed 1 "failed to remove $SG_ID, are you sure that this sg exists?!"
    sleep 6
  done
  pr_info "removed security group $SG_ID"
  pr_info "removing keypair $KEY_NAME"
  $AWS --region ${AZ::-1} ec2 delete-key-pair --key-name $KEY_NAME
  stop_if_failed $? "failed to remove keypair $KEY_NAME"
  pr_end "everything has been cleaned up!"
  exit 0

}

swap_ssh_key() {
    pr_info "changing default private key permissions to 400"
    chmod 400 $PRIVATE_KEY
    stop_if_failed $? "unable to change defualt key permissions"
    pr_info "swapping default key with the one just created"
    $SSHKEYGEN -f $WORKDIR/id_rsa -q -N ''
    stop_if_failed $? "failed to generate the key pair"
    $SCP -o StrictHostKeychecking=no -P $SSH_PORT -i $PRIVATE_KEY $WORKDIR/id_rsa.pub core@$EIP:.
    stop_if_failed $? "failed to upload the public key on the machine @ $EIP"
    $SSH -o StrictHostKeychecking=no -p $SSH_PORT -i $PRIVATE_KEY core@$EIP "cat /home/core/id_rsa.pub > /home/core/.ssh/authorized_keys"
    stop_if_failed $? "failed to swap the key on the authorized_keys  @ $EIP"
    #after swapping on VM private key is replaced by the new one
    PRIVATE_KEY=$WORKDIR/id_rsa
}

inject_and_run_cluster_setup() {
    pr_info "injecting the setup script on the machine & running it [next logs will be remote]"
    $SCP -o StrictHostKeychecking=no -P $SSH_PORT -i $PRIVATE_KEY $WORKDIR/cluster_setup.sh core@$EIP:/var/home/core/
    $SSH -o StrictHostKeychecking=no -p $SSH_PORT -i $PRIVATE_KEY core@$EIP "chmod +x /var/home/core/cluster_setup.sh"
    $SSH -o StrictHostKeychecking=no -p $SSH_PORT -i $PRIVATE_KEY core@$EIP "sudo /var/home/core/cluster_setup.sh"
}

tail_cluster_setup() {
    pr_info "reading from VM logs from remote instance"
    pr_info "waiting the log to be created, hang on...."
    $SSH -o StrictHostKeychecking=no -p $SSH_PORT -i $PRIVATE_KEY core@$EIP "while [ ! -f /tmp/$RANDOM_SUFFIX.log ]; do sleep 1; done"
    pr_info "log detected, printing remote output on $EIP:"
    PREVIOUS_LINE=""
    while :
    do
        LINE=$($SSH -o StrictHostKeychecking=no -p $SSH_PORT -i $PRIVATE_KEY core@$EIP "sudo tail -n 1 /tmp/$RANDOM_SUFFIX.log")
        stop_if_failed $? "impossible to get the logs from $EIP"
        #xargs used to remove leading space
        if [[ $LINE =~ "[ERR]" ]]
        then
            CLEANLINE=${LINE//"[ERR]"}
            stop_if_failed 1 "$EIP -> $(echo $CLEANLINE | xargs)"
        elif [[ $LINE =~ "[END]" ]]
        then
            CLEANLINE=${LINE//"[END]"}
            pr_end "$EIP -> cluster web console: https://$(echo $CLEANLINE | xargs)"
            break
        else
            CLEANLINE=${LINE//"[INF]"}
            if [[ "$CLEANLINE" != "$PREVIOUS_LINE" ]]
            then
                pr_info "$EIP -> $(echo $CLEANLINE | xargs)"
                PREVIOUS_LINE=$CLEANLINE
            fi
        fi
    done
}

get_remote_log() {
    pr_info "getting remote setup log"
    $SCP -o StrictHostKeychecking=no -P $SSH_PORT -i $PRIVATE_KEY core@$EIP:/tmp/$RANDOM_SUFFIX.log $WORKDIR/remote.log
    stop_if_failed $? "impossible to get the logs from $EIP"
}

set_cluster_infos() {
    pr_info "writing cluster informations in $CLUSTER_INFOS"
    jq '.api.address="https://api.'$1'.nip.io"' $CLUSTER_INFOS_TEMPLATE > $CLUSTER_INFOS
    cat <<< $($JQ '.api.port="6443"' $CLUSTER_INFOS) > $CLUSTER_INFOS
    cat <<< $($JQ '.console.address="https://console-openshift-console.apps.'$1'.nip.io"' $CLUSTER_INFOS) > $CLUSTER_INFOS
    cat <<< $($JQ '.console.port="443"' $CLUSTER_INFOS) > $CLUSTER_INFOS
}

create () {
    SECONDS=0
    prepare_workdir
    create_ec2_resources

    INSTANCE_ID=`get_instance_id $WORKDIR/$INSTANCE_DESCRIPTION`
    IIP=`get_instance_private_ip $WORKDIR/$INSTANCE_DESCRIPTION`
    EIP=`get_instance_public_ip $INSTANCE_ID`

    wait_instance_readiness $EIP
    swap_ssh_key
    prepare_cluster_setup
    inject_and_run_cluster_setup > /dev/null 2>&1 &
    tail_cluster_setup
    get_remote_log
    set_cluster_infos $EIP
    duration=$SECONDS
    pr_end "OpenShift cluster baked in $(($duration / 60)) minutes and $(($duration % 60)) seconds"
}


teardown() {
    WORKDIR="$WORKDIR_PATH/$TEARDOWN_RUN"
    [ ! -d $WORKDIR ] && stop_if_failed 1 "$WORKDIR not found, please provide a correct path"
    destroy_ec2_resources
}





usage() {
    echo ""
    echo "*********** OpenSpot NG ***********"
    
    usage="
Cluster Creation :

$(basename "$0") -C -p pull secret path [-d developer user password] [-k kubeadmin user password] [-r redhat user password] [-a AMI ID] [-t Instance type]
where:
    -C  Cluster Creation mode
    -p  pull secret file path (download from https://console.redhat.com/openshift/create/local) 
    -d  developer user password (optional, default: $PASS_DEVELOPER)
    -k  kubeadmin user password (optional, default: $PASS_KUBEADMIN)
    -r  redhat    user password (optional, default: $PASS_REDHAT)
    -a  AMI ID (Amazon Machine Image) from which the VM will be Instantiated (optional, default: $AMI_ID)
    -i  EC2 Instance Type (optional, default; $INSTANCE_TYPE)
    -h  show this help text

Cluster Teardown:

$(basename "$0") -T [-v run id]
    -T  Cluster Teardown mode
    -v  The Id of the run that is gonna be destroyed, corresponds with the numeric name of the folders created in workdir (optional, default: latest)
    -h  show this help text 
    "
    echo "$usage"
    exit 1
}


##COMMANDS CHECK + OS CHECK
[[ `uname` != "Linux" ]] && stop_if_failed 1 "sorry, but $(basename "$0") can be run only on Linux (for the moment), please run the containerized version"
CURL=`which curl 2>/dev/null`
[[ $? != 0 ]] && stop_if_failed 1 "[DEPENDENCY MISSING]: curl, please install it and try again"
JQ=`which jq 2>/dev/null`
[[ $? != 0 ]] && stop_if_failed 1 "[DEPENDENCY MISSING]: jq, please install it and try again"
MD5SUM=`which md5sum 2>/dev/null`
[[ $? != 0 ]] && stop_if_failed 1 "[DEPENDENCY MISSING]: md5sum, please install it and try again"
AWS=`which aws 2>/dev/null`
[[ $? != 0 ]] && stop_if_failed 1 "[DEPENDENCY MISSING]: aws cli, please install it and try again"
HEAD=`which head 2>/dev/null`
[[ $? != 0 ]] && stop_if_failed 1 "[DEPENDENCY MISSING]: head, please install it and try again"
SSHKEYGEN=`which ssh-keygen 2>/dev/null`
[[ $? != 0 ]] && stop_if_failed 1 "[DEPENDENCY MISSING]: ssh-keygen, please install it and try again"
SED=`which sed 2>/dev/null`
[[ $? != 0 ]] && stop_if_failed 1 "[DEPENDENCY MISSING]: GNU sed, please install it and try again"
NC=`which nc 2>/dev/null`
[[ $? != 0 ]] && stop_if_failed 1 "[DEPENDENCY MISSING]: nc (netcat), please install it and try again"
SSH=`which ssh 2>/dev/null`
[[ $? != 0 ]] && stop_if_failed 1 "[DEPENDENCY MISSING]: ssh, please install it and try again"
SCP=`which scp 2>/dev/null`
[[ $? != 0 ]] && stop_if_failed 1 "[DEPENDENCY MISSING]: scp, please install it and try again"





##DEFAULT VALUES THAT CAN BE OVERRIDDEN BY ENV (podman/docker)
[ -z $PASS_DEVELOPER ] && PASS_DEVELOPER="developer"
[ -z $KUBEADMIN ] && PASS_KUBEADMIN="kubeadmin"
[ -z $PASS_REDHAT ] && PASS_REDHAT="redhat"
[ -z $AMI_ID ] && AMI_ID="ami-0569ce8a44f2351be"
[ -z $INSTANCE_TYPE ] && INSTANCE_TYPE="c6in.2xlarge"
[ -z $WORKDIR_PATH ] && WORKDIR_PATH="workdir"
[ -z $WORKING_MODE ] && WORKING_MODE=""
[ -z $TEARDOWN_RUN ] && TEARDOWN_RUN="latest"

##CONST
SSH_PORT="22"
RUN_TIMESTAMP=`date +%s`
INSTANCE_DESCRIPTION="instance_description.json"
RANDOM_SUFFIX=`echo $RANDOM | $MD5SUM | $HEAD -c 8`
WORKDIR="$WORKDIR_PATH/$RUN_TIMESTAMP"
TEMPLATES="templates"
LOG_FILE="$WORKDIR/local.log"
TEARDOWN_LOGFILE="$WORKDIR_PATH/teardown_$RUN_TIMESTAMP.log"
TEARDOWN_MAX_RETRIES=500
RANDOM_SUFFIX_FILE="$WORKDIR/suffix"
PRIVATE_KEY="id_ecdsa_crc"
CLUSTER_INFOS_FILE="cluster_infos.json"
CLUSTER_INFOS_TEMPLATE="$TEMPLATES/$CLUSTER_INFOS_FILE"
CLUSTER_INFOS="$WORKDIR/$CLUSTER_INFOS_FILE"

##ARGS
options=':h:CTp:d:k:r:a:t:v:'
while getopts $options option; do
  case "$option" in
    h) usage;;
    C) WORKING_MODE="C";;
    T) WORKING_MODE="T";;
    p) PULL_SECRET_PATH=$OPTARG;;
    d) PASS_DEVELOPER=$OPTARG;;
    k) PASS_KUBEADMIN=$OPTARG;;
    r) PASS_REDHAT=$OPTARG;;
    a) AMI_ID=$OPTARG;;
    t) INSTANCE_TYPE=$OPTARG;;
    v) TEARDOWN_RUN=$OPTARG;;
    :) printf "missing argument for -%s\n" "$OPTARG" >&2; usage;;
   \?) printf "illegal option: -%s\n" "$OPTARG" >&2; usage;;
  esac
done

##VARIABLE SANITY CHECKS
#WORKING MODE CHECK
[[ (-z $WORKING_MODE ) ]] && echo -e "\nERROR: Working mode must be set\n" && usage
[[ ( $WORKING_MODE != "C" ) && ( $WORKING_MODE != "T" )  ]] && echo -e "\nERROR: Working mode Must be either -C (creation) or -T (teardown), not $WORKING_MODE\n" && usage
#CHECK MANDATORY ARGS FOR CREATION
[[ ($WORKING_MODE == "C" ) && ( ! "$PULL_SECRET_PATH" ) ]] && echo -e "\nERROR: in creation mode argument -p <pull_secret_path> must be provided\n" && usage 
#CHECK PULL SECRET PATH
[[ ($WORKING_MODE == "C" ) && ( ! -f $PULL_SECRET_PATH ) ]] && echo -e "\nERROR: $PULL_SECRET_PATH pull secret file not found" && usage


##ENTRYPOINT: if everything is ok, run the script.
if [[ $WORKING_MODE == "C" ]]
then
    create
elif [[ $WORKING_MODE == "T" ]]
then
    teardown
else
    usage
fi


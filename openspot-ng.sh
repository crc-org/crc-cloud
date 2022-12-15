#!/bin/bash
source ./common.sh
trap "cleanup" INT QUIT TERM EXIT 

## FUNCTIONS

cleanup() {
        local pids=$(jobs -pr)
        [ -n "$pids" ] && kill $pids
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
        $SED "s#_IIP_#$IIP#" templates/cluster_setup.sh > $WORKDIR/cluster_setup.sh
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

swap_ssh_key() {
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
            if [[ $CLEANLINE != $PREVIOUS_LINE ]]
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
    duration=$SECONDS
    pr_end "CRC cluster baked in $(($duration / 60)) minutes and $(($duration % 60)) seconds"
}



usage() {
    echo ""
    echo "*********** OpenSpot NG ***********"
    
    usage="
Cluster Creation :

$(basename "$0") -C -p pull secret path [-d developer user password] [-k kubeadmin user password] [-r redhat user password] [-a AMI ID] [-t Instance type]
where:
    -C  Cluster Creation mode
    -p  CRC pull secret file path (download from https://console.redhat.com/openshift/create/local) 
    -d  CRC developer user password (optional, default: $PASS_DEVELOPER)
    -k  CRC kubeadmin user password (optional, default: $PASS_KUBEADMIN)
    -r  CRC redhat    user password (optional, default: $PASS_REDHAT)
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


##COMMANDS
CURL=`which curl`
JQ=`which jq`
MD5SUM=`which md5sum`
AWS=`which aws`
HEAD=`which head`
SSHKEYGEN=`which ssh-keygen`
SED=`which sed`
NC=`which nc`
SSH=`which ssh`
SCP=`which scp`
PRIVATE_KEY="id_ecdsa_crc"


##CONST
SSH_PORT="22"
RUN_TIMESTAMP=`date +%s`
INSTANCE_DESCRIPTION="instance_description.json"
RANDOM_SUFFIX=`echo $RANDOM | $MD5SUM | $HEAD -c 8`
RANDOM_SUFFIX_FILE="$WORKDIR/suffix"
LOG_FILE="$WORKDIR/local.log"

##DEFAULT VALUES THAT CAN BE OVERRIDDEN BY ENV (podman/docker)
[ -z $PASS_DEVELOPER ] && PASS_DEVELOPER="developer"
[ -z $KUBEADMIN ] && PASS_KUBEADMIN="kubeadmin"
[ -z $PASS_REDHAT ] && PASS_REDHAT="redhat"
[ -z $AMI_ID ] && AMI_ID="ami-0569ce8a44f2351be"
[ -z $INSTANCE_TYPE ] && INSTANCE_TYPE="c6i.4xlarge"
[ -z $WORKDIR_PATH ] && WORKDIR_PATH="workdir"
[ -z $WORKING_MODE ] && WORKING_MODE=""

WORKDIR="$WORKDIR_PATH/$RUN_TIMESTAMP"

##ARGS
options=':h:CTp:d:k:r:a:t:'
while getopts $options option; do
echo "$option"
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
    :) printf "missing argument for -%s\n" "$OPTARG" >&2; usage;;
   \?) printf "illegal option: -%s\n" "$OPTARG" >&2; usage;;
  esac
done

##CHECKS
#WORKING MODE CHECK
[[ (-z $WORKING_MODE ) ]] && echo -e "\nERROR: Working mode must be set\n" && usage
[[ ( $WORKING_MODE != "C" ) && ( $WORKING_MODE != "T" )  ]] && echo -e "\nERROR: Working mode Must be either -C (creation) or -T (teardown), not $WORKING_MODE\n" && usage
#CHECK MANDATORY ARGS FOR CREATION
[[ ($WORKING_MODE == "C" ) && ( ! "$PULL_SECRET_PATH" ) ]] && echo -e "\nERROR: in creation mode argument -p <pull_secret_path> must be provided\n" && usage 
#CHECK PULL SECRET PATH
[[ ! -f $PULL_SECRET_PATH ]] && echo -e "\nERROR: $PULL_SECRET_PATH pull secret file not found" && usage


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
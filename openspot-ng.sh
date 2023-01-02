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
    #link only if not running in container (teardown will need mandatory run id)
    [[ ! $CONTAINER ]] && ln -s $(readlink -f $WORKDIR) $(readlink -f $WORKDIR_PATH)/latest
    pr_info "preparing working directory"
}


prepare_cluster_setup() {
    pr_info "compiling the remote setup script"
    if [[ $CONTAINER ]]
    then
        [[ -z $PULL_SECRET ]] && stop_if_failed 1 "PULL_SECRET environment variable not set"
    else
        [[ -z $PULL_SECRET_PATH ]] && stop_if_failed 1 "PULL_SECRET_PATH not set"
        [[ ! -f $PULL_SECRET_PATH ]] && stop_if_failed 1 "$PULL_SECRET_PATH file not found"
        PULL_SECRET="$($BASE64 -w 0 $PULL_SECRET_PATH)"
    fi

    
    if [[ $IIP != '' && $EIP != '' && $RANDOM_SUFFIX != '' ]]
    then
        $SED "s#_IIP_#$IIP#" $TEMPLATES/cluster_setup.sh > $WORKDIR/cluster_setup.sh
        $SED -i "s#_EIP_#$EIP#g" $WORKDIR/cluster_setup.sh
        $SED -i "s#_RANDOM_SUFFIX_#$RANDOM_SUFFIX#g" $WORKDIR/cluster_setup.sh
        #remove linebreaks (eg. running from podman -e PULL_SECRET="$(base64)" if used with quotes line breaks are introduced)
        $SED -i "s#_PULL_SECRET_#${PULL_SECRET//$'\n'/}#g" $WORKDIR/cluster_setup.sh
        $SED -i "s#_PASS_DEVELOPER_#$PASS_DEVELOPER#g" $WORKDIR/cluster_setup.sh
        $SED -i "s#_PASS_KUBEADMIN_#$PASS_KUBEADMIN#g" $WORKDIR/cluster_setup.sh
        $SED -i "s#_PASS_REDHAT_#$PASS_REDHAT#g" $WORKDIR/cluster_setup.sh
    else
        stop_if_failed 1 "internal IP, external IP, random suffix  are you calling ${FUNCNAME[0]} correctly?"
    fi
}

create_ec2_resources() {
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
    $AWS ec2 run-instances --no-paginate --tag-specifications "ResourceType=instance,Tags=[{Key=Name,Value=$RESOURCES_NAME}]" \
    --image-id $AMI_ID --instance-type $INSTANCE_TYPE --security-group-ids $GROUPID --key-name $RESOURCES_NAME > $WORKDIR/$INSTANCE_DESCRIPTION

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

create_gcp_resources() {
    pr_info "creating GCP resources"
    RESOURCES_NAME="openspot-ng-$RANDOM_SUFFIX"
    # Check if the Instance_type available
    $GCP compute machine-types describe $INSTANCE_TYPE --quiet
    stop_if_failed $? "failed to get $INSTANCE_TYPE on gcp"
    # Create separate network and subnet
    $GCP compute networks create "${RESOURCES_NAME}" --subnet-mode=custom --bgp-routing-mode=regional --quiet
    stop_if_failed $? "failed to create GCP compute network"
    $GCP compute networks subnets create "${RESOURCES_NAME}" --network "${RESOURCES_NAME}" --range=10.0.0.0/9 --quiet
    stop_if_failed $? "failed to create GCP compute subnet"
    # Create firewall rules and add to network
    $GCP compute firewall-rules create "${RESOURCES_NAME}" --network "${RESOURCES_NAME}" --allow tcp:22,tcp:6443,tcp:443,tcp:80 --quiet
    stop_if_failed $? "failed to create GCP network resource"
    #CREATE INSTANCE
    $GCP compute instances create "${RESOURCES_NAME}" --machine-type $INSTANCE_TYPE --image=crc  --tags $RESOURCES_NAME \
    --subnet "${RESOURCES_NAME}" --network "${RESOURCES_NAME}" --format=json > $WORKDIR/$INSTANCE_DESCRIPTION
    stop_if_failed $? "failed to launch GCP instance"
}

destroy_gcp_resources() {
    pr_info "deleting GCP resources"
    ID=$WORKDIR/$INSTANCE_DESCRIPTION
    [ ! -f  $ID ] && stop_if_failed 1 "Missing GCP resource descriptor $INSTANCE_DESCRIPTION in $WORKDIR"
    INSTANCE_ID=`get_instance_id_gcp $WORKDIR/$INSTANCE_DESCRIPTION`
    $GCP compute instances delete "${INSTANCE_ID}" --quiet
    stop_if_failed $? "failed to delete GCP instance"
    $GCP compute firewall-rules delete "${INSTANCE_ID}" --quiet
    stop_if_failed $? "failed to delete GCP firewall-rules"
    $GCP compute networks subnets delete "${INSTANCE_ID}" --quiet
    stop_if_failed $? "failed to delete GCP subnet"
    $GCP  compute networks delete "${INSTANCE_ID}" --quiet
    stop_if_failed $? "failed to delete GCP network"
    pr_end "everything has been cleaned up!"
    exit 0
}

create_openstack_resources() {
    pr_info "creating Openstack resources"
    RESOURCES_NAME="openspot-ng-$RANDOM_SUFFIX"
    # Check if the Instance_type available
    $OSP flavor show $INSTANCE_TYPE -f json > /dev/null
    stop_if_failed $? "failed to get $INSTANCE_TYPE on openstack"
    # Create separate network and subnet
    $OSP security group create --description "CRC security group" "${RESOURCES_NAME}" -f json > /dev/null
    stop_if_failed $? "failed to create openstack security group"
    $OSP security group rule create --proto tcp --dst-port 22 "${RESOURCES_NAME}" -f json > /dev/null
    stop_if_failed $? "failed to add port 22 to ${RESOURCES_NAME} group"
    $OSP security group rule create --proto tcp --dst-port 80 "${RESOURCES_NAME}" -f json > /dev/null
    stop_if_failed $? "failed to add port 80 to ${RESOURCES_NAME} group"
    $OSP security group rule create --proto tcp --dst-port 443 "${RESOURCES_NAME}" -f json > /dev/null
    stop_if_failed $? "failed to add port 443 to ${RESOURCES_NAME} group"
    $OSP security group rule create --proto tcp --dst-port 6443 "${RESOURCES_NAME}" -f json > /dev/null
    stop_if_failed $? "failed to add port 6443 to ${RESOURCES_NAME} group"
    #CREATE INSTANCE
    $OSP server create --flavor $INSTANCE_TYPE --image crc --nic net-id=provider_net_cci_5 --security-group $RESOURCES_NAME \
        $RESOURCES_NAME -f json --wait > $WORKDIR/$INSTANCE_DESCRIPTION
    stop_if_failed $? "failed to launch openstack instance"
}

destroy_openstack_resources() {
    pr_info "deleting openstack resources"
    ID=$WORKDIR/$INSTANCE_DESCRIPTION
    [ ! -f  $ID ] && stop_if_failed 1 "Missing openstack resource descriptor $INSTANCE_DESCRIPTION in $WORKDIR"
    INSTANCE_ID=`get_instance_id_osp $WORKDIR/$INSTANCE_DESCRIPTION`
    $OSP server delete "${INSTANCE_ID}" --wait --quiet
    stop_if_failed $? "failed to delete openstack instance"
    $OSP security group delete "${INSTANCE_ID}" --quiet
    stop_if_failed $? "failed to delete openstack security group"
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
    $SCP -q -o StrictHostKeychecking=no -o UserKnownHostsFile=/dev/null -P $SSH_PORT -i $PRIVATE_KEY $WORKDIR/id_rsa.pub core@$EIP:.
    stop_if_failed $? "failed to upload the public key on the machine @ $EIP"
    $SSH -q -o StrictHostKeychecking=no -o UserKnownHostsFile=/dev/null -p $SSH_PORT -i $PRIVATE_KEY core@$EIP "cat /home/core/id_rsa.pub > /home/core/.ssh/authorized_keys"
    stop_if_failed $? "failed to swap the key on the authorized_keys  @ $EIP"
    #after swapping on VM private key is replaced by the new one
    GEN_PRIVATE_KEY=$WORKDIR/id_rsa
}

inject_and_run_cluster_setup() {
    pr_info "injecting the setup script on the machine & running it [next logs will be remote]"
    $SCP -o StrictHostKeychecking=no -o UserKnownHostsFile=/dev/null -P $SSH_PORT -i $PRIVATE_KEY -i $GEN_PRIVATE_KEY $WORKDIR/cluster_setup.sh core@$EIP:/var/home/core/
    $SSH -q -o StrictHostKeychecking=no -o UserKnownHostsFile=/dev/null -p $SSH_PORT -i $PRIVATE_KEY -i $GEN_PRIVATE_KEY core@$EIP "chmod +x /var/home/core/cluster_setup.sh"
    $SSH -q -o StrictHostKeychecking=no -o UserKnownHostsFile=/dev/null -p $SSH_PORT -i $PRIVATE_KEY -i $GEN_PRIVATE_KEY core@$EIP "sudo /var/home/core/cluster_setup.sh"
}

tail_cluster_setup() {
    pr_info "reading from VM logs from remote instance"
    pr_info "waiting the log to be created, hang on...."
    $SSH -q -o StrictHostKeychecking=no -o UserKnownHostsFile=/dev/null -p $SSH_PORT -i $PRIVATE_KEY -i $GEN_PRIVATE_KEY core@$EIP "while [ ! -f /tmp/$RANDOM_SUFFIX.log ]; do sleep 1; done"
    pr_info "log detected, printing remote output on $EIP:"
    PREVIOUS_LINE=""
    while :
    do
        LINE=$($SSH -q -o StrictHostKeychecking=no -o UserKnownHostsFile=/dev/null -p $SSH_PORT -i $PRIVATE_KEY -i $GEN_PRIVATE_KEY core@$EIP "sudo tail -n 1 /tmp/$RANDOM_SUFFIX.log")
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
    $SCP -o StrictHostKeychecking=no -o UserKnownHostsFile=/dev/null -P $SSH_PORT -i $PRIVATE_KEY -i $GEN_PRIVATE_KEY core@$EIP:/tmp/$RANDOM_SUFFIX.log $WORKDIR/remote.log
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
    case $CLOUD_PROVIDER in
        "aws")
            create_ec2_resources
            INSTANCE_ID=`get_instance_id_aws $WORKDIR/$INSTANCE_DESCRIPTION`
            IIP=`get_instance_private_ip_aws $WORKDIR/$INSTANCE_DESCRIPTION`
            EIP=`get_instance_public_ip_aws $INSTANCE_ID`
        ;;
        "gcp")
            create_gcp_resources
            INSTANCE_ID=`get_instance_id_gcp $WORKDIR/$INSTANCE_DESCRIPTION`
            IIP=`get_instance_private_ip_gcp $WORKDIR/$INSTANCE_DESCRIPTION`
            EIP=`get_instance_public_ip_gcp $WORKDIR/$INSTANCE_DESCRIPTION`
        ;;
        "openstack")
            create_openstack_resources
            INSTANCE_ID=`get_instance_id_osp $WORKDIR/$INSTANCE_DESCRIPTION`
            IIP=`get_instance_private_ip_osp $WORKDIR/$INSTANCE_DESCRIPTION`
            EIP=`get_instance_public_ip_osp $WORKDIR/$INSTANCE_DESCRIPTION`
        ;;
        *)
            echo "Unknown cloud provider"
            usuage
            exit 1
    esac


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
    WORKDIR="$WORKDIR_PATH/$TEARDOWN_RUN_ID"
    [ ! -d $WORKDIR ] && stop_if_failed 1 "$WORKDIR not found, please provide a correct path"
    case $CLOUD_PROVIDER in
        "aws")
            destroy_ec2_resources
            ;;
        "gcp")
            destroy_gcp_resources
            ;;
        "openstack")
            destroy_openstack_resources
            ;;
         *)
            echo "Unknown cloud provider"
            exit 1
            ;;
     esac
}

set_workdir_dependent_variables() {
    WORKDIR="$WORKDIR_PATH/$RUN_TIMESTAMP"
    LOG_FILE="$WORKDIR/local.log"
    TEARDOWN_LOGFILE="$WORKDIR_PATH/teardown_$RUN_TIMESTAMP.log"
    RANDOM_SUFFIX_FILE="$WORKDIR/suffix"
    CLUSTER_INFOS="$WORKDIR/$CLUSTER_INFOS_FILE"
}

usage() {
    echo ""
    echo "*********** OpenSpot-NG ***********"
    
    usage="
Cluster Creation :

$(basename "$0") -C -p pull secret path [-i cloud provider] [-d developer user password] [-k kubeadmin user password] [-r redhat user password] [-a AMI ID] [-t Instance type]
where:
    -C  Cluster Creation mode
    -i  Cloud/Infra provider (optional, default: $CLOUD_PROVIDER)
    -p  pull secret file path (download from https://console.redhat.com/openshift/create/local) 
    -d  developer user password (optional, default: $PASS_DEVELOPER)
    -k  kubeadmin user password (optional, default: $PASS_KUBEADMIN)
    -r  redhat    user password (optional, default: $PASS_REDHAT)
    -a  Image ID (Cloud provider Machine Image) from which the VM will be Instantiated (optional, default: $AMI_ID)
    -t  Cloud provider Instance Type (optional, default; $INSTANCE_TYPE)
    -h  show this help text

Cluster Teardown:

$(basename "$0") -T [-i cloud provider] [-v run id]
    -T  Cluster Teardown mode
    -i  Cloud/Infra provider (optional, default: $CLOUD_PROVIDER)
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
BASE64=`which base64 2>/dev/null`
[[ $? != 0 ]] && stop_if_failed 1 "[DEPENDENCY MISSING]: base64, please install it and try again"
FIGLET=`which figlet 2>/dev/null`
[[ $CONTAINER && ( $? != 0 ) ]] && stop_if_failed 1 "[DEPENDENCY MISSING]: figlet (container mode only), please install it and try again"


##DEFAULT VALUES THAT CAN BE OVERRIDDEN BY ENV (podman/docker)
[ -z $PASS_DEVELOPER ] && PASS_DEVELOPER="developer"
[ -z $KUBEADMIN ] && PASS_KUBEADMIN="kubeadmin"
[ -z $PASS_REDHAT ] && PASS_REDHAT="redhat"
[ -z $INSTANCE_TYPE ] && INSTANCE_TYPE="c6in.2xlarge"
[ -z $WORKDIR_PATH ] && WORKDIR_PATH="workdir"
[ -z $WORKING_MODE ] && WORKING_MODE=""
[ -z $TEARDOWN_RUN_ID ] && TEARDOWN_RUN_ID="latest"
[ -z $CLOUD_PROVIDER ] && CLOUD_PROVIDER="aws"

##CONST
SSH_PORT="22"
RUN_TIMESTAMP=`date +%s`
INSTANCE_DESCRIPTION="instance_description.json"
RANDOM_SUFFIX=`echo $RANDOM | $MD5SUM | $HEAD -c 8`
PRIVATE_KEY="id_ecdsa_crc"
GEN_PRIVATE_KEY=""
CLUSTER_INFOS_FILE="cluster_infos.json"
TEMPLATES="templates"
TEARDOWN_MAX_RETRIES=500
CLUSTER_INFOS_TEMPLATE="$TEMPLATES/$CLUSTER_INFOS_FILE"
AMI_ID="ami-0569ce8a44f2351be"

##ARGS
#collects args from commandline only if not in container otherwise variables are fed by -e VAR=VALUE 
if [ $CONTAINER ] 
then
    WORKDIR_PATH="/workdir"
    set_workdir_dependent_variables
    #check working mode
    [[ (-z $WORKING_MODE ) ]] && stop_if_failed 1 "WORKING_MODE environment variable must be set"
    [[ ( $WORKING_MODE != "C" ) && ( $WORKING_MODE != "T" )  ]] && \
    stop_if_failed 1 "WORKING_MODE value must be either C (create) or T(teardown) $WORKING_MODE is not a valid value"
    #check pull secret
    [[ ($WORKING_MODE == "C") && ( -z $PULL_SECRET ) ]] && stop_if_failed 1 "PULL_SECRET environment variable must be set and must contain a valid base64 encoded pull_secret, please refer to the README.md"
    #check workdir mount write permissions 
    [[ ! -d $WORKDIR_PATH ]] && stop_if_failed 1 "please mount the workdir filesystem, refer to README.md for further instructions"
    [[ ! -w $WORKDIR_PATH ]] && \
    stop_if_failed 1 "please grant write permissions to the host folder mounted as volume, please refer to README.md for further instructions"
    #check the cloud provider env and set env variable accordingly
    case $CLOUD_PROVIDER in
        "aws")
            #check AWS credentials
            [[ -z $AWS_ACCESS_KEY_ID ]] && stop_if_failed 1 "AWS_ACCESS_KEY_ID must be set, please refer to AWS CLI documentation https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html"
            [[ -z $AWS_SECRET_ACCESS_KEY ]] && stop_if_failed 1 "AWS_ACCESS_KEY_ID must be set, please refer to AWS CLI documentation https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html"
            [[ -z $AWS_DEFAULT_REGION ]] && stop_if_failed 1 "AWS_ACCESS_KEY_ID must be set, please refer to AWS CLI documentation https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html"
            ;;
         "gcp")
            #check GCP credentials
            [[ -z $GOOGLE_APPLICATION_CREDENTIALS ]] && stop_if_failed 1 "GOOGLE_APPLICATION_CREDENTIALS must be set, please refer to GCP CLI documentation https://cloud.google.com/docs/authentication/application-default-credentials"
            [[ -z $CLOUDSDK_COMPUTE_REGION ]] && stop_if_failed 1 "AWS_ACCESS_KEY_ID must be set, please refer to GCP CLI documentation https://cloud.google.com/compute/docs/gcloud-compute"
            ;;
        *)
            echo "Unknown provider"
            exit 1
            ;;
    esac
    [[ ( $WORKING_MODE == "T" ) && ($TEARDOWN_RUN_ID == "latest") ]] && stop_if_failed 1 "TEARDOWN_RUN_ID must be set in container mode. Please set this value with the run id that you want to teardown, refer to README.md for further instructions"
else
    set_workdir_dependent_variables
    options=':h:CTp:d:k:r:a:t:v:i:'
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
        v) TEARDOWN_RUN_ID=$OPTARG;;
        i) CLOUD_PROVIDER=$OPTARG;;
        :) printf "missing argument for -%s\n" "$OPTARG" >&2; usage;;
    \?) printf "illegal option: -%s\n" "$OPTARG" >&2; usage;;
    esac
    done

    ##VARIABLE SANITY CHECKS

    #WORKING MODE CHECK
    case $CLOUD_PROVIDER in
    "aws")
        AWS=`which aws 2>/dev/null`
        [[ $? != 0 ]] && stop_if_failed 1 "[DEPENDENCY MISSING]: aws cli, please install it and try again"
    ;;
    "gcp")
        GCP=`which gcloud 2>/dev/null`
        [[ $? != 0 ]] && stop_if_failed 1 "[DEPENDENCY MISSING]: gcloud cli, please install it and try again"
    ;;
    "openstack")
        OSP=`which openstack 2>/dev/null`
        [[ $? != 0 ]] && stop_if_failed 1 "[DEPENDENCY MISSING]: openstack cli, please install it and try again"
    ;;
    *)
        echo "unknown cloud provider"
        usage
        exit 1
    esac

    [[ (-z $WORKING_MODE ) ]] && echo -e "\nERROR: Working mode must be set\n" && usage
    [[ ( $WORKING_MODE != "C" ) && ( $WORKING_MODE != "T" )  ]] && echo \
    -e "\nERROR: Working mode Must be either -C (creation) or -T (teardown), not $WORKING_MODE\n" && usage
    #CHECK MANDATORY ARGS FOR CREATION
    [[ ($WORKING_MODE == "C" ) && ( ! "$PULL_SECRET_PATH" ) ]] && \
    echo -e "\nERROR: in creation mode argument -p <pull_secret_path> must be provided\n" && usage 
    #CHECK PULL SECRET PATH
    [[ ($WORKING_MODE == "C" ) && ( ! -f $PULL_SECRET_PATH ) ]] && \
    echo -e "\nERROR: $PULL_SECRET_PATH pull secret file not found" && usage
fi





##ENTRYPOINT: if everything is ok, run the script.
[[ $CONTAINER ]] && figlet -f smslant -c "OpenSpot-NG" && echo -e "\n\n"
if [[ $WORKING_MODE == "C" ]]
then
    create
elif [[ $WORKING_MODE == "T" ]]
then
    teardown
else
    usage
fi


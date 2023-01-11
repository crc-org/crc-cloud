#!/bin/bash
trap "cleanup" INT QUIT TERM SIGHUP SIGINT SIGTERM



### FUNCTIONS
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
    # link only if not running in container (teardown will need mandatory run id)
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
        # remove linebreaks (eg. running from podman -e PULL_SECRET="$(base64)" if used with quotes line breaks are introduced)
        $SED -i "s#_PULL_SECRET_#${PULL_SECRET//$'\n'/}#g" $WORKDIR/cluster_setup.sh
        $SED -i "s#_PASS_DEVELOPER_#$PASS_DEVELOPER#g" $WORKDIR/cluster_setup.sh
        $SED -i "s#_PASS_KUBEADMIN_#$PASS_KUBEADMIN#g" $WORKDIR/cluster_setup.sh
        $SED -i "s#_PASS_REDHAT_#$PASS_REDHAT#g" $WORKDIR/cluster_setup.sh
    else
        stop_if_failed 1 "internal IP, external IP, random suffix  are you calling ${FUNCNAME[0]} correctly?"
    fi
}


swap_ssh_key() {
    pr_info "changing default private key permissions to 400"
    chmod 400 $PRIVATE_KEY
    stop_if_failed $? "unable to change defualt key permissions"
    pr_info "swapping default key with the one just created"
    $SSHKEYGEN -f $WORKDIR/id_rsa -q -N ''
    stop_if_failed $? "failed to generate the key pair"
    $SCP -q -o UserKnownHostsFile=/dev/null -o StrictHostKeychecking=no -P $SSH_PORT -i $PRIVATE_KEY $WORKDIR/id_rsa.pub core@$EIP:.
    stop_if_failed $? "failed to upload the public key on the machine @ $EIP"
    $SSH -q -o UserKnownHostsFile=/dev/null -o StrictHostKeychecking=no -p $SSH_PORT -i $PRIVATE_KEY core@$EIP "cat /home/core/id_rsa.pub > /home/core/.ssh/authorized_keys"
    stop_if_failed $? "failed to swap the key on the authorized_keys  @ $EIP"
    # after swapping on VM private key is replaced by the new one
    PRIVATE_KEY=$WORKDIR/id_rsa
}

inject_and_run_cluster_setup() {
    pr_info "injecting the setup script on the machine & running it [next logs will be remote]"
    $SCP -q -o UserKnownHostsFile=/dev/null -o StrictHostKeychecking=no -P $SSH_PORT -i $PRIVATE_KEY $WORKDIR/cluster_setup.sh core@$EIP:/var/home/core/
    $SSH -q -o UserKnownHostsFile=/dev/null -o StrictHostKeychecking=no -p $SSH_PORT -i $PRIVATE_KEY core@$EIP "chmod +x /var/home/core/cluster_setup.sh"
    $SSH -q -o UserKnownHostsFile=/dev/null -o StrictHostKeychecking=no -p $SSH_PORT -i $PRIVATE_KEY core@$EIP "sudo /var/home/core/cluster_setup.sh"
}

tail_cluster_setup() {
    pr_info "reading from VM logs from remote instance"
    pr_info "waiting the log to be created, hang on...."
    $SSH -q -o UserKnownHostsFile=/dev/null -o StrictHostKeychecking=no -p $SSH_PORT -i $PRIVATE_KEY core@$EIP "while [ ! -f /tmp/$RANDOM_SUFFIX.log ]; do sleep 1; done"
    pr_info "log detected, printing remote output on $EIP:"
    PREVIOUS_LINE=""
    while :
    do
        LINE=$($SSH -q -o UserKnownHostsFile=/dev/null -o StrictHostKeychecking=no -p $SSH_PORT -i $PRIVATE_KEY core@$EIP "sudo tail -n 1 /tmp/$RANDOM_SUFFIX.log")
        stop_if_failed $? "impossible to get the logs from $EIP"
        # xargs used to remove leading space
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
    $SCP -q -o UserKnownHostsFile=/dev/null -o StrictHostKeychecking=no -P $SSH_PORT -i $PRIVATE_KEY core@$EIP:/tmp/$RANDOM_SUFFIX.log $WORKDIR/remote.log
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
    deployer_load_dependencies
    deployer_create $@
    IIP=`deployer_get_iip`
    EIP=`deployer_get_eip`

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
    deployer_load_dependencies
    deployer_teardown $@
}

set_workdir_dependent_variables() {
    WORKDIR="$WORKDIR_PATH/$RUN_TIMESTAMP"
    LOG_FILE="$WORKDIR/local.log"
    TEARDOWN_LOGFILE="$WORKDIR_PATH/teardown_$RUN_TIMESTAMP.log"
    RANDOM_SUFFIX_FILE="$WORKDIR/suffix"
    CLUSTER_INFOS="$WORKDIR/$CLUSTER_INFOS_FILE"
}

parse_args () {
    while getopts $BASE_OPTIONS option; do
    case "$option" in
        h) SHOW_HELP=1;;
        C) WORKING_MODE="C";pr_info "working mode: CREATE";;
        D) DEPLOYER_API=$OPTARG;pr_info "deployer api: $OPTARG";;
        T) WORKING_MODE="T";pr_info "working mode: TEARDOWN";;
        p) PULL_SECRET_PATH=$OPTARG;pr_info "pull secret path: $OPTARG";;
        d) PASS_DEVELOPER=$OPTARG;pr_info "developer user pass has been overridden";;
        k) PASS_KUBEADMIN=$OPTARG;pr_info "kubeadmin user pass has been overridden";;
        r) PASS_REDHAT=$OPTARG;pr_info "redhat user pass has been overridden";;
        v) TEARDOWN_RUN_ID=$OPTARG;pr_info "tearing down run id: $OPTARG";;
        :) printf "missing argument for -%s\n" "$OPTARG" >&2; usage;;
    \?);;
    esac
    done
    OPTIND=1
    OPTARG=""
    option=""
}

check_command_dependencies() {
    ## COMMANDS CHECK + OS CHECK
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
    CAT=`which cat 2>/dev/null`
    [[ $? != 0 ]] && stop_if_failed 1 "[DEPENDENCY MISSING]: cat, please install it and try again"
    GREP=`which grep 2>/dev/null`
    [[ $? != 0 ]] && stop_if_failed 1 "[DEPENDENCY MISSING]: grep, please install it and try again"
    EXPR=`which expr 2>/dev/null`
    [[ $? != 0 ]] && stop_if_failed 1 "[DEPENDENCY MISSING]: expr, please install it and try again"
    FIGLET=`which figlet 2>/dev/null`
    [[ $CONTAINER && ( $? != 0 ) ]] && stop_if_failed 1 "[DEPENDENCY MISSING]: figlet (container mode only), please install it and try again"
}

set_const() {
    ## CONST
    SSH_PORT="22"
    RUN_TIMESTAMP=`date +%s`
    INSTANCE_DESCRIPTION="instance_description.json"
    RANDOM_SUFFIX=`echo $RANDOM | $MD5SUM | $HEAD -c 8`
    PRIVATE_KEY="id_ecdsa_crc"
    CLUSTER_INFOS_FILE="cluster_infos.json"
    TEMPLATES="templates"
    TEARDOWN_MAX_RETRIES=500
    CLUSTER_INFOS_TEMPLATE="$TEMPLATES/$CLUSTER_INFOS_FILE"
    DEFAULT_DEPLOYER="bash"
    BASE_OPTIONS=":hCTp:D:d:k:r:v:"
    PLUGIN_FOLDER="plugin"
    PLUGIN_DEPLOYER_FOLDER="$PLUGIN_FOLDER/deployer"
    [[ $CONTAINER ]] && WORKDIR_PATH="/workdir"
}

usage() {
    echo -e "\n*********** CRC-Cloud ***********"
    usage="
Cluster Creation :

$(basename "$0") -C -p pull secret path [-D infrastructure_deployer] [-d developer user password] [-k kubeadmin user password] [-r redhat user password]
where:
    -D  Infrastructure Deployer (default: $DEFAULT_DEPLOYER) *NOTE* Must match with the folder name placed in plugin/deployer (please refer to the deployer documentation in api/README.md)
    -C  Cluster Creation mode
    -p  pull secret file path (download from https://console.redhat.com/openshift/create/local) 
    -d  developer user password (optional, default: $PASS_DEVELOPER)
    -k  kubeadmin user password (optional, default: $PASS_KUBEADMIN)
    -r  redhat    user password (optional, default: $PASS_REDHAT)
    -h  show this help text

Cluster Teardown:

$(basename "$0") -T [-D infrastructure_deployer] [-v run id]
    -D  Infrastructure Deployer (default: $DEFAULT_DEPLOYER) *NOTE* Must match with the folder name placed in plugin/deployer (please refer to the deployer documentation in api/README.md)
    -T  Cluster Teardown mode
    -v  The Id of the run that is gonna be destroyed, corresponds with the numeric name of the folders created in workdir (optional, default: latest)
    -h  show this help text 
    "
    echo "$usage"

    echo -e "\n*********** Deployer: $DEPLOYER_API ***********"
    deployer_usage
    exit 0
}
### END FUNCTIONS

### INIT
[[ $CONTAINER ]] && figlet -f slant -c "CRC-Cloud `cat VERSION`" && echo -e "\n\n"
check_command_dependencies
set_const

## IMPORT API COMMON
source ./api/common.sh

## DEFAULT VALUES THAT CAN BE OVERRIDDEN BY ENV (podman/docker)
[ -z $PASS_DEVELOPER ] && PASS_DEVELOPER="developer"
[ -z $KUBEADMIN ] && PASS_KUBEADMIN="kubeadmin"
[ -z $PASS_REDHAT ] && PASS_REDHAT="redhat"
[ -z $WORKDIR_PATH ] && WORKDIR_PATH="workdir"
[ -z $WORKING_MODE ] && WORKING_MODE=""
[ -z $TEARDOWN_RUN_ID ] && TEARDOWN_RUN_ID="latest"
[ -z $DEPLOYER_API ] && DEPLOYER_API=$DEFAULT_DEPLOYER 

set_workdir_dependent_variables



## PARSE & CHECK ARGS
if [ $CONTAINER ] 
then
    # check working mode
    [[ (-z $WORKING_MODE ) ]] && stop_if_failed 1 "WORKING_MODE environment variable must be set"
    [[ ( $WORKING_MODE != "C" ) && ( $WORKING_MODE != "T" )  ]] && \
    stop_if_failed 1 "WORKING_MODE value must be either C (create) or T(teardown) $WORKING_MODE is not a valid value"
    # check pull secret
    [[ ($WORKING_MODE == "C") && ( -z $PULL_SECRET ) ]] && stop_if_failed 1 "PULL_SECRET environment variable must be set and must contain a valid base64 encoded pull_secret, please refer to the README.md"
    # check workdir mount write permissions 
    [[ ! -d $WORKDIR_PATH ]] && stop_if_failed 1 "please mount the workdir filesystem, refer to README.md for further instructions"
    [[ ! -w $WORKDIR_PATH ]] && \
    stop_if_failed 1 "please grant write permissions to the host folder mounted as volume, please refer to README.md for further instructions"
    [[ ( $WORKING_MODE == "T" ) && ($TEARDOWN_RUN_ID == "latest") ]] && stop_if_failed 1 "TEARDOWN_RUN_ID must be set in container mode. Please set this value with the run id that you want to teardown, refer to README.md for further instructions"
else
    parse_args $@
    ## VARIABLES SANITY CHECKS
    # WORKING MODE CHECK
    if [[ ! $SHOW_HELP ]]
    then
        [[ (-z $WORKING_MODE ) ]] && echo -e "\nERROR: Working mode must be set\n" && usage
        [[ ( $WORKING_MODE != "C" ) && ( $WORKING_MODE != "T" )  ]] && echo \
        -e "\nERROR: Working mode Must be either -C (creation) or -T (teardown), not $WORKING_MODE\n" && usage
        # CHECK MANDATORY ARGS FOR CREATION
        [[ ($WORKING_MODE == "C" ) && ( ! "$PULL_SECRET_PATH" ) ]] && \
        echo -e "\nERROR: in creation mode argument -p <pull_secret_path> must be provided\n" && usage 
        # CHECK PULL SECRET PATH
        [[ ($WORKING_MODE == "C" ) && ( ! -f $PULL_SECRET_PATH ) ]] && \
        echo -e "\nERROR: $PULL_SECRET_PATH pull secret file not found" && usage
    fi

fi

## LOAD DEPLOYER API
api_load_deployer $DEPLOYER_API

[[ $SHOW_USAGE ]] && usage

## ENTRYPOINT: if everything is ok, run the script.
if [[ $WORKING_MODE == "C" ]]
then
    create $@
elif [[ $WORKING_MODE == "T" ]]
then
    teardown $@
else
    usage
fi


#!/bin/sh

pr_info() {
    if [[ $WORKING_MODE == "C" ]]
    then
        echo "[INF] $1" | (tee -a $LOG_FILE 2>/dev/null) 
    else
        echo "[INF] $1" | (tee -a $TEARDOWN_LOGFILE  2>/dev/null)
    fi
}

pr_error() {
    if [[ $WORKING_MODE == "C" ]]
    then
        echo "[ERR] $1" | (tee -a $LOG_FILE 2>/dev/null)
    else
        echo "[ERR] $1" | (tee -a $TEARDOWN_LOGFILE 2>/dev/null)
    fi
}

pr_end() {
    if [[ $WORKING_MODE == "C" ]]
    then
        echo "[END] $1" | (tee -a $LOG_FILE 2>/dev/null) 
    else
        echo "[END] $1" | (tee -a $TEARDOWN_LOGFILE 2>/dev/null) 
    fi
}

stop_if_failed(){
	EXIT_CODE=$1
	MESSAGE=$2
	if [[ $EXIT_CODE != 0 ]]
	then
		pr_error "$MESSAGE" 
		exit $EXIT_CODE
	fi
}

check_ssh(){
    $NC -z $1 $SSH_PORT > /dev/null 2>&1
    return $?
}

wait_instance_readiness(){
    RES=1
    while [[ $RES != 0 ]] 
    do
        check_ssh $1
        RES=$?
        sleep 1
        pr_info "waiting sshd to become ready on $1, hang on...."
    done
}


get_instance_id_aws() {
    $JQ -r '.Instances[0].InstanceId' $1
}

get_instance_public_ip_aws(){
    INSTANCE_IP=""
    while [ -z $INSTANCE_IP ]
    do
        INSTANCE_IP=`$AWS ec2 describe-instances --instance-ids $1 --query 'Reservations[*].Instances[*].PublicIpAddress' --output text`
    done
    echo "$INSTANCE_IP"
}

get_instance_private_ip_aws(){
    $JQ -r '.Instances[0].PrivateIpAddress' $1
}


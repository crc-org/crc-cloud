#!/bin/sh

pr_info() {
    echo "[INF] $1" | tee -a $LOG_FILE
}

pr_error() {
    echo "[ERR] $1" | tee -a $LOG_FILE
}

pr_end() {
    echo "[END] " | tee -a $LOG_FILE
}


check_ssh(){
    $NC -z $1 22 > /dev/null 2>&1
    return $?
}

wait_instance_readiness(){
    RES=1
    while [[ $RES != 0 ]] 
    do
        check_ssh $1
        RES=$?
        sleep 1
    done
}


get_instance_id() {
    $JQ -r '.Instances[0].InstanceId' $1
}

get_instance_public_ip(){
    INSTANCE_IP=""
    while [ -z $INSTANCE_IP ]
    do
        INSTANCE_IP=`$AWS ec2 describe-instances --instance-ids $1 --query 'Reservations[*].Instances[*].PublicIpAddress' --output text`
    done
    echo "$INSTANCE_IP"
}

get_instance_private_ip(){
    $JQ -r '.Instances[0].PrivateIpAddress' $1
}


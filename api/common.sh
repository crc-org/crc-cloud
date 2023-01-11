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


api_load_deployer() {
    [ ! -d "$PLUGIN_DEPLOYER_FOLDER/$1" ] && stop_if_failed 1 "Deployer API $1 folder not found in $PLUGIN_DEPLOYER_FOLDER/$1, please refer api/README.md for API specifications"
    [ ! -f "$PLUGIN_DEPLOYER_FOLDER/$1/main.sh" ] && stop_if_failed 1 "main.sh not found for deployer $1 in folder $PLUGIN_DEPLOYER_FOLDER/$1, please refer api/README.md for API specifications"
    source $PLUGIN_DEPLOYER_FOLDER/$1/main.sh
    [[ ! `declare -F deployer_create` ]] &&\
    stop_if_failed 1 "deployer_create method not found in main.sh implementation for $1 infrastructure deployer api, please refer api/README.md for API specifications"

    [[ ! `declare -F deployer_teardown` ]] &&\
    stop_if_failed 1 "deployer_teardown method not found in main.sh implementation for $1 infrastructure deployer api, please refer api/README.md for API specifications"

    [[ ! `declare -F deployer_get_eip` ]] &&\
    stop_if_failed 1 "deployer_get_eip method not found in main.sh implementation for $1 infrastructure deployer api, please refer api/README.md for API specifications"

    [[ ! `declare -F deployer_get_iip` ]] &&\
    stop_if_failed 1 "deployer_get_iip method not found in main.sh implementation for $1 infrastructure deployer api, please refer api/README.md for API specifications"

    [[ ! `declare -F deployer_usage` ]] &&\
    stop_if_failed 1 "deployer_usage method not found in main.sh implementation for $1 infrastructure deployer api, please refer api/README.md for API specifications"


    # CHECKS FOR NOT INTERFACE METHDODS NAME NOT STARTING WITH UNDERSCORE, IF THE INTERFACE WILL BE EXTENDED ADD THEM TO THE SWITCH CASE
    for i in `$CAT $PLUGIN_DEPLOYER_FOLDER/$1/main.sh | $GREP -P "^\s*.+\s*\(\)\s*\{"|$SED -r 's/(.+)\(\)\s*\{/\1/'`
    do 
        case "$i" in
            deployer_create);;
            deployer_teardown);;
            deployer_get_eip);;
            deployer_get_iip);;
            deployer_usage);;
            *) 
            [[ ${i::1} != '_' ]] && stop_if_failed 1 "$i() is not a valid private method name, non interface method must start with underscore '_'"
            ;;
        esac
    done
    pr_info "successfully loaded $1 deployer plugin"
}


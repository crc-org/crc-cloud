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
# PARAMS $1: file path $2 plugin name
api_check_method_name() {
    # CHECKS FOR NOT INTERFACE METHDODS NAME NOT STARTING WITH UNDERSCORE, IF THE INTERFACE WILL BE EXTENDED ADD THEM TO THE SWITCH CASE
    # IF THE INTERFACE WILL BE EXTENDED WITH OTHER METHODS ADD THEM HERE
    PLUGIN_NAME_LENGTH=`expr length $2`
    for i in `$CAT $1 | $GREP -P "^\s*.+\s*\(\)\s*\{"|$SED -r 's/(.+)\(\)\s*\{/\1/'`
    do 
        case "$i" in
            deployer_create);;
            deployer_teardown);;
            deployer_get_eip);;
            deployer_get_iip);;
            deployer_usage);;
            deployer_load_dependencies);;
            *) 
            [[ ${i:0:$PLUGIN_NAME_LENGTH+2} != "_$2_" ]] && stop_if_failed 1 "$i() in $1 is not a valid private method name, non interface method must start with _<plugin_name>_ in that case '_$2_' "
            ;;
        esac
    done
}


api_load_deployer() {

    case $1 in
        (*[![:lower:]_]*) stop_if_failed 1 "plugin name must contain only lowercase letters and underscores";;
        (*);;
    esac
    # IF THE INTERFACE WILL BE EXTENDED WITH OTHER METHODS ADD THEM HERE
    [ ! -d "$PLUGIN_DEPLOYER_FOLDER/$1" ] && stop_if_failed 1 "deployer API $1 folder not found in $PLUGIN_DEPLOYER_FOLDER/$1, please refer api/README.md for API specifications"
    [ ! -f "$PLUGIN_DEPLOYER_FOLDER/$1/main.sh" ] && stop_if_failed 1 "main.sh not found for deployer $1 in folder $PLUGIN_DEPLOYER_FOLDER/$1, please refer api/README.md for API specifications"
    source $PLUGIN_DEPLOYER_FOLDER/$1/main.sh
    [[ ! `declare -F deployer_create` ]] &&\
    stop_if_failed 1 "deployer_create method not found in main.sh implementation for $1 infrastructure deployer plugin, please refer api/README.md for API specifications"

    [[ ! `declare -F deployer_teardown` ]] &&\
    stop_if_failed 1 "deployer_teardown method not found in main.sh implementation for $1 infrastructure deployer plugin, please refer api/README.md for API specifications"

    [[ ! `declare -F deployer_get_eip` ]] &&\
    stop_if_failed 1 "deployer_get_eip method not found in main.sh implementation for $1 infrastructure deployer plugin, please refer api/README.md for API specifications"

    [[ ! `declare -F deployer_get_iip` ]] &&\
    stop_if_failed 1 "deployer_get_iip method not found in main.sh implementation for $1 infrastructure deployer plugin, please refer api/README.md for API specifications"

    [[ ! `declare -F deployer_usage` ]] &&\
    stop_if_failed 1 "deployer_usage method not found in main.sh implementation for $1 infrastructure deployer plugin, please refer api/README.md for API specifications"
    
    [[ ! `declare -F deployer_load_dependencies` ]] &&\
    stop_if_failed 1 "deployer_load_dependencies method not found in main.sh implementation for $1 infrastructure deployer plugin, please refer api/README.md for API specifications"


    # CHECKS FOR INTERFACE METHDODS NAME NOT STARTING WITH _<plugin_name>_ , 
    api_check_method_name $PLUGIN_DEPLOYER_FOLDER/$1/main.sh $1
    # CHECKS main.sh INCLUDES FOR METHOD NAME COMPLIANCE
    for i in `$CAT $PLUGIN_DEPLOYER_FOLDER/$1/main.sh | $GREP -P "^\s+source .+\.sh" | sed -r 's#source\s+.+/(.*)#\1#'`
    do 
        api_check_method_name $PLUGIN_DEPLOYER_FOLDER/$1/$i $1
    done

    PLUGIN_ROOT_FOLDER=$PLUGIN_DEPLOYER_FOLDER/$1
    pr_info "successfully loaded $1 deployer plugin"
}




_bash_parse_args() {
    local opts=''
    # NOTE ON ARGUMENTS:
    # DUE TO A GETOPTS STRANGE BEHAVIOUR, IN ORDER TO ADD PLUGIN SPECIFIC ARGUMENTS 
    # IT IS *MANDATORY* TO CONCATENATE WITH THE MAIN SCRIPT ARGUMENT STRING ($BASE_OPTIONS)
    while getopts "${BASE_OPTIONS}t:a:i:" option; do
    case "$option" in
        a) AMI_ID=$OPTARG;pr_info "provided ami-id: $AMI_ID";;
        t) INSTANCE_TYPE=$OPTARG;pr_info "provided instance type: $INSTANCE_TYPE";;
        i) TARGET_INFRASTRUCTURE=$OPTARG; pr_info "selected target infrastructure: $TARGET_INFRASTRUCTURE";;
        :) printf "missing argument for -%s\n" "$OPTARG" >&2; usage;;
    \?) 
    esac
    done
}

# INTERFACE METHODS
deployer_load_dependencies() {
    source $PLUGIN_ROOT_FOLDER/aws.sh
}

deployer_create() {
    [[ ! $CONTAINER ]] && _bash_parse_args $@
    case $TARGET_INFRASTRUCTURE in
        aws)
            _bash_ec2_init
            _bash_create_ec2_resources
        ;;
        gcp)pr_info "coming soon";exit 0;;
        "") stop_if_failed 1 "please select a target infrastructure [aws,gcp], if running in container mode please set environment variable TARGET_INFRASTRUCTURE";;
        *) stop_if_failed 1 "$TARGET_INFRASTRUCTURE is not a valid target infrastructure";;
    esac
}

deployer_teardown() {
    [[ ! $CONTAINER ]] && _bash_parse_args $@
    case $TARGET_INFRASTRUCTURE in
        aws)
            _bash_ec2_init
            _bash_destroy_ec2_resources
        ;;
        gcp)pr_info "coming soon";exit 0;;
        "") stop_if_failed 1 "please select a target infrastructure [aws,gcp], if running in container mode please set environment variable TARGET_INFRASTRUCTURE";;
        *) stop_if_failed 1 "$TARGET_INFRASTRUCTURE is not a valid target infrastructure";;
    esac
}

deployer_get_eip() {
    INSTANCE_ID=`_bash_get_ec2_instance_id $WORKDIR/$INSTANCE_DESCRIPTION`
    EIP=`_bash_get_ec2_instance_public_ip $INSTANCE_ID`
    echo $EIP
}

deployer_get_iip() {
    IIP=`_bash_get_ec2_instance_private_ip $WORKDIR/$INSTANCE_DESCRIPTION`
    echo $IIP
}

deployer_usage() {
        usage="
Deployer Options :

[-a AMI ID] [-t Instance type]
where:
    -a  AMI ID (Amazon Machine Image) from which the VM will be Instantiated (optional, default: $AMI_ID)
    -i  EC2 Instance Type (optional, default; $INSTANCE_TYPE)
"
    echo "$usage"
}
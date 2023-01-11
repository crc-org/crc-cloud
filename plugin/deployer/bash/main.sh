

_bash_parse_args() {
    local opts=''
    # NOTE ON ARGUMENTS:
    # DUE TO A GETOPTS STRANGE BEHAVIOUR, IN ORDER TO ADD PLUGIN SPECIFIC ARGUMENTS 
    # IT IS *MANDATORY* TO CONCATENATE WITH THE MAIN SCRIPT ARGUMENT STRING ($BASE_OPTIONS)
    while getopts "${BASE_OPTIONS}t:a:" option; do
    case "$option" in
        a) AMI_ID=$OPTARG;pr_info "provided ami-id: $AMI_ID";;
        t) INSTANCE_TYPE=$OPTARG;pr_info "provided instance type: $INSTANCE_TYPE";;
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
    _bash_parse_args $@
    _bash_create_ec2_resources
}

deployer_teardown() {
    _bash_destroy_ec2_resources
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
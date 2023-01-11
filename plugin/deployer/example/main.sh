# THIS METHOD IS NOT PART OF THE INTERFACE HAS BEEN PUT HERE AS AN EXAMPLE ON HOW
# TO PARSE ADDITIONAL ARGUMENTS FOR A DEPLOYER PLUGIN
# NON INTERFACE METHODS *MUST* START WITH _<plugin_name>_  SO IN THAT CASE, CONSIDERING THAT THE PLUGIN NAME IS example
# ALL THE METHODS MUST START WITH _example_ TO AVOID CONFLICTS AND KEEP THE CODE CLEAR

_example_parse_args() {
    local opts=''
    #NOTE ON ARGUMENTS:
    #DUE TO A GETOPTS STRANGE BEHAVIOUR, IN ORDER TO ADD PLUGIN SPECIFIC ARGUMENTS 
    #IT IS *MANDATORY* TO CONCATENATE WITH THE MAIN SCRIPT ARGUMENT STRING ($BASE_OPTIONS)
    #IN THIS EXAMPLE WE ADDED -t AS A "DEPLOYER SPECIFIC" OPTION
    while getopts "${BASE_OPTIONS}t:" option; do
    case "$option" in
    #example -t argument
        t)echo "I'm -t argument $OPTARG";exit 0;; 
        :) printf "missing argument for -%s\n" "$OPTARG" >&2; usage;;
    \?) 
    esac
    done
    echo $OPTIND
}




# INTERFACE METHODS

deployer_create() {
    #all the command line args will be passed to that function
    pr_info "creates the infrastructure"
    pr_info "\$CONTAINER: $CONTAINER"
    pr_info "\$RANDOM_SUFFIX: $RANDOM_SUFFIX"
    pr_info "\$WORKDIR: $WORKDIR"
    pr_info "\$PLUGIN_ROOT_FOLDER: $PLUGIN_ROOT_FOLDER"
    exit 0
}

deployer_load_dependencies() {
    pr_info "loads (if needed otherwise keep it empty) the plugin dependencies"
    source $PLUGIN_ROOT_FOLDER/example_include.sh
}

deployer_teardown() {
    #all the command line args will be passed to that function
    pr_info "destroys infrastructure"
    exit 0
}

deployer_get_eip() {
    echo "return the external (public ip) of the VM"
}

deployer_get_iip() {
    echo "return internal ip (within the cloud infrastructure) of the VM"
}

deployer_usage() {
    echo "prints the usage of the deployer including all its cli parameters"
}
# CRC-Cloud Infrastructure Deployer API

The Infrastructure Deployer API has been designed to abstract the Infrastructure provisioning from the OpenShift instance provisioning. The first version of CRC-Cloud was relying on AWS and AWS CLI, but as soon as the project was gaining interest, we started considering to support other cloud providers so we decided to implement this abstraction to easily implement other deployment technologies such as IaC tools like [Ansible](https://www.redhat.com/it/engage/delivery-with-ansible-20170906?sc_cid=7013a000002w14JAAQ&gclid=EAIaIQobChMIwLPlpZG9_AIVA5zVCh2EPw9VEAAYASAAEgJJXfD_BwE&gclsrc=aw.ds), [Terraform](https://terraform.io), [Pulumi](https://https://www.pulumi.com/) etc.

## Plugin loading and API implementation
In order to be loaded, an Infrastructure Deployer Plugin must have a folder in ```plugin/deployer```. This folder must have the name of the plugin and the plugin name must contain **only** lowercase letters and underscores, that will be passed to ```crc-cloud.sh``` with the ```-D``` option, so for example, if you want to create a plugin named ```my_deployer``` the plugin code and resources will be stored in ```<openspot_path>/plugin/deployer/my-deployer```.
You can find an example implementation from which start to develop a new plugin in ```plugin/deployer/example``` (and you can even run it!!).
The plugin folder must contain a ```main.sh``` script that is the entrypoint of the plugin. The ```main.sh``` **must** implement the following methods:

```

deployer_load_dependencies() {
    pr_info "loads (if needed otherwise keep it empty) the plugin dependencies"
    source $PLUGIN_ROOT_FOLDER/example_include.sh
}

deployer_create() {
    #all the command line args will be passed to that function
    pr_info "creates the infrastructure"
    exit 0
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
    echo "prints the usage of the infrastructure deployer plugin including all its cli parameters"
}
```

The **CRC-Cloud** engine will check if these methods are implemented and will eventually exit it they're not.
All the plugin specific resources (other scripts, IaC definitions etc.) should be kept into the plugin folder and the paths that refer to them must be valorized accordingly by the developer.

## Global Variables

The **CRC-Cloud** engine will expose to the plugin some variables that must be used to keep the logic consistent with the engine itself

| Variable | Type| Description |
| --- | --- | --- |
| $CONTAINER | Boolean | It's valorized if the script is running inside a container |
| $RANDOM_SUFFIX | String | It's the random suffix applied to the resources created inside the cloud provider in order to avoid conflicts with other **CRC-Cloud** instances running in the same namespace (can be ignored if the deployment methods provides it's own logic) |
| $WORKDIR | String | That's the folder where all the deployment status infos must be stored (will be created by the engine) |
| $PLUGIN_ROOT_FOLDER | String | That's the folder containing the loaded plugin, this can be used as starting path for plugin resources | 

 ## *Private* methods names conventions

 In order to increase code readability and avoid conflicts non interface method names (kinda private) must start with _<plugin_name>_ for example if you name your plugin *example_plugin* all the non interface methods defined inside the plugin must start with _example_plugin_
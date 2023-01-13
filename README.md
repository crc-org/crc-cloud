# CRC-Cloud
### Disposable OpenShift instances on cloud in minutes

![CRC Cloud](assets/crc-cloud.png)

## Disclaimer
This project has been developed for **experimental** purpose only and it's not **absolutely** meant to run production clusters. The author is not responsible in any manner of any cost on which the user may incur for inexperience or software failure.
Before running the script **be sure** to have an adequate experience to safely create and destroy resources on AWS or any other cloud provider that will be supported **without** the help of this software in order to recover **manually** from possible issues.

## Why? (TL;DR)
I needed to test a chaos engineering tool (https://github.com/redhat-chaos/krkn) against a disposable (single-node) OpenShift cluster that could have been setup and destroyed from inside a CI/CD pipeline as quick as possible, unattended and with a reasonable cost per run.
I stumbled upon OpenSpot (https://github.com/ksingh7/openspot) made by my colleagues  [@ksingh7](https://github.com/ksingh7) and [@praveenkumar](https://github.com/praveenkumar). I found the idea amazing, unfortunately was relying on EC2 Spot Instances, that, if from a cost perspective are more affordable, do not guarantee that the machine is instantiated.
Moreover the solution was based on CRC that creates a qemu VM to run the (single-node) cluster, so bare metal instances were needed and the startup time was too long for the purpose.
We had a meeting and they gave me all the detailed instructions on how to run the qemu image directly on AWS standard EC2 instances and configure properly the OpenShift single-node cluster, only the code was missing....

## Cloud Providers
For the moment only AWS is supported. Other will be added soon.
<br/>
<br/>
**Note:** AWS AMIs (Amazon Machine Images) are regional resources so,for the moment, the only supported region is **us-west-2**. In the next few days the AMI will be copied to other regions, please be patient, it will take a while.

## Usage

#### Prerequisites
<a name="prereq"></a>
The basic requirements to run a single-node OpenShift cluster with **CRC-Cloud** are:
- register a Red Hat account and get a pull secret from https://console.redhat.com/openshift/create/local 
- create an access key for your AWS account and grab the *ACCESS_KEY_ID* and the *SECRET_ACCESS_KEY* (instructions can be found [here](https://docs.aws.amazon.com/general/latest/gr/aws-sec-cred-types.html)) 
<br/>
<br/>

The AWS instance type of choice is *c6in.2xlarge* with 8vcpu and 16 GB of RAM.
This instance will cost ~0.45$ per hour (price may vary depending from the region) and will take ~11 minutes to have a working cluster.
Increasing or decreasing the resources will affect the deployment time together with the price per hour. If you want to change instance type keep in mind that the minimum hardware requirements to run CRC (on which this solution is based) are 4vcpus and 8GB of RAM, please refer to the [documentation](https://developers.redhat.com/blog/2019/09/05/red-hat-openshift-4-on-your-laptop-introducing-red-hat-codeready-containers) for further informations.

**WARNING:** Running VM instances will cost you **real money** so be extremely careful to verify that all the resources instantiated are **removed** once you're done and remember that you're running them at **your own risk and cost**



### Containers (the easy way)

Running **CRC-Cloud** from a container (podman/docker) is strongly recommended for the following reasons:
- Compatible with any platform (Linux/MacOs/Windows)
- No need to satisfy any software dependency in you're OS since everything is packed into the container
- In CI/CD systems (eg. Jenkins) won't be necessary to propagate dependencies to the agents (only podman/docker needed)
- In Cloud Native CI/CD systems (eg. Tekton) everything runs in container so that's the natural choice

#### Working directory
<a name="workdir"></a>
In the working directory that will be mounted into the container, **CRC-Cloud** will store all the cluster metadata including those needed to teardown the cluster once you'll be done.
Per each run **CRC-Cloud** will create a folder named with the run timestamp, this folder name will be referred as *TEARDOWN_RUN_ID* and will be used to cleanup the cluster in teardown mode and to store all the logs and the infos related to the cluster deployment.

Please **be careful** on deleting the working directory content because without the metadata **CRC-Cloud** won't be able to teardown the cluster and associated resources from AWS.

**NOTE (podman only):** In order to make the mounted workdir read-write accessible from the container is need to change the SELinux security context related to the folder with the following command 
```chcon -Rt svirt_sandbox_file_t <HOST_WORKDIR_PATH>```

#### Single node cluster creation
```
<podman|docker> run -v <HOST_WORKDIR_PATH>:/workdir\
 -e WORKING_MODE=C\
 -e PULL_SECRET="`base64 <PULL_SECRET_PATH>`"\
 -e AWS_ACCESS_KEY_ID=<AWS_ACCESS_KEY_ID>\
 -e AWS_SECRET_ACCESS_KEY=<AWS_SECRET_ACCESS>\
 -e AWS_DEFAULT_REGION=<AWS_REGION_OF_YOUR_CHOICE>\
 -ti quay.io/crcont/crc-cloud
```

#### Single node cluster teardown
```
<podman|docker> run -v <HOST_WORKDIR_ABSOLUTE_PATH>:/workdir\
 -e WORKING_MODE=T\
 -e TEARDOWN_RUN_ID=<TEARDOWN_RUN_ID>\ 
 -e AWS_ACCESS_KEY_ID=<AWS_ACCESS_KEY_ID>\
 -e AWS_SECRET_ACCESS_KEY=<AWS_SECRET_ACCESS>\
 -e AWS_DEFAULT_REGION=us-west-2\
 quay.io/crcont/crc-cloud
```
(check [here](#workdir) for **TEARDOWN_RUN_ID** infos and **WORKDIR** setup instructions )

#### Environment variables
Environment variables will be passed to the container from the command line invocation with the ```-e VARIABLE=VALUE``` option that you can find above.
##### Mandatory Variables

**Cluster creation**

|  VARIABLE | DESCRIPTION  |
|---|---|
|  WORKING_MODE | C (creation mode)  |
|  PULL_SECRET | base64 string of the Red Hat account pull secret ( it is recommended to use the command substitution to generate the string as described above)  |
| AWS_ACCESS_KEY_ID  | AWS access key (infos [here](#prereq))  |
| AWS_SECRET_ACCESS_KEY | AWS secret access key (infos [here](#prereq)) |
| AWS_DEFAULT_REGION | AWS region where the cluster will be deployed ( currently us-west-2 is the only supported)



**Cluster teardown**

| VARIABLE | DESCRIPTION |
| --- | ---|
| WORKING_MODE | T (teardown) |
| TEARDOWN_ID | the name of the folder created inside the working directory, containing all the metadata needed to teardown the cluster (if specified in creation will match with CREATE_RUN_ID)|
| AWS_ACCESS_KEY_ID  | AWS access key (infos [here](#prereq))  |
| AWS_SECRET_ACCESS_KEY | AWS secret access key (infos [here](#prereq)) |
| AWS_DEFAULT_REGION | AWS region where the cluster has been deployed ( currently us-west-2 is the only supported) |

#### Optional Variables


|  VARIABLE | DESCRIPTION  |
|---|---|
| CREATE_RUN_ID | Overrides the folder where the **CRC-Cloud** will put the output files, if not set will use the current timestamp, useful for CI/CD pipelines or everywhere there is the need to know this value upfront |
|  PASS_DEVELOPER |  overrides the default password (developer) for developer account  |
|  PASS_KUBEADMIN | overrides the default password (kubeadmin) for kubeadmin account   |
|  PASS_REDHAT |  overrides the default password (redhat) for redhat account |
| INSTANCE_TYPE | overrides the default AWS instance type (c6in.2xlarge, infos [here](#prereq)) |



### Linux Bash (the hard path)
#### Dependencies 
To run **CRC-Cloud** from your command line you must be on Linux, be sure to have installed and configured the following programs in your box

- bash (>=v4)
- AWS CLI
- jq
- md5sum
- curl
- head
- ssh-keygen
- GNU sed
- nc (netcat)
- ssh client
- scp
  
#### Single node cluster creation
Once copied and downloaded the pull secret somewhere in your filesystem you'll be able to run the cluster with

```./crc-cloud.sh -C -p <pull_secret_path>```

A folder with the run id will be created under ```<openspot_path>/workspace``` containing all the logs, the keypair needed to login into the VM, and the VM metadata. The last run will be also linked automatically to ```<openspot_path>/latest```
<br/>
<br/>
**WARNING:** if you delete the working directory **CRC-Cloud** won't be able to teardown the cluster so be **extremely careful** with the workspace folder content.
<br/>
<br/>
at the end of the process the script will print the public address of the console.
Below you'll find all the options available

```
./crc-cloud.sh -C -p pull secret path [-d developer user password] [-k kubeadmin user password] [-r redhat user password] [-a AMI ID] [-t Instance type]
where:
    -C  Cluster Creation mode
    -p  pull secret file path (download from https://console.redhat.com/openshift/create/local) 
    -d  developer user password (optional, default: developer)
    -k  kubeadmin user password (optional, default: kubeadmin)
    -r  redhat    user password (optional, default: redhat)
    -a  AMI ID (Amazon Machine Image) from which the VM will be Instantiated (optional, default: ami-0569ce8a44f2351be)
    -i  EC2 Instance Type (optional, default; c6in.2xlarge)
    -h  show this help text
```
#### Single node cluster teardown
To teardown the single node cluster the basic command is 
```./crc-cloud.sh -T```
this will refer to the *latest* run found in ```<openspot_path>/workspace```, if you have several run folders in your workspace, you can specify the one you want to teardown with the parameter ```-v <run_id>``` where ```<run_id>``` corresponds to the numeric folder name containing the metadata of the cluster that will be deleted

```
./crc-cloud.sh -T [-v run id]
    -T  Cluster Teardown mode
    -v  The Id of the run that is gonna be destroyed, corresponds with the numeric name of the folders created in workdir (optional, default: latest)
    -h  show this help text 

```
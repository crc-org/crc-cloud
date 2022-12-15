# OpenSpot-NG
### Disposable OpenShift instances on AWS in minutes

| ![space-1.jpg](assets/openspot_ng_logo.png) |
|:--:|
| <sub><sup>kindly created by OpenAI DALL-E (https://openai.com/dall-e-2) </sup></sub>|

## Disclaimer
This project has been developed for **experimental** purpose only and it's not **absolutely** meant to run production clusters. The author is not responsible in any manner of any cost on which the user may incur for inexperience or software malfunctioning.
Before running the script **be sure** to have an adequate experience to safely create and destroy resources on AWS or any other cloud provider that will be supported, **without** the help of this software in order to recover **manually** to possible issues.

## Why? (TL;DR)
I needed to test a chaos engineering tool (https://github.com/redhat-chaos/krkn) against a disposable (single-node) OpenShift cluster that could have been setup and destroyed from inside a CI/CD pipeline as quick as possible, unattended and with a reasonable cost per run.
I stumbled upon OpenSpot (https://github.com/ksingh7/openspot) made by my colleagues  [@ksingh7](https://github.com/ksingh7) and [@praveenkumar](https://github.com/praveenkumar). I found the idea amazing, unfortunately was relying on EC2 Spot Instances, that, if from a cost perspective are more affordable, do not guarantee that the machine is instantiated.
Moreover the solution was based on CRC that creates a qemu VM to run the (single-node) cluster, so bare metal instances were needed and the startup time was too long for the purpose.
We had a meeting and they gave me all the instructions on how to run the qemu image directly in AWS and configure properly the OpenShift single-node cluster, only the code was missing....

## Cloud Providers
For the moment only AWS is supported. Other will be added soon.
<br/>
<br/>
**Note:** AWS AMIs (Amazon Machine Images) are regional resources so,for the moment, the only supported region is **us-west-2**. In the next few days the AMI will be copied to other regions, please be patient, it will take a while.

## Usage
### Prerequisites
To run **OpenSpot-NG** from your command line you must be on Linux, very soon will be ready a containerized version that will allow you to run it from every OS that supports Podman/Docker.

Be sure to have installed and configured the following programs in your box

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

### Linux (Bash)
#### Single node cluster creation

The basic requirements to run a single-node OpenShift cluster with **OpenSpot-NG** are:
- register a Red Hat account and get a pull secret from https://console.redhat.com/openshift/create/local 
- configure the AWS CLI with your AWS credentials (region **us-west-2**)
<br/>
<br/>

**WARNING:** Running VM instances will cost you **real money** so be extremely careful to verify that all the resources instantiated are **removed** once you're done and remember that you're running them at **your own risk and cost**

<br/>
<br/>
Once copied and downloaded the pull secret somewhere in your filesystem you'll be able to run the cluster with

```./openspot-ng.sh -C -p <pull_secret_path>```

A folder with the run id will be created under ```<openspot_path>/workspace``` containing all the logs, the keypair needed to login into the VM, and the VM metadata. The last run will be also linked automatically to ```<openspot_path>/latest```
<br/>
<br/>
**WARNING:** if you delete the working directory **OpenSpot-NG** won't be able to teardown the cluster so be **extremely careful** with the workspace folder content.
<br/>
<br/>
at the end of the process the script will print the public address of the console.
Below you'll find all the options available

```
./openspot-ng.sh -C -p pull secret path [-d developer user password] [-k kubeadmin user password] [-r redhat user password] [-a AMI ID] [-t Instance type]
where:
    -C  Cluster Creation mode
    -p  CRC pull secret file path (download from https://console.redhat.com/openshift/create/local) 
    -d  CRC developer user password (optional, default: developer)
    -k  CRC kubeadmin user password (optional, default: kubeadmin)
    -r  CRC redhat    user password (optional, default: redhat)
    -a  AMI ID (Amazon Machine Image) from which the VM will be Instantiated (optional, default: ami-0569ce8a44f2351be)
    -i  EC2 Instance Type (optional, default; c6in.2xlarge)
    -h  show this help text
```
#### Single node cluster teardown
To teardown the single node cluster the basic command is 
```./openspot-ng.sh -T```
this will refer to the *latest* run found in ```<openspot_path>/workspace```, if you have several run folders in your workspace, you can specify the one you want to teardown with the parameter ```-v <run_id>``` where ```<run_id>``` corresponds at the numeric folder name containing the metadata of the cluster that will be deleted

```
./openspot-ng.sh -T [-v run id]
    -T  Cluster Teardown mode
    -v  The Id of the run that is gonna be destroyed, corresponds with the numeric name of the folders created in workdir (optional, default: latest)
    -h  show this help text 

```

## Containers
### coming soon.

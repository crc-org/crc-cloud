# CRC Cloud - Runs Containers in the Cloud

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

#### Single node cluster creation
```
podman run -d --rm \
    -v ${PWD}:/workspace:z \
    -e AWS_ACCESS_KEY_ID=XXX \
    -e AWS_SECRET_ACCESS_KEY=XXX \
    -e AWS_DEFAULT_REGION=eu-west-1 \
    quay.io/ariobolo/crc-cloud:v0.0.1 create \
        --project-name "crc-ocp412" \
        --backed-url "file:///workspace" \
        --output "/workspace" \
        --provider "aws" \
        --aws-ami-id "ami-0ab26eb25f41697ef" \
        --pullsecret-filepath "/workspace/pullsecret" \
        --key-filepath "/workspace/id_ecdsa"
```

#### Single node cluster teardown
```
podman run -d --rm \
    -v ${PWD}:/workspace:z \
    -e AWS_ACCESS_KEY_ID=XXX \
    -e AWS_SECRET_ACCESS_KEY=XXX \
    -e AWS_DEFAULT_REGION=eu-west-1 \
    quay.io/ariobolo/crc-cloud:v0.0.1 destroy \
        --project-name "crc-ocp412" \
        --backed-url "file:///workspace" \
        --provider "aws"
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
| TEARDOWN_ID | the name (unix timestamp format) of the folder created inside the working directory, containing all the metadata needed to teardown the cluster |
| AWS_ACCESS_KEY_ID  | AWS access key (infos [here](#prereq))  |
| AWS_SECRET_ACCESS_KEY | AWS secret access key (infos [here](#prereq)) |
| AWS_DEFAULT_REGION | AWS region where the cluster has been deployed ( currently us-west-2 is the only supported) |

#### Optional Variables


|  VARIABLE | DESCRIPTION  |
|---|---|
|  PASS_DEVELOPER |  overrides the default password (developer) for developer account  |
|  PASS_KUBEADMIN | overrides the default password (kubeadmin) for kubeadmin account   |
|  PASS_REDHAT |  overrides the default password (redhat) for redhat account |
| INSTANCE_TYPE | overrides the default AWS instance type (c6in.2xlarge, infos [here](#prereq)) |


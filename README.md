# CRC Cloud - Runs Containers in the Cloud

## Disposable OpenShift instances on cloud in minutes

![CRC Cloud](assets/crc-cloud.png)

This project is stumbled upon [OpenSpot](https://github.com/ksingh7/openspot) made by [@ksingh7](https://github.com/ksingh7) and all the improvements made by [@tsebastiani](https://github.com/tsebastiani) creating the next generation for openspot were he got rid off bare metal hard requirement for running the single-node cluster on the cloud.

## Disclaimer

This project has been developed for **experimental** purpose only and it's not **absolutely** meant to run production clusters.

The authors are not responsible in any manner of any cost on which the user may incur for inexperience or software failure.

Before running the script **be sure** to have an adequate experience to safely create and destroy resources on AWS or any other cloud provider that will be supported **without** the help of this software in order to recover **manually** from possible issues.

## Overview  

This is a side project of [`Openshift Local` formerly `CRC`](https://github.com/crc-org), while `CRC` and `crc cli` main purpose is spin `Openshift Single Node` clusters on local development environments (it works multi-platform and multi-arch), `crc-cloud` will offer those clusters on cloud (multi-provider).

The following diagram shows what is the expected interaction between an user of `crc-cloud` and the assets provided by `CRC`:

![crc-cloud flow](docs/crc-cloud-flow.svg?raw=true)

## Usage

To facilite the usage of `crc-cloud`, a [container image](https://quay.io/repository/crcont/crc-cloud) is offered with all required dependecies. Using the container all 3 supported operation can be executed

### Authetication  

All operations require to set the authentication mechanism in place. As so any `aws` authentication mechanism is supported by `crc-cloud`:

- long term credentials `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` as environment variables  
- short lived credentials (in addition to `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` would require `AWS_SESSION_TOKEN`)
- credentials on config file (default file ~/.aws/config), in case of multiple profiles it will also accepts `AWS_PROFILE`

### Restrictions

The `import` operation downloads and transform the bundle offered by crc into an image supported by `AWS`, as so there are some disk demanding operation. So there is a requirement of at least 70G free on disk to run this operation.  

The AWS instance type of choice is *c6a.2xlarge* with 8vcpu and 16 GB of RAM. This will be customizable in the future, for the moment this fixed type imposes some [restrictions](https://aws.amazon.com/about-aws/whats-new/2022/12/amazon-ec2-m6a-c6a-instances-additional-regions/) on available regions to run crc cloud, those regions are:

- us-east-1 and us-east-2  
- us-west-1 and us-west-2  
- ap-south-1, ap-southeast-1, ap-southeast-2 and ap-northeast-1  
- eu-west-1, eu-central-1 and eu-west-2  

### Operations

#### Import

`import` operation uses crc official bundles, transform them and import as an AMI on the user account. It is required to run `import` operation on each region where the user wants to sping the cluster.  

Usage:

```bash
import crc cloud image

Usage:
  crc-cloud import [flags]

Flags:
      --backed-url string              backed for stack state. Can be a local path with format file:///path/subpath or s3 s3://existing-bucket
      --bundle-shasumfile-url string   custom url to download the shasum file to verify the bundle artifact
      --bundle-url string              custom url to download the bundle artifact
  -h, --help                           help for import
      --output string                  path to export assets
      --project-name string            project name to identify the instance of the stack
      --provider string                target cloud provider
```

Outputs:

- `image-id` file with the ami-id of the imported image  
- `id_ecdsa` this is key required to spin the image. (It will be required on `create` operation, is user responsability to store this key)  

Sample

```bash
podman run -d --rm \
    -v ${PWD}:/workspace:z \
    -e AWS_ACCESS_KEY_ID=${access_key_value} \
    -e AWS_SECRET_ACCESS_KEY=${secret_key_value} \
    -e AWS_DEFAULT_REGION=eu-west-1 \
    quay.io/crcont/crc-cloud:v0.0.2 import \
        --project-name "ami-ocp412" \
        --backed-url "file:///workspace" \
        --output "/workspace" \
        --provider "aws" \
        --bundle-url "https://developers.redhat.com/content-gateway/file/pub/openshift-v4/clients/crc/bundles/openshift/4.12.0/crc_libvirt_4.12.0_amd64.crcbundle" \
        --bundle-shasumfile-url "https://developers.redhat.com/content-gateway/file/pub/openshift-v4/clients/crc/bundles/openshift/4.12.0/sha256sum.txt"

```

#### Create  

`create` operation is responsible for create all required resources on the cloud provider to spin the Openshift Single Node Cluster.  

Usage:

```bash
create crc cloud instance on AWS

Usage:
  crc-cloud create aws [flags]

Flags:
      --aws-ami-id string   AMI identifier
  -h, --help                help for aws

Global Flags:
      --backed-url string            backed for stack state. Can be a local path with format file:///path/subpath or s3 s3://existing-bucket
      --key-filepath string          path to init key obtained when importing the image
      --output string                path to export assets
      --project-name string          project name to identify the instance of the stack
      --pullsecret-filepath string   path for pullsecret file
```

Outputs:

- `host` file containing host address running the cluster
- `username` file containing the username to connect the remote host
- `id_rsa` key to connect the remote host
- `password` password generated for `kubeadmin` and `developer` default cluster users

Sample

```bash
podman run -d --rm \
    -v ${PWD}:/workspace:z \
    -e AWS_ACCESS_KEY_ID=${access_key_value} \
    -e AWS_SECRET_ACCESS_KEY=${secret_key_value} \
    -e AWS_DEFAULT_REGION=eu-west-1 \
    quay.io/crcont/crc-cloud:v0.0.2 create aws \
        --project-name "crc-ocp412" \
        --backed-url "file:///workspace" \
        --output "/workspace" \
        --aws-ami-id "ami-xxxx" \
        --pullsecret-filepath "/workspace/pullsecret" \
        --key-filepath "/workspace/id_ecdsa"
```

#### Destroy

`destroy` operation will remove any resource created at the cloud provider, it uses the files holding the state of the infrastructure which has been store at location defined by parameter `backed-url` on `create` operation.  

Usage:

```bash
destroy crc cloud instance

Usage:
  crc-cloud destroy [flags]

Flags:
      --backed-url string     backed for stack state. Can be a local path with format file:///path/subpath or s3 s3://existing-bucket
  -h, --help                  help for destroy
      --project-name string   project name to identify the instance of the stack
      --provider string       target cloud provider
```

Sample

```bash
podman run -d --rm \
    -v ${PWD}:/workspace:z \
    -e AWS_ACCESS_KEY_ID=${access_key_value} \
    -e AWS_SECRET_ACCESS_KEY=${secret_key_value} \
    -e AWS_DEFAULT_REGION=eu-west-1 \
    quay.io/crcont/crc-cloud:v0.0.2 destroy \
        --project-name "crc-ocp412" \
        --backed-url "file:///workspace" \
        --provider "aws" 
```

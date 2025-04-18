# Deploy CRC cloud directly on CRC host

## Introduction

The CRC cloud tool deploys the CRC host using its own binary, but it
requires access to the cloud provider credentials. That functionality
is useless when the CRC cloud needs to be used in external CI, where
the crc-cloud tool can not be used - especially when the CI is responsible
to check how many VMs are spawned or checking job results.

## CRC QCOW2 image

There is a simple way to bootstrap the CRC without using the crc-cloud
tool and upload it to your cloud provider:

- download libvirt bundle (you can see url when crc is in debug mode)
- extract libvirt bundle using tar command with zst:

```shell
tar xaf crc_libvirt_$VERSION_amd64.crcbundle
```

- upload crc.qcow2 image to your cloud provider
- take the id_ecdsa_crc file located in the root directory of
  the extracted crcbundle archive after unpacking the crcbundle file

## Bootstrap crc-cloud directly on host

Now after spawning VM using the `crc.qcow2` image:

- prepare `inventory.yaml` file:

```shell
CRC_VM_IP=<ip address>

cat << EOF > inventory.yaml
---
all:
  hosts:
    crc:
      ansible_port: 22
      ansible_host: $CRC_VM_IP
      ansible_user: core
      ansible_ssh_private_key_file: upacked_crcbundle_dir/id_ecdsa
  vars:
    alternative_domain: true
    pass_developer: 12345678
    pass_kubeadmin: 12345678
    pass_redhat: 12345678
    openshift_pull_secret: |
      < PULL SECRET >
EOF
```

- clone crc-cloud project

```shell
git clone https://github.com/crc-org/crc-cloud
```

- run playbook to bootstrap the container that later would start deploy-crc-cloud role

```shell
ansible-playbook -i inventory.yaml crc-cloud/ansible/playbooks/bootstrap.yaml
```

Then just wait until the bootstrap container finish :)
You can follow Ansible execution steps by checking the container logs
on the remote CRC host:

```shell
sudo podman logs -f crc-cloud-bootstrap
```

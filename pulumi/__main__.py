"""A Python Pulumi Project used to test the creation of an AWS EC2 instance for crc purposes"""

import pulumi
import pulumi_tls as tls
import pulumi_aws as aws

config = pulumi.Config()
key_name = config.require('key_name')
instance_type = config.require('instance_type')
instance_ami = config.require('instance_ami')

ssh_key = tls.PrivateKey(
    "generated",
    algorithm="RSA",
    rsa_bits=4096,
)

aws_key = aws.ec2.KeyPair(
    "generated",
    key_name=key_name,
    public_key=ssh_key.public_key_openssh,
    opts=pulumi.ResourceOptions(parent=ssh_key),
)

server = aws.ec2.Instance('my-pulumi-instance',
    key_name = aws_key.key_name,
    instance_type=instance_type,
    ami=instance_ami)

pulumi.export('private_key_pem', ssh_key.private_key_pem)
pulumi.export('public_ip', server.public_ip)
pulumi.export('public_dns', server.public_dns)

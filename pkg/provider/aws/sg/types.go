package sg

import "github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"

type IngressRule struct {
	Description string
	FromPort    int
	ToPort      int
	Protocol    string
	CidrBlocks  string
	SG          *ec2.SecurityGroup
}

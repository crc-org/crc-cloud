package sg

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	NETWORKING_CIDR_ANY_IPV4 string = "0.0.0.0/0"
)

func Create(ctx *pulumi.Context, rules []IngressRule,
	name, description string) (pulumi.IDOutput, error) {
	t := true
	vpc, err := ec2.LookupVpc(ctx, &ec2.LookupVpcArgs{Default: &t})
	if err != nil {
		return pulumi.IDOutput{}, err
	}
	sg, err := ec2.NewSecurityGroup(
		ctx,
		name,
		&ec2.SecurityGroupArgs{
			Description: pulumi.String(description),
			VpcId:       pulumi.String(vpc.Id),
			Ingress:     getSecurityGroupIngressArray(rules),
			Egress:      ec2.SecurityGroupEgressArray{egressAll},
			Tags: pulumi.StringMap{
				"Name": pulumi.String(name),
			},
		})
	if err != nil {
		return pulumi.IDOutput{}, err
	}
	return sg.ID(), nil
}

func getSecurityGroupIngressArray(rules []IngressRule) (sgia ec2.SecurityGroupIngressArray) {
	for _, r := range rules {
		args := &ec2.SecurityGroupIngressArgs{
			Description: pulumi.String(r.Description),
			FromPort:    pulumi.Int(r.FromPort),
			ToPort:      pulumi.Int(r.ToPort),
			Protocol:    pulumi.String(r.Protocol),
		}
		if r.SG != nil {
			args.SecurityGroups = pulumi.StringArray{r.SG.ID()}
		} else if len(r.CidrBlocks) > 0 {
			args.CidrBlocks = pulumi.StringArray{pulumi.String(r.CidrBlocks)}
		} else {
			args.CidrBlocks = pulumi.StringArray{pulumi.String(NETWORKING_CIDR_ANY_IPV4)}
		}
		sgia = append(sgia, args)
	}
	return sgia
}

var egressAll = &ec2.SecurityGroupEgressArgs{
	FromPort: pulumi.Int(0),
	ToPort:   pulumi.Int(0),
	Protocol: pulumi.String("-1"),
	CidrBlocks: pulumi.StringArray{
		pulumi.String(NETWORKING_CIDR_ANY_IPV4),
	}}

package aws

import (
	"fmt"
	"strconv"

	"github.com/crc/crc-cloud/pkg/bundle"
	"github.com/crc/crc-cloud/pkg/bundle/setup"
	"github.com/crc/crc-cloud/pkg/manager/context"
	providerAPI "github.com/crc/crc-cloud/pkg/manager/provider/api"
	"github.com/crc/crc-cloud/pkg/provider/aws/sg"
	"github.com/crc/crc-cloud/pkg/util"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi-tls/sdk/v4/go/tls"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type createRequest struct {
	projectName               string
	amiID                     string
	instanceType              string
	diskSize                  int
	bootingPrivateKeyFilePath string
	ocpPullSecretFilePath     string
}

func fillCreateRequest(projectName, bootingPrivateKeyFilePath, ocpPullSecretFilePath string,
	args map[string]string) (*createRequest, error) {
	amiIDValue, ok := args[amiID]
	if !ok {
		return nil, fmt.Errorf("amiID not found")
	}
	it := ocpInstanceType
	if customInstanceType, ok := args[instanceType]; ok {
		it = customInstanceType
	}
	ds := ocpDefaultRootBlockDeviceSize
	if customDiskSizeAsString, ok := args[diskSize]; ok {
		customDiskSize, err := strconv.Atoi(customDiskSizeAsString)
		if err != nil {
			return nil, fmt.Errorf("error creating request for cluster machine: %v", err)
		}
		ds = customDiskSize
	}
	return &createRequest{
		projectName:               projectName,
		amiID:                     amiIDValue,
		instanceType:              it,
		diskSize:                  ds,
		bootingPrivateKeyFilePath: bootingPrivateKeyFilePath,
		ocpPullSecretFilePath:     ocpPullSecretFilePath}, nil
}

func (r createRequest) runFunc(ctx *pulumi.Context) error {
	securityGroupsIds, err := securityGroupsIds(ctx)
	if err != nil {
		return err
	}
	privateKey, awsKeyPair, err := createKey(ctx)
	if err != nil {
		return err
	}
	args := ec2.InstanceArgs{
		Ami:                      pulumi.String(r.amiID),
		InstanceType:             pulumi.String(r.instanceType),
		KeyName:                  awsKeyPair.KeyName,
		AssociatePublicIpAddress: pulumi.Bool(true),
		VpcSecurityGroupIds:      securityGroupsIds,
		RootBlockDevice: ec2.InstanceRootBlockDeviceArgs{
			VolumeSize: pulumi.Int(r.diskSize),
		},
		Tags: context.GetTags(),
	}
	instance, err := ec2.NewInstance(ctx, r.projectName, &args)
	if err != nil {
		return err
	}
	password, err := util.CreatePassword(ctx, "OpenshiftLocal-OCP")
	if err != nil {
		return err
	}
	_, err =
		setup.SwapKeys(ctx, &instance.PublicIp,
			r.bootingPrivateKeyFilePath, &privateKey.PublicKeyOpenssh)
	if err != nil {
		return err
	}
	kubeconfig, _, err := setup.Setup(ctx,
		&instance.PublicIp, &privateKey.PrivateKeyOpenssh,
		setup.Data{
			PrivateIP:             &instance.PrivateIp,
			PublicIP:              &instance.PublicIp,
			Password:              &password.Result,
			OCPPullSecretFilePath: r.ocpPullSecretFilePath,
		})
	if err != nil {
		return err
	}
	ctx.Export(providerAPI.Kubeconfig, kubeconfig)
	ctx.Export(providerAPI.OutputKey, privateKey.PrivateKeyPem)
	ctx.Export(providerAPI.OutputHost, instance.PublicIp)
	ctx.Export(providerAPI.OutputUsername, pulumi.String(bundle.ImageUsername))
	ctx.Export(providerAPI.OutputPassword, password.Result)
	return nil
}

func securityGroupsIds(ctx *pulumi.Context) (pulumi.StringArrayInput, error) {
	ingressRules := []sg.IngressRule{
		{Description: "SSH", FromPort: 22, ToPort: 22, Protocol: "tcp"},
		{Description: "Cluster API", FromPort: 443, ToPort: 443, Protocol: "tcp"},
		{Description: "Console", FromPort: 6443, ToPort: 6443, Protocol: "tcp"}}
	sgID, err := sg.Create(ctx, ingressRules, "OpenshiftLocal-OCP", "OpenshiftLocal OCP ingress rules")
	if err != nil {
		return nil, err
	}
	return pulumi.StringArray{sgID}, nil
}

func createKey(ctx *pulumi.Context) (*tls.PrivateKey, *ec2.KeyPair, error) {
	pk, err := tls.NewPrivateKey(
		ctx,
		"OpenshiftLocal-OCP",
		&tls.PrivateKeyArgs{
			Algorithm: pulumi.String("RSA"),
			RsaBits:   pulumi.Int(4096),
		})
	if err != nil {
		return nil, nil, err
	}
	kp, err := ec2.NewKeyPair(ctx,
		"OpenshiftLocal-OCP",
		&ec2.KeyPairArgs{
			PublicKey: pk.PublicKeyOpenssh,
			Tags:      context.GetTags()})
	if err != nil {
		return nil, nil, err
	}
	return pk, kp, nil
}

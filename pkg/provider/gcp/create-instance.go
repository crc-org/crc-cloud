package gcp

import (
	"fmt"
	"strconv"

	"github.com/crc/crc-cloud/pkg/bundle"
	"github.com/crc/crc-cloud/pkg/bundle/setup"
	providerAPI "github.com/crc/crc-cloud/pkg/manager/provider/api"
	"github.com/crc/crc-cloud/pkg/provider/constants"
	"github.com/crc/crc-cloud/pkg/util"
	crctls "github.com/crc/crc-cloud/pkg/util/tls"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/compute"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type createRequest struct {
	projectName               string
	imageID                   string
	instanceType              string
	diskSize                  int
	bootingPrivateKeyFilePath string
	ocpPullSecretFilePath     string
}

func fillCreateRequest(projectName, bootingPrivateKeyFilePath, ocpPullSecretFilePath string,
	args map[string]string) (*createRequest, error) {
	imageIDValue, ok := args[imageID]
	if !ok {
		return nil, fmt.Errorf("imageID not found")
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
		imageID:                   imageIDValue,
		instanceType:              it,
		diskSize:                  ds,
		bootingPrivateKeyFilePath: bootingPrivateKeyFilePath,
		ocpPullSecretFilePath:     ocpPullSecretFilePath}, nil
}

func (r createRequest) runFunc(ctx *pulumi.Context) error {
	privateKey, err := crctls.CreateKey(ctx)
	if err != nil {
		return err
	}
	// Create a new network for the virtual machine.
	network, err := compute.NewNetwork(ctx, r.projectName, &compute.NetworkArgs{
		AutoCreateSubnetworks: pulumi.Bool(false),
	})
	if err != nil {
		return err
	}

	// Create a subnet on the network.
	subnet, err := compute.NewSubnetwork(ctx, r.projectName, &compute.SubnetworkArgs{
		IpCidrRange: pulumi.String("10.0.1.0/24"),
		Network:     network.ID(),
	})
	if err != nil {
		return err
	}

	// Create a firewall allowing inbound access over ports 80 (for HTTP) and 22 (for SSH).
	firewall, err := compute.NewFirewall(ctx, r.projectName, &compute.FirewallArgs{
		Network: network.SelfLink,
		Allows: compute.FirewallAllowArray{
			compute.FirewallAllowArgs{
				Protocol: pulumi.String("tcp"),
				Ports: pulumi.ToStringArray([]string{
					strconv.Itoa(constants.HTTPPort),
					strconv.Itoa(constants.HTTPSPort),
					strconv.Itoa(constants.APIPort),
					strconv.Itoa(constants.SSHPort),
				}),
			},
		},
		Direction: pulumi.String("INGRESS"),
		SourceRanges: pulumi.ToStringArray([]string{
			"0.0.0.0/0",
		}),
		TargetTags: pulumi.ToStringArray([]string{r.projectName}),
	})
	if err != nil {
		return err
	}

	args := compute.InstanceArgs{
		MachineType: pulumi.String(r.instanceType),
		BootDisk: compute.InstanceBootDiskArgs{
			InitializeParams: compute.InstanceBootDiskInitializeParamsArgs{
				Image: pulumi.String(r.imageID),
			},
		},
		NetworkInterfaces: compute.InstanceNetworkInterfaceArray{
			compute.InstanceNetworkInterfaceArgs{
				Network:    network.ID(),
				Subnetwork: subnet.ID(),
				AccessConfigs: compute.InstanceNetworkInterfaceAccessConfigArray{
					compute.InstanceNetworkInterfaceAccessConfigArgs{
						// NatIp:       nil,
						// NetworkTier: nil,
					},
				},
			},
		},
		ServiceAccount: compute.InstanceServiceAccountArgs{
			Scopes: pulumi.ToStringArray([]string{
				"https://www.googleapis.com/auth/cloud-platform",
			}),
		},
		AllowStoppingForUpdate: pulumi.Bool(true),
		Tags:                   pulumi.ToStringArray([]string{r.projectName}),
	}

	instance, err := compute.NewInstance(ctx, r.projectName, &args, pulumi.DependsOn([]pulumi.Resource{firewall}))
	if err != nil {
		return err
	}

	internalIP := instance.NetworkInterfaces.Index(pulumi.Int(0)).NetworkIp().Elem()
	publicIP := instance.NetworkInterfaces.Index(pulumi.Int(0)).AccessConfigs().Index(pulumi.Int(0)).NatIp().Elem()

	password, err := util.CreatePassword(ctx, "OpenshiftLocal-OCP")
	if err != nil {
		return err
	}
	_, err = setup.SwapKeys(ctx, &publicIP,
		r.bootingPrivateKeyFilePath, &privateKey.PublicKeyOpenssh)
	if err != nil {
		return err
	}
	kubeconfig, _, err := setup.Setup(ctx,
		&publicIP, &privateKey.PrivateKeyOpenssh,
		setup.Data{
			PrivateIP:             &internalIP,
			PublicIP:              &publicIP,
			Password:              &password.Result,
			OCPPullSecretFilePath: r.ocpPullSecretFilePath,
		})
	if err != nil {
		return err
	}
	ctx.Export(providerAPI.Kubeconfig, kubeconfig)
	ctx.Export(providerAPI.OutputKey, privateKey.PrivateKeyPem)
	ctx.Export(providerAPI.OutputHost, publicIP)
	ctx.Export(providerAPI.OutputUsername, pulumi.String(bundle.ImageUsername))
	ctx.Export(providerAPI.OutputPassword, password.Result)
	return nil
}

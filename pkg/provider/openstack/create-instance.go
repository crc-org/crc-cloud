package openstack

import (
	"fmt"
	"strconv"

	"github.com/crc/crc-cloud/pkg/bundle"
	"github.com/crc/crc-cloud/pkg/bundle/setup"
	providerAPI "github.com/crc/crc-cloud/pkg/manager/provider/api"
	"github.com/crc/crc-cloud/pkg/provider/constants"
	"github.com/crc/crc-cloud/pkg/util"
	crctls "github.com/crc/crc-cloud/pkg/util/tls"
	"github.com/pulumi/pulumi-openstack/sdk/v3/go/openstack/blockstorage"
	"github.com/pulumi/pulumi-openstack/sdk/v3/go/openstack/compute"
	"github.com/pulumi/pulumi-openstack/sdk/v3/go/openstack/images"
	"github.com/pulumi/pulumi-openstack/sdk/v3/go/openstack/networking"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type createRequest struct {
	projectName               string
	imageID                   string
	instanceType              string
	diskSize                  int
	bootingPrivateKeyFilePath string
	ocpPullSecretFilePath     string
	networkName               string
}

func fillCreateRequest(projectName, bootingPrivateKeyFilePath, ocpPullSecretFilePath string,
	args map[string]string) (*createRequest, error) {
	imageIDValue, ok := args[imageID]
	if !ok {
		return nil, fmt.Errorf("imageID not found")
	}
	networkNameValue, ok := args[networkName]
	if !ok {
		return nil, fmt.Errorf("openstack-network-name not found")
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
		ocpPullSecretFilePath:     ocpPullSecretFilePath,
		networkName:               networkNameValue}, nil
}

func (r createRequest) runFunc(ctx *pulumi.Context) error {
	privateKey, err := crctls.CreateKey(ctx)
	if err != nil {
		return err
	}

	// Check if provided network exists
	_, err = networking.LookupNetwork(ctx, &networking.LookupNetworkArgs{
		Name: pulumi.StringRef(r.networkName)}, nil)
	if err != nil {
		return fmt.Errorf("provided %s network does not exist. Err: %v", r.networkName, err)
	}

	// Create the security group for the instance
	secGroup, err := networking.NewSecGroup(ctx, r.projectName, &networking.SecGroupArgs{
		Description: pulumi.String("My neutron security group"),
	})
	if err != nil {
		return err
	}

	// Define an array of the ports to allow
	ports := []int{constants.HTTPPort, constants.HTTPSPort, constants.APIPort, constants.SSHPort}

	// Iterate over the ports and create a security group rule for each
	for _, port := range ports {
		_, err := networking.NewSecGroupRule(ctx, fmt.Sprintf("secgroupRule%d", port), &networking.SecGroupRuleArgs{
			Direction:       pulumi.String("ingress"),
			Ethertype:       pulumi.String("IPv4"),
			SecurityGroupId: secGroup.ID(),
			PortRangeMin:    pulumi.Int(port),
			PortRangeMax:    pulumi.Int(port),
			Protocol:        pulumi.String("tcp"),
			RemoteIpPrefix:  pulumi.String("0.0.0.0/0"),
		})
		if err != nil {
			return err
		}
	}

	imageRef, err := images.LookupImage(ctx, &images.LookupImageArgs{Name: pulumi.StringRef(r.imageID)}, nil)
	if err != nil {
		return err
	}

	vol, err := blockstorage.NewVolume(ctx, "myVolume", &blockstorage.VolumeArgs{
		Size:             pulumi.Int(r.diskSize), // Size of the volume in GB
		AvailabilityZone: pulumi.String("nova"),  // The Availability Zone in which to create the volume
		ImageId:          pulumi.String(imageRef.Id),
	})
	if err != nil {
		return err
	}

	args := compute.InstanceArgs{
		FlavorName: pulumi.String(r.instanceType),
		//ImageName:  pulumi.String(r.imageID),
		SecurityGroups: pulumi.StringArray{
			secGroup.Name,
		},
		Networks: compute.InstanceNetworkArray{
			&compute.InstanceNetworkArgs{
				Name: pulumi.String(r.networkName),
			},
		},
		BlockDevices: compute.InstanceBlockDeviceArray{
			&compute.InstanceBlockDeviceArgs{
				Uuid:                vol.ID(),
				SourceType:          pulumi.String("volume"),
				DestinationType:     pulumi.String("volume"),
				BootIndex:           pulumi.Int(0),
				DeleteOnTermination: pulumi.Bool(true),
			},
		},
	}

	instance, err := compute.NewInstance(ctx, r.projectName, &args, pulumi.DependsOn([]pulumi.Resource{secGroup}))
	if err != nil {
		return err
	}

	internalIP := instance.AccessIpV4
	publicIP := internalIP

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

package upi

import (
	"fmt"
	"net"
	"os"

	"github.com/crc/crc-cloud/pkg/bundle"
	"github.com/crc/crc-cloud/pkg/bundle/setup"
	providerAPI "github.com/crc/crc-cloud/pkg/manager/provider/api"
	"github.com/crc/crc-cloud/pkg/util"
	crctls "github.com/crc/crc-cloud/pkg/util/tls"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type createRequest struct {
	projectName           string
	hostIP                string
	sshKeyPath            string
	ocpPullSecretFilePath string
}

func fillCreateRequest(projectName, _ /* bootingPrivateKeyFilePath */, ocpPullSecretFilePath string,
	args map[string]string) (*createRequest, error) {
	hostIPValue, ok := args[hostIP]
	if !ok || hostIPValue == "" {
		return nil, fmt.Errorf("host-ip not found")
	}
	if net.ParseIP(hostIPValue) == nil {
		return nil, fmt.Errorf("invalid host-ip: %q", hostIPValue)
	}
	sshKeyPathValue, ok := args[sshKeyPath]
	if !ok || sshKeyPathValue == "" {
		return nil, fmt.Errorf("ssh-key-path not found")
	}

	// Verify SSH key file exists and is accessible
	if _, err := os.Stat(sshKeyPathValue); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("SSH key file not found at path: %s", sshKeyPathValue)
		}
		return nil, fmt.Errorf("unable to access SSH key at %s: %w", sshKeyPathValue, err)
	}

	return &createRequest{
		projectName:           projectName,
		hostIP:                hostIPValue,
		sshKeyPath:            sshKeyPathValue,
		ocpPullSecretFilePath: ocpPullSecretFilePath,
	}, nil
}

func (r createRequest) runFunc(ctx *pulumi.Context) error {
	// Create a new SSH key pair for the cluster
	privateKey, err := crctls.CreateKey(ctx)
	if err != nil {
		return err
	}

	// Convert IP string to pulumi.StringOutput
	publicIP := pulumi.String(r.hostIP).ToStringOutput()

	// For UPI provider, we assume internal and public IP are the same
	// since the user is managing the VM themselves
	internalIP := publicIP

	// Generate password for the cluster
	password, err := util.CreatePassword(ctx, "OpenshiftLocal-OCP")
	if err != nil {
		return err
	}

	// Swap keys: use the user-provided key to connect, then add the new key
	_, err = setup.SwapKeys(ctx, &publicIP,
		r.sshKeyPath, &privateKey.PublicKeyOpenssh)
	if err != nil {
		return err
	}

	// Run the cluster setup with the new key
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

	// Export outputs
	ctx.Export(providerAPI.Kubeconfig, kubeconfig)
	ctx.Export(providerAPI.OutputKey, privateKey.PrivateKeyPem)
	ctx.Export(providerAPI.OutputHost, publicIP)
	ctx.Export(providerAPI.OutputUsername, pulumi.String(bundle.ImageUsername))
	ctx.Export(providerAPI.OutputPassword, password.Result)

	return nil
}

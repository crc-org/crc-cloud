package upi

import (
	"fmt"

	providerAPI "github.com/crc/crc-cloud/pkg/manager/provider/api"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Provider struct{}

func GetProvider() *Provider {
	return &Provider{}
}

func (a *Provider) GetPlugin() *providerAPI.PluginInfo {
	// UPI provider doesn't need a Pulumi plugin since it doesn't create infrastructure
	return nil
}

func (a *Provider) ImportImageRunFunc(_, _, _ string) (pulumi.RunFunc, error) {
	// UPI provider operates on an existing VM and does not import images
	return nil, fmt.Errorf("image import is not supported by the UPI provider")
}

func (a *Provider) CreateParams() map[string]string {
	return map[string]string{
		hostIP:     hostIPDesc,
		sshKeyPath: sshKeyPathDesc,
	}
}

func (a *Provider) CreateParamsMandatory() []string {
	return []string{hostIP, sshKeyPath}
}

func (a *Provider) CreateRunFunc(projectName, bootingPrivateKeyFilePath, ocpPullSecretFilePath string,
	args map[string]string) (pulumi.RunFunc, error) {
	r, err := fillCreateRequest(projectName, bootingPrivateKeyFilePath, ocpPullSecretFilePath, args)
	if err != nil {
		return nil, err
	}
	return r.runFunc, nil
}

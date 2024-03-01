package openstack

import (
	providerAPI "github.com/crc/crc-cloud/pkg/manager/provider/api"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Provider struct{}

func GetProvider() *Provider {
	return &Provider{}
}

func (a *Provider) GetPlugin() *providerAPI.PluginInfo {
	return &providerAPI.PluginInfo{
		Name:    "openstack",
		Version: "v3.15.1"}
}

func (a *Provider) ImportImageRunFunc(_, _, _ string) (pulumi.RunFunc, error) {
	return nil, nil
}

func (a *Provider) CreateParams() map[string]string {
	return map[string]string{
		imageID:      imageIDDesc,
		instanceType: instanceTypeDesc,
		diskSize:     diskSizeDesc,
		networkName:  networkNameDesc,
	}
}

func (a *Provider) CreateParamsMandatory() []string {
	return []string{imageID, networkName}
}

func (a *Provider) CreateRunFunc(projectName, bootingPrivateKeyFilePath, ocpPullSecretFilePath string,
	args map[string]string) (pulumi.RunFunc, error) {
	r, err := fillCreateRequest(projectName, bootingPrivateKeyFilePath, ocpPullSecretFilePath, args)
	if err != nil {
		return nil, err
	}
	return r.runFunc, nil
}

package gcp

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
		Name:    "gcp",
		Version: "v6.65.0"}
}

func (a *Provider) ImportImageRunFunc(projectName, bundleDownloadURL, shasumfileDownloadURL string) (pulumi.RunFunc, error) {
	return nil, nil
}

func (a *Provider) CreateParams() map[string]string {
	return map[string]string{
		imageID:      imageIDDesc,
		instanceType: instanceTypeDesc,
		diskSize:     diskSizeDesc,
	}
}

func (a *Provider) CreateParamsMandatory() []string {
	return []string{imageID}
}

func (a *Provider) CreateRunFunc(projectName, bootingPrivateKeyFilePath, ocpPullSecretFilePath string,
	args map[string]string) (pulumi.RunFunc, error) {
	r, err := fillCreateRequest(projectName, bootingPrivateKeyFilePath, ocpPullSecretFilePath, args)
	if err != nil {
		return nil, err
	}
	return r.runFunc, nil
}

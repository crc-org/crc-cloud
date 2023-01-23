package aws

import (
	providerAPI "github.com/crc/crc-cloud/pkg/manager/provider/api"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type AWSProvider struct{}

func GetProvider() *AWSProvider {
	return &AWSProvider{}
}

func (a *AWSProvider) GetPlugin() *providerAPI.PluginInfo {
	return &providerAPI.PluginInfo{
		Name:    "aws",
		Version: "v5.27.0"}
}

func (a *AWSProvider) ImportImageRunFunc(projectName, bundleDownloadURL, shasumfileDownloadURL string) (pulumi.RunFunc, error) {
	r, err := fillImportRequest(projectName, bundleDownloadURL, shasumfileDownloadURL)
	if err != nil {
		return nil, err
	}
	return (pulumi.RunFunc)(r.runFunc), nil
}

func (a *AWSProvider) CreateParams() map[string]string {
	return map[string]string{
		amiID: amiIDDesc,
	}
}

func (a *AWSProvider) CreateParamsMandatory() []string {
	return []string{amiID}
}

func (a *AWSProvider) CreateRunFunc(projectName, bootingPrivateKeyFilePath, ocpPullSecretFilePath string,
	args map[string]string) (pulumi.RunFunc, error) {
	r, err := fillCreateRequest(projectName, bootingPrivateKeyFilePath, ocpPullSecretFilePath, args)
	if err != nil {
		return nil, err
	}
	return (pulumi.RunFunc)(r.runFunc), nil
}

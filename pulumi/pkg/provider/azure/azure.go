package azure

import (
	providerAPI "github.com/crc/crc-cloud/pkg/manager/provider/api"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	azureNativePluginName    = "azure-native"
	azureNativePluginVersion = "v1.92.0"
)

type AzureProvider struct{}

func GetProvider() *AzureProvider {
	return &AzureProvider{}
}

func (a *AzureProvider) GetPlugin() *providerAPI.PluginInfo {
	return &providerAPI.PluginInfo{
		Name:    azureNativePluginName,
		Version: azureNativePluginVersion}
}

func (a *AzureProvider) ImportImageRunFunc(projectName, bundleDownloadURL, shasumfileDownloadURL string) (pulumi.RunFunc, error) {
	r, err := fillImportRequest(projectName, bundleDownloadURL, shasumfileDownloadURL)
	if err != nil {
		return nil, err
	}
	return (pulumi.RunFunc)(r.runFunc), nil
}

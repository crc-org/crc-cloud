package api

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Struct with information about
// plugin to manage provider resources
type PluginInfo struct {
	Name    string
	Version string
}

// Struct holding the information for
// the pulumi stack
type Stack struct {
	ProjectName string
	StackName   string
	BackedURL   string
	DeployFunc  pulumi.RunFunc
	Plugin      PluginInfo
}

type Provider interface {
	// Plugin information, required to dynamically install the plugin
	//for the specific provider
	GetPlugin() *PluginInfo

	// Manage all the image import process for the specific provider
	ImportImageRunFunc(projectName,
		bundleDownloadURL, shasumfileDownloadURL string) (pulumi.RunFunc, error)

	// Set of params tied to provider to customize the create operation
	CreateParams() map[string]string
	// Subset of create params which are mandatory
	CreateParamsMandatory() []string
	// Creates all resources for the specific provider required on the create operation
	CreateRunFunc(projectName,
		bootingPrivateKeyFilePath, ocpPullSecretFilePath string,
		args map[string]string) (pulumi.RunFunc, error)
}

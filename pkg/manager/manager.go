package manager

import (
	providerAPI "github.com/crc/crc-cloud/pkg/manager/provider/api"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"golang.org/x/exp/maps"
)

const (
	stackCreate      string = "crcCloud-Create"
	stackImportImage string = "crcCloud-ImageImport"
)

func CreateParams() (params map[string]string) {
	params = map[string]string{}
	for _, p := range getSupportedProviders() {
		maps.Copy(params, p.CreateParams())
	}
	return
}

func Import(projectName, backerURL, outputFoler string,
	bundleDownloadURL, shasumfileDownloadURL string, provider Provider) error {
	// Pick the import function according to the provider
	p, err := getProvider(provider)
	if err != nil {
		return err
	}
	importRunFunc, err :=
		p.ImportImageRunFunc(projectName, bundleDownloadURL, shasumfileDownloadURL)
	if err != nil {
		return err
	}
	// Create a stack based on the import function and create it
	stack := providerAPI.Stack{
		ProjectName: projectName,
		StackName:   stackImportImage,
		BackedURL:   backerURL,
		DeployFunc:  importRunFunc,
		Plugin:      *p.GetPlugin()}
	stackResult, err := upStack(stack)
	if err != nil {
		return err
	}
	err = manageImageImportResults(stackResult, outputFoler)
	if err != nil {
		return nil
	}

	// Current exec create temporary resources to enable the import
	// we delete it as they are only temporary
	return destroyStack(stack)
}

func manageImageImportResults(stackResult auto.UpResult, destinationFolder string) error {
	if err := writeOutputs(stackResult, destinationFolder, map[string]string{
		providerAPI.OutputBootKey: "id_ecdsa",
		providerAPI.OutputImageID: "image-id",
	}); err != nil {
		return err
	}
	return nil
}

func Create(projectName, backerURL, outputFoler string,
	provider Provider, providerArgs map[string]string,
	ocpPullSecretFilePath, bootKeyFilePath string) error {
	// this will return a provider which implements the api.Provider interface
	p, err := getProvider(provider)
	if err != nil {
		return err
	}
	// TODO think best option to pass params to provider
	// may serialize all params and let provider validate and pick the required
	// as a provider the params are specs and manager requires to know about them?
	err = validateParams(providerArgs, p.CreateParamsMandatory())
	if err != nil {
		return err
	}

	createFunc, err :=
		p.CreateRunFunc(projectName, bootKeyFilePath, ocpPullSecretFilePath,
			providerArgs)
	if err != nil {
		return err
	}

	createStack := providerAPI.Stack{
		ProjectName: projectName,
		StackName:   stackCreate,
		BackedURL:   backerURL,
		DeployFunc:  createFunc,
		Plugin:      *p.GetPlugin()}
	stackResult, err := upStack(createStack)
	if err != nil {
		return err
	}
	return manageCreateResults(stackResult, outputFoler)
}

func Destroy(projectName, backedURL string, provider Provider) error {
	// this will return a provider which implements the api.Provider interface
	p, err := getProvider(provider)
	if err != nil {
		return err
	}
	stack := providerAPI.Stack{
		ProjectName: projectName,
		StackName:   stackCreate,
		BackedURL:   backedURL,
		Plugin:      *p.GetPlugin()}
	return destroyStack(stack)
}

func manageCreateResults(stackResult auto.UpResult, destinationFolder string) error {
	if err := writeOutputs(stackResult, destinationFolder, map[string]string{
		providerAPI.OutputKey:      "id_rsa",
		providerAPI.OutputHost:     "host",
		providerAPI.OutputUsername: "username",
		providerAPI.OutputPassword: "password",
	}); err != nil {
		return err
	}
	return nil
}

package azure

import (
	"github.com/pulumi/pulumi-azure-native-sdk/resources"
	"github.com/pulumi/pulumi-azure-native-sdk/storage"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
// defaultVolumeSize int = 100
)

type importRequest struct {
	projectName           string
	bundleDownloadURL     string
	shasumfileDownloadURL string
}

func fillImportRequest(projectName, bundleDownloadURL, shasumfileDownloadURL string) (*importRequest, error) {
	return &importRequest{
		projectName:           projectName,
		bundleDownloadURL:     bundleDownloadURL,
		shasumfileDownloadURL: shasumfileDownloadURL,
	}, nil
}

func (r importRequest) runFunc(ctx *pulumi.Context) error {
	_, err := createTempBlobStorageContainer(ctx)
	if err != nil {
		return err
	}
	// bundleExtractAssets, bootkey, err := bundleExtract.Extract(ctx, r.bundleDownloadURL, r.shasumfileDownloadURL)
	// if err != nil {
	// 	return err
	// }

	// ctx.Export(providerAPI.OutputBootKey, *bootkey)
	// ctx.Export(providerAPI.OutputImageID, ami.ID())
	return nil
}

// This function creates a temporary bucket to upload the disk image to be imported
// As a temporary bucket bucket name is not set so a s
func createTempBlobStorageContainer(ctx *pulumi.Context) (*storage.BlobContainer, error) {
	rg, err := resources.NewResourceGroup(ctx,
		"crcCloudImporterRG",
		&resources.ResourceGroupArgs{
			Location: pulumi.String("West Europe"),
		})
	if err != nil {
		return nil, err
	}
	sa, err := storage.NewStorageAccount(ctx,
		"crcCloudImporterSA",
		&storage.StorageAccountArgs{
			ResourceGroupName: rg.Name,
			Location:          rg.Location,
			AccessTier:        storage.AccessTierHot})
	if err != nil {
		return nil, err
	}
	return storage.NewBlobContainer(ctx,
		"crcCloudImporterBC",
		&storage.BlobContainerArgs{
			AccountName:       sa.Name,
			ResourceGroupName: rg.Name,
			PublicAccess:      storage.PublicAccessNone})
}

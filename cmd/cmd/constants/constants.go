package constants

const (
	ProjectName               string = "project-name"
	ProjectNameDesc           string = "project name to identify the instance of the stack"
	BackedURL                 string = "backed-url"
	BackedURLDesc             string = "backed for stack state. Can be a local path with format file:///path/subpath or s3 s3://existing-bucket"
	OutputFolder              string = "output"
	OutputFolderDesc          string = "path to export assets"
	Provider                  string = "provider"
	ProviderDesc              string = "target cloud provider"
	OcpPullSecretFilePath     string = "pullsecret-filepath"
	OcpPullSecretFilePathDesc string = "path for pullsecret file"
	KeyFilePath               string = "key-filepath"
	KeyFilePathDesc           string = "path to init key obtained when importing the image"

	BundleDownloadURL         string = "bundle-url"
	BundleDownloadURLDesc     string = "custom url to download the bundle artifact"
	ShasumfileDownloadURL     string = "bundle-shasumfile-url"
	ShasumfileDownloadURLDesc string = "custom url to download the shasum file to verify the bundle artifact"
)

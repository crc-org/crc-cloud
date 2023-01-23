package cmd

const (
	projectName               string = "project-name"
	projectNameDesc           string = "project name to identify the instance of the stack"
	backedURL                 string = "backed-url"
	backedURLDesc             string = "backed for stack state. Can be a local path with format file:///path/subpath or s3 s3://existing-bucket"
	outputFolder              string = "output"
	outputFolderDesc          string = "path to export assets"
	provider                  string = "provider"
	providerDesc              string = "target cloud provider"
	ocpPullSecretFilePath     string = "pullsecret-filepath"
	ocpPullSecretFilePathDesc string = "path for pullsecret file"
	keyFilePath               string = "key-filepath"
	keyFilePathDesc           string = "path to init key obtained when importing the image"

	bundleDownloadURL         string = "bundle-url"
	bundleDownloadURLDesc     string = "custom url to download the bundle artifact"
	shasumfileDownloadURL     string = "bundle-shasumfile-url"
	shasumfileDownloadURLDesc string = "custom url to download the shasum file to verify the bundle artifact"
)

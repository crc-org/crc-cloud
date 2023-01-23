package bundle

import (
	"fmt"
	"regexp"
)

const (
	ImageUsername                   string = "core"
	ImageSSHPort                    int    = 22
	ImageInternalKubeconfigFilepath string = "/opt/kubeconfig"

	bundleDescription    = "openshift-local"
	bundleVersionUnknown = "unknown"
)

var (
	bundleVersionRegex string = "\\d.\\d+.\\d+"
	bundleNameRegex    string = fmt.Sprintf(
		"crc_libvirt_%s_amd64.crcbundle", bundleVersionRegex)
)

// Bundle name format contains the version number we are managing
// this function will return a description having that versioning info in it
// if bundle name does not mach the default format version will be reported as unknown
func GetDescription(bundleURL string) (*string, error) {
	var bundleVersion string
	rn, err := regexp.Compile(bundleNameRegex)
	if err != nil {
		return nil, err
	}
	bundleName := rn.FindString(bundleURL)
	if len(bundleName) > 0 {
		rv, err := regexp.Compile(bundleVersionRegex)
		if err != nil {
			return nil, err
		}
		bundleVersion = rv.FindString(bundleName)
		if len(bundleVersion) == 0 {
			bundleVersion = bundleVersionUnknown
		}
	}
	description :=
		fmt.Sprintf("%s-%s", bundleDescription, bundleVersion)
	return &description, nil
}

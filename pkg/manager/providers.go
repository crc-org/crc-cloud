package manager

import (
	"fmt"

	providerAPI "github.com/crc/crc-cloud/pkg/manager/provider/api"
	"github.com/crc/crc-cloud/pkg/provider/aws"
	"github.com/crc/crc-cloud/pkg/provider/gcp"
	"github.com/crc/crc-cloud/pkg/provider/openstack"
)

type Provider string

const (
	AWS Provider = "aws"
	GCP Provider = "gcp"
	OSP Provider = "openstack"
	AZ  Provider = "azure"
)

func getProvider(provider Provider) (providerAPI.Provider, error) {
	switch provider {
	case AWS, AZ:
		return aws.GetProvider(), nil
	case GCP:
		return gcp.GetProvider(), nil
	case OSP:
		return openstack.GetProvider(), nil
	}
	return nil, fmt.Errorf("%s: provider not supported", provider)
}

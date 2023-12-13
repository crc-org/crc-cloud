package manager

import (
	"fmt"

	providerAPI "github.com/crc/crc-cloud/pkg/manager/provider/api"
	"github.com/crc/crc-cloud/pkg/provider/aws"
)

type Provider string

const (
	AWS Provider = "aws"
	GCP Provider = "gcp"
	AZ  Provider = "azure"
)

func getProvider(provider Provider) (providerAPI.Provider, error) {
	switch provider {
	case AWS, AZ:
		return aws.GetProvider(), nil
	}
	return nil, fmt.Errorf("%s: provider not supported", provider)
}

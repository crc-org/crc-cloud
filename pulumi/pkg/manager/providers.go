package manager

import (
	"fmt"

	providerAPI "github.com/crc/crc-cloud/pkg/manager/provider/api"
	"github.com/crc/crc-cloud/pkg/provider/aws"
)

func getProvider(provider string) (providerAPI.Provider, error) {
	switch provider {
	case "aws":
		return aws.GetProvider(), nil
	}
	return nil, fmt.Errorf("provider not supported")
}

func getSupportedProviders() (sp []providerAPI.Provider) {
	sp = append(sp, aws.GetProvider())
	return
}

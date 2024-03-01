package create

import (
	"fmt"
	"os"

	"github.com/crc/crc-cloud/cmd/cmd/constants"
	"github.com/crc/crc-cloud/pkg/manager"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	gcpProviderName        string = "gcp"
	gcpProviderDescription string = "create crc cloud instance on GCP"
)

func getGCPProviderCmd() *cobra.Command {
	gcpProviderCmd := &cobra.Command{
		Use:   gcpProviderName,
		Short: gcpProviderDescription,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := viper.BindPFlags(cmd.Flags()); err != nil {
				return err
			}
			// Provider dependent params
			providerParams := make(map[string]string)
			for name := range manager.CreateParams(manager.GCP) {
				if viper.IsSet(name) {
					providerParams[name] = viper.GetString(name)
				}
			}
			fmt.Printf("providerParams: %v", providerParams)
			if err := manager.Create(
				viper.GetString(constants.ProjectName),
				viper.GetString(constants.BackedURL),
				viper.GetString(constants.OutputFolder),
				manager.GCP,
				providerParams,
				viper.GetString(constants.OcpPullSecretFilePath),
				viper.GetString(constants.KeyFilePath),
				viper.GetStringMapString(constants.Tags)); err != nil {
				fmt.Printf("error creating the cluster with %s provider: %s\n", manager.GCP, err)
				os.Exit(1)
			}
			return nil
		},
	}

	flagSet := pflag.NewFlagSet(gcpProviderName, pflag.ExitOnError)
	// Provider dependent params
	for name, description := range manager.CreateParams(manager.GCP) {
		flagSet.StringP(name, "", "", description)
	}

	gcpProviderCmd.Flags().AddFlagSet(flagSet)
	return gcpProviderCmd
}

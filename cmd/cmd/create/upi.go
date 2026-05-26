package create

import (
	"fmt"

	"github.com/crc/crc-cloud/cmd/cmd/constants"
	"github.com/crc/crc-cloud/pkg/manager"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	upiProviderName        string = "upi"
	upiProviderDescription string = "create crc cloud instance on User Provided Infrastructure"
)

func getUPIProviderCmd() *cobra.Command {
	upiProviderCmd := &cobra.Command{
		Use:   upiProviderName,
		Short: upiProviderDescription,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := viper.BindPFlags(cmd.Flags()); err != nil {
				return err
			}
			// Provider dependent params
			providerParams := make(map[string]string)
			for name := range manager.CreateParams(manager.UPI) {
				if viper.IsSet(name) {
					providerParams[name] = viper.GetString(name)
				}
			}
			if err := manager.Create(
				viper.GetString(constants.ProjectName),
				viper.GetString(constants.BackedURL),
				viper.GetString(constants.OutputFolder),
				manager.UPI,
				providerParams,
				viper.GetString(constants.OcpPullSecretFilePath),
				viper.GetString(constants.KeyFilePath),
				viper.GetStringMapString(constants.Tags)); err != nil {
				return fmt.Errorf("error creating the cluster with %s provider: %w", manager.UPI, err)
			}
			return nil
		},
	}

	flagSet := pflag.NewFlagSet(upiProviderName, pflag.ExitOnError)
	// Provider dependent params
	for name, description := range manager.CreateParams(manager.UPI) {
		flagSet.StringP(name, "", "", description)
	}

	upiProviderCmd.Flags().AddFlagSet(flagSet)
	return upiProviderCmd
}

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
	openstackProviderName        string = "openstack"
	openstackProviderDescription string = "create crc cloud instance on OpenStack"
)

func getOpenStackProviderCmd() *cobra.Command {
	openStackProviderCmd := &cobra.Command{
		Use:   openstackProviderName,
		Short: openstackProviderDescription,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := viper.BindPFlags(cmd.Flags()); err != nil {
				return err
			}
			// Provider dependent params
			providerParams := make(map[string]string)
			for name := range manager.CreateParams(manager.OSP) {
				if viper.IsSet(name) {
					providerParams[name] = viper.GetString(name)
				}
			}
			fmt.Printf("providerParams: %v", providerParams)
			if err := manager.Create(
				viper.GetString(constants.ProjectName),
				viper.GetString(constants.BackedURL),
				viper.GetString(constants.OutputFolder),
				manager.OSP,
				providerParams,
				viper.GetString(constants.OcpPullSecretFilePath),
				viper.GetString(constants.KeyFilePath),
				viper.GetStringMapString(constants.Tags)); err != nil {
				fmt.Printf("error creating the cluster with %s provider: %s\n", manager.OSP, err)
				os.Exit(1)
			}
			return nil
		},
	}

	flagSet := pflag.NewFlagSet(openstackProviderName, pflag.ExitOnError)
	// Provider dependent params
	for name, description := range manager.CreateParams(manager.OSP) {
		flagSet.StringP(name, "", "", description)
	}

	openStackProviderCmd.Flags().AddFlagSet(flagSet)
	return openStackProviderCmd
}

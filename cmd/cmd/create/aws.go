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
	awsProviderName        string = "aws"
	awsProviderDescription string = "create crc cloud instance on AWS"
)

func getAWSProviderCmd() *cobra.Command {
	awsProviderCmd := &cobra.Command{
		Use:   awsProviderName,
		Short: awsProviderDescription,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := viper.BindPFlags(cmd.Flags()); err != nil {
				return err
			}
			// Provider dependent params
			providerParams := make(map[string]string)
			for name := range manager.CreateParams(manager.AWS) {
				if viper.IsSet(name) {
					providerParams[name] = viper.GetString(name)
				}
			}
			if err := manager.Create(
				viper.GetString(constants.ProjectName),
				viper.GetString(constants.BackedURL),
				viper.GetString(constants.OutputFolder),
				manager.AWS,
				providerParams,
				viper.GetString(constants.OcpPullSecretFilePath),
				viper.GetString(constants.KeyFilePath),
				viper.GetStringMapString(constants.Tags)); err != nil {
				fmt.Printf("error creating the cluster with %s provider: %s\n", manager.AWS, err)
				os.Exit(1)
			}
			return nil
		},
	}

	flagSet := pflag.NewFlagSet(awsProviderName, pflag.ExitOnError)
	// Provider dependent params
	for name, description := range manager.CreateParams(manager.AWS) {
		flagSet.StringP(name, "", "", description)
	}

	awsProviderCmd.Flags().AddFlagSet(flagSet)
	return awsProviderCmd
}

package create

import (
	"github.com/crc/crc-cloud/cmd/cmd/constants"
	"github.com/crc/crc-cloud/pkg/manager"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
)

const (
	awsProviderName        string = "aws"
	awsProviderDescription string = "create crc cloud instance on AWS"
)

func getAWSProviderCmd() *cobra.Command {
	awsProviderCmd := &cobra.Command{
		Use:   awsProviderName,
		Short: awsProviderDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := viper.BindPFlags(cmd.Flags()); err != nil {
				return err
			}
			// Provider dependent params
			providerParams := make(map[string]string)
			for name := range manager.CreateParams() {
				providerParams[name] = viper.GetString(name)
			}
			if err := manager.Create(
				viper.GetString(constants.ProjectName),
				viper.GetString(constants.BackedURL),
				viper.GetString(constants.OutputFolder),
				viper.GetString(constants.Provider),
				providerParams,
				viper.GetString(constants.OcpPullSecretFilePath),
				viper.GetString(constants.KeyFilePath)); err != nil {
				os.Exit(1)
			}
			return nil
		},
	}

	flagSet := pflag.NewFlagSet(awsProviderName, pflag.ExitOnError)
	// Provider dependent params
	for name, description := range manager.CreateParams() {
		flagSet.StringP(name, "", "", description)
	}

	awsProviderCmd.Flags().AddFlagSet(flagSet)
	return awsProviderCmd
}

package create

import (
	"github.com/crc/crc-cloud/cmd/cmd/constants"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	createCmdName        string = "create"
	createCmdDescription string = "create crc cloud instance"
)

func GetCreateCmd() *cobra.Command {
	createCmd := &cobra.Command{
		Use:   createCmdName,
		Short: createCmdDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := viper.BindPFlags(cmd.Flags()); err != nil {
				return err
			}
			return nil
		},
	}
	flagSet := pflag.NewFlagSet(createCmdName, pflag.ExitOnError)
	// Fixed params
	flagSet.StringP(constants.ProjectName, "", "", constants.ProjectNameDesc)
	flagSet.StringP(constants.BackedURL, "", "", constants.BackedURLDesc)
	flagSet.StringP(constants.OutputFolder, "", "", constants.OutputFolderDesc)
	flagSet.StringP(constants.OcpPullSecretFilePath, "", "", constants.OcpPullSecretFilePathDesc)
	flagSet.StringP(constants.KeyFilePath, "", "", constants.KeyFilePathDesc)
	createCmd.PersistentFlags().AddFlagSet(flagSet)

	createCmd.AddCommand(getAWSProviderCmd())

	return createCmd
}

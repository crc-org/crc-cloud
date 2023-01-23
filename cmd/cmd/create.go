package cmd

import (
	"os"

	"github.com/crc/crc-cloud/pkg/manager"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	createCmdName        string = "create"
	createCmdDescription string = "create crc cloud instance"
)

func init() {
	rootCmd.AddCommand(crcCloudCreateCmd)
	flagSet := pflag.NewFlagSet(createCmdName, pflag.ExitOnError)
	// Fixed params
	flagSet.StringP(projectName, "", "", projectNameDesc)
	flagSet.StringP(backedURL, "", "", backedURLDesc)
	flagSet.StringP(outputFolder, "", "", outputFolderDesc)
	flagSet.StringP(provider, "", "", providerDesc)
	flagSet.StringP(ocpPullSecretFilePath, "", "", ocpPullSecretFilePathDesc)
	flagSet.StringP(keyFilePath, "", "", keyFilePathDesc)
	// Provider dependent params
	for name, description := range manager.CreateParams() {
		flagSet.StringP(name, "", "", description)
	}
	crcCloudCreateCmd.Flags().AddFlagSet(flagSet)
}

var crcCloudCreateCmd = &cobra.Command{
	Use:   createCmdName,
	Short: createCmdDescription,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		// Provider dependet params
		providerParams := make(map[string]string)
		for name := range manager.CreateParams() {
			providerParams[name] = viper.GetString(name)
		}
		if err := manager.Create(
			viper.GetString(projectName),
			viper.GetString(backedURL),
			viper.GetString(outputFolder),
			viper.GetString(provider),
			providerParams,
			viper.GetString(ocpPullSecretFilePath),
			viper.GetString(keyFilePath)); err != nil {
			os.Exit(1)
		}
		return nil
	},
}

package cmd

import (
	"github.com/crc/crc-cloud/cmd/cmd/constants"
	"os"

	"github.com/crc/crc-cloud/pkg/manager"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	importImageCmdName        string = "import"
	importImageCmdDescription string = "import crc cloud image"
)

func init() {
	rootCmd.AddCommand(crcCloudImportCmd)
	flagSet := pflag.NewFlagSet(importImageCmdName, pflag.ExitOnError)
	// Fixed params
	flagSet.StringP(constants.ProjectName, "", "", constants.ProjectNameDesc)
	flagSet.StringP(constants.BackedURL, "", "", constants.BackedURLDesc)
	flagSet.StringP(constants.Provider, "", "", constants.ProviderDesc)
	flagSet.StringP(constants.OutputFolder, "", "", constants.OutputFolderDesc)
	flagSet.StringP(constants.BundleDownloadURL, "", "", constants.BundleDownloadURLDesc)
	flagSet.StringP(constants.ShasumfileDownloadURL, "", "", constants.ShasumfileDownloadURLDesc)
	crcCloudImportCmd.Flags().AddFlagSet(flagSet)
}

var crcCloudImportCmd = &cobra.Command{
	Use:   importImageCmdName,
	Short: importImageCmdDescription,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		if err := manager.Import(
			viper.GetString(constants.ProjectName),
			viper.GetString(constants.BackedURL),
			viper.GetString(constants.Provider),
			viper.GetString(constants.OutputFolder),
			viper.GetString(constants.BundleDownloadURL),
			viper.GetString(constants.ShasumfileDownloadURL)); err != nil {
			os.Exit(1)
		}
		return nil
	},
}

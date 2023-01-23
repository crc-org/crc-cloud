package cmd

import (
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
	flagSet.StringP(projectName, "", "", projectNameDesc)
	flagSet.StringP(backedURL, "", "", backedURLDesc)
	flagSet.StringP(provider, "", "", providerDesc)
	flagSet.StringP(outputFolder, "", "", outputFolderDesc)
	flagSet.StringP(bundleDownloadURL, "", "", bundleDownloadURLDesc)
	flagSet.StringP(shasumfileDownloadURL, "", "", shasumfileDownloadURLDesc)
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
			viper.GetString(projectName),
			viper.GetString(backedURL),
			viper.GetString(provider),
			viper.GetString(outputFolder),
			viper.GetString(bundleDownloadURL),
			viper.GetString(shasumfileDownloadURL)); err != nil {
			os.Exit(1)
		}
		return nil
	},
}

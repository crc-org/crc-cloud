package cmd

import (
	"os"

	"github.com/crc/crc-cloud/pkg/manager"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	destroyCmdName        string = "destroy"
	destroyCmdDescription string = "destroy crc cloud instance"
)

func init() {
	rootCmd.AddCommand(crcCloudDestroyCmd)
	flagSet := pflag.NewFlagSet(createCmdName, pflag.ExitOnError)
	flagSet.StringP(projectName, "", "", projectNameDesc)
	flagSet.StringP(backedURL, "", "", backedURLDesc)
	flagSet.StringP(provider, "", "", providerDesc)
	crcCloudDestroyCmd.Flags().AddFlagSet(flagSet)
}

var crcCloudDestroyCmd = &cobra.Command{
	Use:   destroyCmdName,
	Short: destroyCmdDescription,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		if err := manager.Destroy(
			viper.GetString(projectName),
			viper.GetString(backedURL),
			viper.GetString(provider)); err != nil {
			os.Exit(1)
		}
		return nil
	},
}

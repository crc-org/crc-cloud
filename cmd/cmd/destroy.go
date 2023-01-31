package cmd

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
	destroyCmdName        string = "destroy"
	destroyCmdDescription string = "destroy crc cloud instance"
)

func init() {
	rootCmd.AddCommand(crcCloudDestroyCmd)
	flagSet := pflag.NewFlagSet(destroyCmdName, pflag.ExitOnError)
	flagSet.StringP(constants.ProjectName, "", "", constants.ProjectNameDesc)
	flagSet.StringP(constants.BackedURL, "", "", constants.BackedURLDesc)
	flagSet.StringP(constants.Provider, "", "", constants.ProviderDesc)
	crcCloudDestroyCmd.Flags().AddFlagSet(flagSet)
	if err := crcCloudDestroyCmd.MarkFlagRequired("provider"); err != nil {
		fmt.Printf("Error setting provider as required flag: %s", err)
		os.Exit(1)
	}
}

var crcCloudDestroyCmd = &cobra.Command{
	Use:   destroyCmdName,
	Short: destroyCmdDescription,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		if err := manager.Destroy(
			viper.GetString(constants.ProjectName),
			viper.GetString(constants.BackedURL),
			manager.Provider(viper.GetString(constants.Provider))); err != nil {
			fmt.Printf("error destroying the cluster: %s\n", err)
			os.Exit(1)
		}
		return nil
	},
}

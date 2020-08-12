package cmd

import (
	"github.com/beeeeeeenny/zkcmd/common/version"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "Print the version number of zkcmd",
	Run: func(cmd *cobra.Command, args []string) {
		version.ShowVersion()
	},
}

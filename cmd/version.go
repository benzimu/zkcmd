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
	Short:   "Print version information of zkcmd and quit",
	Run: func(cmd *cobra.Command, args []string) {
		version.ShowVersion()
	},
}

package cmd

import (
	"github.com/benzimu/zkcmd/common/version"

	"github.com/spf13/cobra"
)

func newCmdVersion() *cobra.Command {
	return &cobra.Command{
		Use:     "version",
		Aliases: []string{"v"},
		Short:   "Print version information of zkcmd and quit",
		Run: func(cmd *cobra.Command, args []string) {
			version.ShowVersion()
		},
	}
}

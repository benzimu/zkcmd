package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/beeeeeeenny/zkcmd/common/zookeeper"
	"github.com/samuel/go-zookeeper/zk"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(aclCmd)

	aclCmd.AddCommand(aclGetCmd)
	aclGetCmd.Flags().BoolVarP(&isStat, "stat", "s", false, "znode stat info")

	aclCmd.AddCommand(aclSetCmd)
	aclSetCmd.Flags().BoolVarP(&isStat, "stat", "s", false, "znode stat info")
	aclSetCmd.Flags().StringVarP(&dataVersion, "version", "v", "", "znode data version")
}

var aclCmd = &cobra.Command{
	Use:   "acl",
	Short: "Znode ACL command",
}

var aclGetCmd = &cobra.Command{
	Use:   "get [flags] path",
	Short: "Get znode acl",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		newZKClient()

		acls, stat, err := zkcli.GetACL(args[0])
		checkError(err)

		w := tabwriter.NewWriter(os.Stdout, 6, 4, 3, '\t', 0)
		fmt.Fprintf(w, "ChildrenNum:\t%v\t\n", stat.NumChildren)
		fmt.Fprintf(w, "ACL:        \t%v\t\n", zookeeper.FormatACLs(acls))
		w.Flush()

		if isStat {
			outputStat(stat)
		}
	},
}

var aclSetCmd = &cobra.Command{
	Use:   "set [flags] path acl",
	Short: "Set znode acl",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		newZKClient()

		acls, err := zookeeper.ParseACL(args[1])
		checkError(err)

		exist, stat, err := zkcli.Exists(args[0])
		checkError(err)

		if !exist {
			checkError(zk.ErrNoNode)
		}

		version := checkDataVersion(stat.Aversion)
		stat, err = zkcli.SetACL(args[0], acls, version)
		checkError(err)

		if isStat {
			outputStat(stat)
		}
	},
}

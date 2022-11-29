package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/benzimu/zkcmd/common/zookeeper"
	"github.com/go-zookeeper/zk"
	"github.com/spf13/cobra"
)

func newCmdACL() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "acl",
		Short: "Znode ACL command",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			zkcli = newZKClient()
		},
	}

	cmd.AddCommand(newCmdACLGet())
	cmd.AddCommand(newCmdACLSet())

	return cmd
}

func newCmdACLGet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [flags] path",
		Short: "Get znode acl",
		Args:  cobra.ExactArgs(1),
		Run:   cmdRunACLGet,
	}

	cmd.Flags().BoolVarP(&isStat, "stat", "s", false, "znode stat info")

	return cmd
}

func newCmdACLSet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set [flags] path acl",
		Short: "Set znode acl",
		Args:  cobra.ExactArgs(2),
		Run:   cmdRunACLSet,
	}

	cmd.Flags().BoolVarP(&isStat, "stat", "s", false, "znode stat info")
	cmd.Flags().StringVarP(&dataVersion, "version", "v", "", "znode data version")

	return cmd
}

func cmdRunACLGet(cmd *cobra.Command, args []string) {
	acls, stat, err := zkcli.GetACL(args[0])
	checkError(err)

	w := tabwriter.NewWriter(os.Stdout, 6, 4, 3, '\t', 0)
	fmt.Fprintf(w, "ChildrenNum:\t%v\t\n", stat.NumChildren)
	fmt.Fprintf(w, "ACL:        \t%v\t\n", zookeeper.FormatACLs(acls))
	w.Flush()

	if isStat {
		outputStat(stat)
	}
}

func cmdRunACLSet(cmd *cobra.Command, args []string) {
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
}

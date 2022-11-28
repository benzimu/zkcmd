package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"text/tabwriter"
	"time"

	"github.com/benzimu/zkcmd/common/zookeeper"
	"github.com/go-zookeeper/zk"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	setCreate   bool
	force       bool
	dataVersion string
	isStat      bool

	ephemeral bool
	sequence  bool
)

func newCmdZnode() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "znode",
		Short: "Znode command",
	}

	cmd.AddCommand(newCmdZnodeLs())
	cmd.AddCommand(newCmdZnodeLsn())
	cmd.AddCommand(newCmdZnodeGet())
	cmd.AddCommand(newCmdZnodeDelete())
	cmd.AddCommand(newCmdZnodeSet())
	cmd.AddCommand(newCmdZnodeCreate())

	return cmd
}

func newCmdZnodeLs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls [flags] [path]",
		Short: "List znode children, the path default: /",
		Args:  cobra.MinimumNArgs(0),
		Run:   cmdRunZnodeLs,
	}

	cmd.Flags().BoolVarP(&isStat, "stat", "s", false, "znode stat info")

	return cmd
}

func newCmdZnodeLsn() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ll [flags] [path]",
		Short: "List all znode has no children, the path default: /",
		Args:  cobra.MinimumNArgs(0),
		Run:   cmdRunZnodeLsn,
	}

	return cmd
}

func newCmdZnodeGet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [flags] path",
		Short: "Get znode value",
		Args:  cobra.ExactArgs(1),
		Run:   cmdRunZnodeGet,
	}

	cmd.Flags().BoolVarP(&isStat, "stat", "s", false, "znode stat info")

	return cmd
}

func newCmdZnodeSet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set [flags] path data",
		Short: "Update znode value",
		Args:  cobra.ExactArgs(2),
		Run:   cmdRunZnodeSet,
	}

	cmd.Flags().BoolVarP(&setCreate, "create", "c", false, "will create znode if znode does not exist, but does not directly create multi-level znode")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "will force create multi-level znode if znode does not exist")
	cmd.Flags().StringVarP(&dataVersion, "version", "v", "", "znode data version")
	cmd.Flags().BoolVarP(&isStat, "stat", "s", false, "znode stat info")

	return cmd
}

func newCmdZnodeCreate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [flags] path [data] [acl]",
		Short: "Create znode",
		Example: `  zkcmd znode create /test
	  zkcmd znode create -f /test/1/2
	  zkcmd znode create -f /test/1/2 'data'
	  zkcmd znode create -f /test/1/2 'data' world:anyone:cdrwa`,
		Args: cobra.MinimumNArgs(1),
		Run:  cmdRunZnodeCreate,
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "will force create multi-level znode if znode does not exist")
	cmd.Flags().BoolVarP(&ephemeral, "ephemeral", "e", false, "create ephemeral znode")
	cmd.Flags().BoolVarP(&sequence, "sequence", "s", false, "create sequence znode")

	return cmd
}

func newCmdZnodeDelete() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [flags] path",
		Short: "Delete znode",
		Args:  cobra.ExactArgs(1),
		Run:   cmdRunZnodeDelete,
	}

	cmd.Flags().StringVarP(&dataVersion, "version", "v", "", "znode data version")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "will force delete multi-level znode, like: deleteall")

	return cmd
}

func cmdRunZnodeLs(cmd *cobra.Command, args []string) {
	zkcli := newZKClient()

	path := "/"
	if len(args) > 0 {
		path = args[0]
	}

	cs, stat, err := zkcli.Children(path)
	checkError(err)

	sort.Strings(cs)

	w := tabwriter.NewWriter(os.Stdout, 6, 4, 3, ' ', 0)
	fmt.Fprintf(w, "ID\tPath\tChildrenNum\t\n")

	if stat.NumChildren != 0 {
		for i, c := range cs {
			p := filepath.Join(path, c)
			_, stat, err := zkcli.Children(p)
			checkError(err)

			fmt.Fprintf(w, "%v\t%v\t%v\t\n", i+1, p, stat.NumChildren)
		}
		w.Flush()
	}

	if isStat {
		outputStat(stat)
	}
}

func cmdRunZnodeLsn(cmd *cobra.Command, args []string) {
	zkcli := newZKClient()

	path := "/"
	if len(args) > 0 {
		path = args[0]
	}

	_, stat, err := zkcli.Children(path)
	checkError(err)

	w := tabwriter.NewWriter(os.Stdout, 6, 4, 3, ' ', 0)
	fmt.Fprintf(w, "ID\tPath\t\n")

	if stat.NumChildren != 0 {
		ns, err := zkcli.GetZnodes(path)
		checkError(err)
		for i, n := range ns {
			fmt.Fprintf(w, "%v\t%v\t\n", i+1, n)
		}
		w.Flush()
	}
}

func cmdRunZnodeGet(cmd *cobra.Command, args []string) {
	zkcli := newZKClient()

	d, stat, err := zkcli.Get(args[0])
	checkError(err)

	w := tabwriter.NewWriter(os.Stdout, 6, 4, 3, '\t', 0)
	fmt.Fprintf(w, "ChildrenNum:\t%v\t\n", stat.NumChildren)
	fmt.Fprintf(w, "Value:      \t%v\t\n", string(d))
	w.Flush()

	if isStat {
		outputStat(stat)
	}
}

func cmdRunZnodeSet(cmd *cobra.Command, args []string) {
	zkcli := newZKClient()

	path := args[0]
	data := args[1]

	exist, stat, err := zkcli.Exists(path)
	checkError(err)

	if exist {
		version := checkDataVersion(stat.Version)

		_, err = zkcli.Set(path, []byte(data), version)
		checkError(err)
		return
	}

	if force {
		err = zookeeper.ValidatePath(path, false)
		checkError(err)

		err = zkcli.ForceCreate(path, []byte(data), 0, zk.WorldACL(zk.PermAll))
		checkError(err)
		return
	}

	if setCreate {
		_, err = zkcli.DefaultCreate(path, []byte(data))
		checkError(err)
		return
	}

	exist, stat, err = zkcli.Exists(path)
	checkError(err)

	if exist && isStat {
		outputStat(stat)
	}

	checkError(zk.ErrNoNode)
}

func cmdRunZnodeCreate(cmd *cobra.Command, args []string) {
	zkcli := newZKClient()

	path := args[0]

	var data string
	if len(args) == 2 {
		data = args[1]
	}

	// check acl
	var acl string
	if len(args) > 2 {
		data = args[1]
		acl = args[2]
	}

	acls := zk.WorldACL(zk.PermAll)
	var err error
	if acl != "" {
		acls, err = zookeeper.ParseACL(acl)
		checkError(err)
	}

	// check exist
	exist, _, err := zkcli.Exists(path)
	checkError(err)

	if exist {
		checkError(zk.ErrNodeExists)
	}

	// parse flags
	var flags int32
	if !ephemeral && !sequence {
		flags = int32(0)
	}

	if ephemeral && !sequence {
		flags = int32(zk.FlagEphemeral)
	}

	if !ephemeral && sequence {
		flags = int32(zk.FlagSequence)
	}

	if ephemeral && sequence {
		flags = int32(3)
	}

	if force {
		err = zookeeper.ValidatePath(path, false)
		checkError(err)

		err = zkcli.ForceCreate(path, []byte(data), flags, acls)
		checkError(err)
		return
	}

	_, err = zkcli.Create(path, []byte(data), flags, acls)
	checkError(err)
}

func cmdRunZnodeDelete(cmd *cobra.Command, args []string) {
	zkcli := newZKClient()

	exist, stat, err := zkcli.Exists(args[0])
	checkError(err)

	if exist {
		if force {
			err = zookeeper.ValidatePath(args[0], false)
			checkError(err)

			err = zkcli.ForceDelete(args[0])
			checkError(err)

			return
		}

		version := checkDataVersion(stat.Version)

		err = zkcli.Delete(args[0], version)
		checkError(err)

		return
	}

	checkError(zk.ErrNoNode)
}

func outputStat(stat *zk.Stat) {
	w := tabwriter.NewWriter(os.Stdout, 6, 4, 3, ' ', 0)
	fmt.Fprintf(w, "----------\t\n")
	fmt.Fprintf(w, "Czxid\t%#x\t\n", stat.Czxid)
	fmt.Fprintf(w, "Mzxid\t%#x\t\n", stat.Mzxid)
	fmt.Fprintf(w, "Pzxid\t%#x\t\n", stat.Pzxid)
	fmt.Fprintf(w, "Ctime\t%v\t\n", time.Unix(stat.Ctime/1000, 0))
	fmt.Fprintf(w, "Mtime\t%v\t\n", time.Unix(stat.Mtime/1000, 0))
	fmt.Fprintf(w, "DataVersion\t%v\t\n", stat.Version)
	fmt.Fprintf(w, "Cversion\t%v\t\n", stat.Cversion)
	fmt.Fprintf(w, "AclVersion\t%v\t\n", stat.Aversion)
	fmt.Fprintf(w, "EphemeralOwner\t%v\t\n", stat.EphemeralOwner)
	fmt.Fprintf(w, "DataLength\t%v\t\n", stat.DataLength)
	w.Flush()
}

func checkDataVersion(curVersion int32) int32 {
	if dataVersion != "" {
		dv, err := strconv.Atoi(dataVersion)
		checkError(errors.Wrap(err, "version invalid"))

		curVersion = int32(dv)
	}

	return curVersion
}

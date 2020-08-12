package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"text/tabwriter"
	"time"

	"github.com/beeeeeeenny/zkcmd/common/zookeeper"
	"github.com/pkg/errors"
	"github.com/samuel/go-zookeeper/zk"
	"github.com/spf13/cobra"
)

var (
	setCreate   bool
	force       bool
	dataVersion string
	outStat     bool

	ephemeral bool
	sequence  bool
)

func init() {
	rootCmd.AddCommand(nodeCmd)

	nodeCmd.AddCommand(nodeLsCmd)
	nodeLsCmd.Flags().BoolVarP(&outStat, "stat", "s", false, "node stat info")

	nodeCmd.AddCommand(nodeLsnCmd)

	nodeCmd.AddCommand(nodeGetCmd)
	nodeGetCmd.Flags().BoolVarP(&outStat, "stat", "s", false, "node stat info")

	nodeCmd.AddCommand(nodeDeleteCmd)
	nodeDeleteCmd.Flags().StringVarP(&dataVersion, "version", "v", "", "node data version")
	nodeDeleteCmd.Flags().BoolVarP(&force, "force", "f", false, "will force delete multi-level node, like: deleteall")

	nodeCmd.AddCommand(nodeSetCmd)
	nodeSetCmd.Flags().BoolVarP(&setCreate, "create", "c", false, "will create node if node does not exist, but does not directly create multi-level node")
	nodeSetCmd.Flags().BoolVarP(&force, "force", "f", false, "will force create multi-level node if node does not exist")
	nodeSetCmd.Flags().StringVarP(&dataVersion, "version", "v", "", "node data version")
	nodeSetCmd.Flags().BoolVarP(&outStat, "stat", "s", false, "node stat info")

	nodeCmd.AddCommand(nodeCreateCmd)
	nodeCreateCmd.Flags().BoolVarP(&force, "force", "f", false, "will force create multi-level node if node does not exist")
	nodeCreateCmd.Flags().BoolVarP(&ephemeral, "ephemeral", "e", false, "create ephemeral node")
	nodeCreateCmd.Flags().BoolVarP(&sequence, "sequence", "s", false, "create sequence node")
}

var nodeCmd = &cobra.Command{
	Use: "node",
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

var nodeLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List children, the path default: /",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		newZKClient()

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

		if outStat {
			outputStat(stat)
		}
	},
}

var nodeLsnCmd = &cobra.Command{
	Use:   "lsn",
	Short: "List all node, the path default: /",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		newZKClient()

		path := "/"
		if len(args) > 0 {
			path = args[0]
		}

		_, stat, err := zkcli.Children(path)
		checkError(err)

		w := tabwriter.NewWriter(os.Stdout, 6, 4, 3, ' ', 0)
		fmt.Fprintf(w, "ID\tPath\t\n")

		if stat.NumChildren != 0 {
			ns, err := zkcli.GetNodes(path)
			checkError(err)
			for i, n := range ns {
				fmt.Fprintf(w, "%v\t%v\t\n", i+1, n)
			}
			w.Flush()
		}
	},
}

var nodeGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get node value",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		newZKClient()

		d, stat, err := zkcli.Get(args[0])
		checkError(err)

		w := tabwriter.NewWriter(os.Stdout, 6, 4, 3, '\t', 0)
		fmt.Fprintf(w, "ChildrenNum:\t%v\t\n", stat.NumChildren)
		fmt.Fprintf(w, "Value:      \t%v\t\n", string(d))
		w.Flush()

		if outStat {
			outputStat(stat)
		}
	},
}

var nodeSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Update node value",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		newZKClient()

		path := args[0]
		data := args[1]

		exist, stat, err := zkcli.Exists(path)
		checkError(err)

		if exist {
			version := stat.Version
			if dataVersion != "" {
				dv, err := strconv.Atoi(dataVersion)
				checkError(errors.Wrap(err, "version invalid"))

				version = int32(dv)
			}

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

		if exist && outStat {
			outputStat(stat)
		}

		checkError(zk.ErrNoNode)
	},
}

var nodeCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create node",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		newZKClient()

		path := args[0]
		data := args[1]

		// check acl
		var acl string
		if len(args) > 2 {
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
	},
}

var nodeDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete node",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		newZKClient()

		exist, stat, err := zkcli.Exists(args[0])
		checkError(err)

		if exist {
			fmt.Println("#######: ", force)
			if force {
				err = zookeeper.ValidatePath(args[0], false)
				checkError(err)

				err = zkcli.ForceDelete(args[0])
				checkError(err)

				return
			}

			version := stat.Version
			if dataVersion != "" {
				dv, err := strconv.Atoi(dataVersion)
				checkError(errors.Wrap(err, "version invalid"))

				version = int32(dv)
			}

			err = zkcli.Delete(args[0], version)
			checkError(err)

			return
		}

		checkError(zk.ErrNoNode)
	},
}

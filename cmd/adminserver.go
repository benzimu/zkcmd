package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newCmdAdminServer() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "adminsrv",
		Short: `Zookeeper AdminServer, see: https://zookeeper.apache.org/doc/current/zookeeperAdmin.html#sc_adminserver`,
	}

	cmd.AddCommand(newCmdAdminServerList())
	cmd.AddCommand(newCmdAdminServerExec())

	return cmd
}

func newCmdAdminServerList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [flags]",
		Short: "list AdminServer all commands",
		Args:  cobra.ExactArgs(0),
		Run:   cmdRunAdminServerList,
	}

	cmdAdminServerCommonFlags(cmd)

	return cmd
}

func newCmdAdminServerExec() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "exec [flags] command",
		Short: `exec AdminServer command, command like: stats/stat, ruok, configuration/conf/config, is_read_only/isro`,
		Example: `  zkcmd adminsrv exec stat
	  zkcmd adminsrv exec conf

	  For more commands, see: https://zookeeper.apache.org/doc/current/zookeeperAdmin.html#sc_adminserver`,
		Args: cobra.ExactArgs(1),
		Run:  cmdRunAdminServerExec,
	}

	cmdAdminServerCommonFlags(cmd)

	return cmd
}

func cmdAdminServerCommonFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&zkcmdConf.AdminCommandURL, "adminCommandURL", "", "", fmt.Sprintf(`The AdminServer URL for listing and issuing commands relative to the root URL. (default "%s")`, defaultAdminCommandURL))
	cmd.Flags().StringSliceVarP(&zkcmdConf.AdminServer, "adminServer", "", nil, fmt.Sprintf("zookeeper AdminServer address, multiple addresses with a comma. (default [%s])", defaultAdminServer))
	_ = viper.BindPFlag("adminCommandURL", cmd.Flags().Lookup("adminCommandURL"))
	_ = viper.BindPFlag("adminServer", cmd.Flags().Lookup("adminServer"))
	viper.SetDefault("adminCommandURL", defaultAdminCommandURL)
	viper.SetDefault("adminServer", []string{defaultAdminServer})
}

func cmdRunAdminServerList(cmd *cobra.Command, args []string) {
	for _, srv := range zkcmdConf.AdminServer {
		fmt.Printf("############### AdminServer: %s ###############\n", srv)

		srvAddr := fmt.Sprintf("%s%s", srv, zkcmdConf.AdminCommandURL)
		res, err := doAdminServerReq(srvAddr)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println(res)
	}
}

func cmdRunAdminServerExec(cmd *cobra.Command, args []string) {
	command := args[0]

	for _, srv := range zkcmdConf.AdminServer {
		fmt.Printf("############### AdminServer: %s ###############\n", srv)

		srvAddr := fmt.Sprintf("%s%s/%s", srv, zkcmdConf.AdminCommandURL, command)
		res, err := doAdminServerReq(srvAddr)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println(res)
	}
}

func doAdminServerReq(addr string) (string, error) {
	if !strings.HasPrefix(addr, "http://") && !strings.HasPrefix(addr, "https://") {
		addr = "http://" + addr
	}

	resp, err := resty.New().
		SetTimeout(5 * time.Second).
		R().
		Get(addr)
	if err != nil {
		return "", err
	}

	return resp.String(), nil
}

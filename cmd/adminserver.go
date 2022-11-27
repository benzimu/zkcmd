package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultAdminCommandURL = "/commands"
	defaultAdminServer     = "127.0.0.1:8080"
)

var (
	adminServer     []string
	adminCommandURL string
)

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(adminServerCmd)

	adminServerCmd.AddCommand(adminServerExecCmd)
	adminServerCmd.AddCommand(adminServerListCmd)

	adminServerCmd.PersistentFlags().StringVarP(&adminCommandURL, "adminCmdURL", "", "", fmt.Sprintf(`The AdminServer URL for listing and issuing commands relative to the root URL. (default "%s")`, defaultAdminCommandURL))
	adminServerCmd.PersistentFlags().StringSliceVarP(&adminServer, "adminServer", "", nil, fmt.Sprintf("zookeeper AdminServer address, multiple addresses with a comma. (default [%s])", defaultAdminServer))
	viper.BindPFlag("adminCmdURL", adminServerCmd.PersistentFlags().Lookup("adminCommandURL"))
	viper.BindPFlag("adminServer", adminServerCmd.PersistentFlags().Lookup("adminServer"))

	fmt.Println(viper.GetStringSlice("server"))
}

var adminServerCmd = &cobra.Command{
	Use:   "adminsrv",
	Short: "Zookeeper AdminServer, see: https://zookeeper.apache.org/doc/current/zookeeperAdmin.html#sc_adminserver",
}

var adminServerListCmd = &cobra.Command{
	Use:   "list [flags]",
	Short: "list AdminServer all commands",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println(rootCmd.PersistentFlags().Lookup("server"))
		// fmt.Println(rootCmd.PersistentFlags().HasAvailableFlags())
		// fmt.Println(rootCmd.PersistentFlags().HasFlags())
		fmt.Println("QQQQQQQQ:", viper.GetStringSlice("server"))
		if len(adminServer) == 0 {
			adminServer = []string{defaultAdminServer}
		}

		if adminCommandURL == "" {
			adminCommandURL = defaultAdminCommandURL
		}

		for _, srv := range adminServer {
			fmt.Printf("############### AdminServer: %s ###############\n", srv)

			srvAddr := fmt.Sprintf("%s%s", srv, adminCommandURL)
			res, err := doAdminServerReq(srvAddr)
			if err != nil {
				fmt.Println(err)
				continue
			}

			fmt.Println(res)
		}
	},
}

var adminServerExecCmd = &cobra.Command{
	Use:   "exec [flags] command",
	Short: `exec AdminServer command, command like: stats/stat, ruok, configuration/conf/config, is_read_only/isro`,
	Example: `  zkcmd adminsrv exec stat
  zkcmd adminsrv exec conf

  For more commands, see: https://zookeeper.apache.org/doc/current/zookeeperAdmin.html#sc_adminserver`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		command := args[0]

		if len(adminServer) == 0 {
			adminServer = []string{defaultAdminServer}
		}

		if adminCommandURL == "" {
			adminCommandURL = defaultAdminCommandURL
		}

		for _, srv := range adminServer {
			fmt.Printf("############### AdminServer: %s ###############\n", srv)

			srvAddr := fmt.Sprintf("%s%s/%s", srv, adminCommandURL, command)
			res, err := doAdminServerReq(srvAddr)
			if err != nil {
				fmt.Println(err)
				continue
			}

			fmt.Println(res)
		}
	},
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

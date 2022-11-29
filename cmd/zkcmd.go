package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/benzimu/zkcmd/common/zookeeper"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	verbose bool

	zkcli *zookeeper.Client
)

func Execute() {
	err := newRootCommand().Execute()
	checkError(err)
}

func newRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "zkcmd",
		Short: "A brief description of your application",
		Long: `zkcmd is a command tool for zookeeper cluster management.
  This application can connect zookeeper server and list/create/delete/set znode.
  And use The Four Letter Words command, see: https://zookeeper.apache.org/doc/current/zookeeperAdmin.html#sc_4lw`,
	}

	cmd.PersistentFlags().StringVarP(&cfgFile, "config", "", "", `config file. (default "$HOME/.zkcmd.yaml")`)
	cmd.PersistentFlags().StringSliceVarP(&zkcmdConf.Server, "server", "", nil, fmt.Sprintf("zookeeper server address, multiple addresses with a comma. (default [%s])", defaultServer))
	cmd.PersistentFlags().StringSliceVarP(&zkcmdConf.ACL, "acl", "", nil, `zookeeper cluster ACL, multiple ACL with a comma. EX: "user:password"`)
	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "V", false, "whether to print verbose log")
	_ = viper.BindPFlag("server", cmd.PersistentFlags().Lookup("server"))
	_ = viper.BindPFlag("acl", cmd.PersistentFlags().Lookup("acl"))
	viper.SetDefault("server", []string{defaultServer})

	cobra.OnInitialize(initConfig)

	cmd.AddCommand(newCmd4lw())
	cmd.AddCommand(newCmdACL())
	cmd.AddCommand(newCmdAdminServer())
	cmd.AddCommand(newCmdConfig())
	cmd.AddCommand(newCmdVersion())
	cmd.AddCommand(newCmdZnode())

	return cmd
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		checkError(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".zkcmd")
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("zkcmd")

	if err := viper.ReadInConfig(); err == nil {
		if verbose {
			log.Println("Using config file:", viper.ConfigFileUsed())
		}

		err = viper.Unmarshal(zkcmdConf)
		checkError(err)
	}
}

// newZKClient new zookeeper client, if server non-empty
func newZKClient() *zookeeper.Client {
	zkcli, err := zookeeper.New(zkcmdConf.Server)
	checkError(errors.Wrap(err, "new zk client"))

	zkcli.EnableLogging(verbose)

	for _, a := range zkcmdConf.ACL {
		err = zkcli.AddAuth("digest", []byte(a))
		checkError(errors.Wrap(err, "add auth error"))
	}

	return zkcli
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

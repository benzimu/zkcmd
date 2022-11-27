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
	zkcli *zookeeper.Client
	acl   []string

	cfgFile string
	server  []string
	verbose bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "zkcmd",
	Short: "A brief description of your application",
	Long: `zkcmd is a command tool for zookeeper cluster management.
This application can connect zookeeper server and create/delete/set node.`,
}

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "", "", `config file. (default "$HOME/.zkcmd.yaml")`)
	rootCmd.PersistentFlags().StringSliceVarP(&server, "server", "", nil, "zookeeper server address, multiple addresses with a comma.")
	rootCmd.PersistentFlags().StringSliceVarP(&acl, "acl", "", nil, `zookeeper cluster ACL, multiple ACL with a comma. EX: "user:password"`)
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "V", false, "whether to print verbose log")
	// viper.BindPFlag("server", rootCmd.PersistentFlags().Lookup("server"))
	viper.BindPFlag("acl", rootCmd.PersistentFlags().Lookup("acl"))

	// fmt.Println(rootCmd.PersistentFlags().Lookup("server"))
	// fmt.Println(rootCmd.PersistentFlags().HasAvailableFlags())
	// fmt.Println(rootCmd.PersistentFlags().HasFlags())
	// fmt.Println(rootCmd.Flags().HasAvailableFlags())
	// fmt.Println(rootCmd.Name())

	// viper.SetDefault("server", []string{"127.0.0.1:2181"})
	fmt.Println(viper.GetStringSlice("server"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		// home, err := homedir.Dir()
		home, err := os.UserHomeDir()
		checkError(err)

		// Search config in home directory with name ".zkcmd" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".zkcmd")
	}

	viper.AutomaticEnv() // read in environment variables that match
	viper.SetEnvPrefix("zkcmd")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if verbose {
			log.Println("Using config file:", viper.ConfigFileUsed())
		}

		fmt.Println("EEEEEEE:", viper.GetStringSlice("server"))

		// if len(server) == 0 {
		// 	fmt.Println("SSSSSSS:", viper.GetStringSlice("server"))
		// 	server = viper.GetStringSlice("server")
		// 	fmt.Println("CCCCCCC:", viper.GetStringSlice("server"))
		// }
		// acl = viper.GetStringSlice("acl")
	}
}

// newZKClient new zookeeper client, if server non-empty
func newZKClient() {
	var err error
	zkcli, err = zookeeper.New(server)
	checkError(errors.Wrap(err, "new zk client"))

	zkcli.EnableLogging(verbose)

	for _, a := range acl {
		err = zkcli.AddAuth("digest", []byte(a))
		checkError(errors.Wrap(err, "add auth error"))
	}

}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/beeeeeeenny/zkcmd/common/zookeeper"

	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	zkcli *zookeeper.Client

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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.zkcmd.yaml)")
	rootCmd.PersistentFlags().StringSliceVar(&server, "server", nil,
		"zookeeper server address, multiple addresses with a comma (default is 127.0.0.1:2181)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "V", false, "whether to print verbose log")

	viper.SetDefault("server", []string{"127.0.0.1:2181"})

	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		checkError(err)

		// Search config in home directory with name ".zkcmd" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".zkcmd")
	}

	viper.AutomaticEnv() // read in environment variables that match
	viper.SetEnvPrefix("zkcmd")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if verbose {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}
	}
}

// newZKClient new zookeeper client, if server non-empty
func newZKClient() {
	if len(server) == 0 {
		server = viper.GetStringSlice("server")
	}

	var err error
	zkcli, err = zookeeper.New(server)
	checkError(errors.Wrap(err, "new zk client"))

	zkcli.EnableLogging(verbose)
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

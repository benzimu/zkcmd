package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const (
	defaultServer          = "127.0.0.1:2181"
	defaultAdminServer     = "127.0.0.1:8080"
	defaultAdminCommandURL = "/commands"
)

var zkcmdConf = &zkcmdConfig{}

type zkcmdConfig struct {
	Server          []string `yaml:"server"`
	ACL             []string `yaml:"acl"`
	AdminServer     []string `yaml:"adminServer"`
	AdminCommandURL string   `yaml:"adminCommandURL"`
}

func newCmdConfig() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "zkcmd config init and cat",
	}

	cmd.AddCommand(newCmdConfigInit())
	cmd.AddCommand(newCmdConfigCat())

	return cmd
}

func newCmdConfigInit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "init zkcmd config",
		Run:   cmdRunConfigInit,
	}

	return cmd
}

func newCmdConfigCat() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cat",
		Short: "cat zkcmd config",
		Run:   cmdRunConfigCat,
	}

	return cmd
}

func cmdRunConfigInit(cmd *cobra.Command, args []string) {
	reader := bufio.NewReader(os.Stdin)

	// input cluster address
	zkcmdConf.Server = inputClusterAddress(reader)

	// input cluster ACL
	zkcmdConf.ACL = inputClusterACL(reader)

	// input cluster AdminServer address
	zkcmdConf.AdminServer = inputAdminServerAddress(reader)

	// input cluster AdminServer command root URL
	zkcmdConf.AdminCommandURL = inputAdminServerCommandURL(reader)

	saveConfigFile(zkcmdConf)

	fmt.Println("########################################")
	fmt.Println("zkcmd config path:", getConfigFilePath())
}

func cmdRunConfigCat(cmd *cobra.Command, args []string) {
	cfgPath := getConfigFilePath()
	f, err := os.ReadFile(cfgPath)
	if os.IsNotExist(err) {
		fmt.Println(`The zkcmd configuration has not been initialized. Please init zkcmd config, use command:
	zkcmd config init`)
		os.Exit(1)
	}

	checkError(err)

	fmt.Println("########################################")
	fmt.Println("##### zkcmd config path:", getConfigFilePath())
	fmt.Println("##### zkcmd config data:")
	fmt.Println(string(f))
}

func saveConfigFile(c *zkcmdConfig) {
	cfgFilePath := getConfigFilePath()

	f, err := os.OpenFile(cfgFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	checkError(err)
	defer f.Close()

	enc := yaml.NewEncoder(f)

	err = enc.Encode(c)
	checkError(err)
}

func getConfigFilePath() string {
	home, err := os.UserHomeDir()
	checkError(errors.Wrap(err, "fail to get homedir"))

	return filepath.Join(home, ".zkcmd.yaml")
}

func inputClusterAddress(reader *bufio.Reader) []string {
	fmt.Printf(`Please input zookeeper cluster addresses, multiple addresses with a comma. Defaults to [%s] >
> `, defaultServer)
	server, err := reader.ReadString('\n')
	checkError(err)

	serverTrim := strings.TrimSpace(server)
	if serverTrim == "" {
		return []string{defaultServer}
	}

	return strings.Split(serverTrim, ",")
}

func inputClusterACL(reader *bufio.Reader) []string {
	fmt.Print(`Please input zookeeper cluster ACL, multiple ACL with a comma. EX: "user:password" >
> `)
	acl, err := reader.ReadString('\n')
	checkError(err)

	aclTrim := strings.TrimSpace(acl)
	if aclTrim == "" {
		return nil
	}

	as := strings.Split(aclTrim, ",")
	for _, s := range as {
		ss := strings.Split(s, ":")
		if len(ss) < 2 {
			fmt.Printf("Invalid ACL input: %s, EX: \"user:password\"\n", s)
			os.Exit(1)
		}
	}

	return as
}

func inputAdminServerAddress(reader *bufio.Reader) []string {
	fmt.Printf(`Please input zookeeper AdminServer addresses, multiple addresses with a comma. Defaults to [%s] >
> `, defaultAdminServer)
	server, err := reader.ReadString('\n')
	checkError(err)

	serverTrim := strings.TrimSpace(server)
	if serverTrim == "" {
		return []string{defaultAdminServer}
	}

	return strings.Split(serverTrim, ",")
}

func inputAdminServerCommandURL(reader *bufio.Reader) string {
	fmt.Printf(`Please input zookeeper AdminServer commandURL. Defaults to %s >
> `, defaultAdminCommandURL)
	commandURL, err := reader.ReadString('\n')
	checkError(err)

	commandURLTrim := strings.TrimSpace(commandURL)
	if commandURLTrim == "" {
		return defaultAdminCommandURL
	}

	return commandURLTrim
}

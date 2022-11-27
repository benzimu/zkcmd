package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(configCatCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "zkcmd config init and cat",
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Init zkcmd config",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		// input cluster address
		confServer := inputClusterAddress(reader)

		// input cluster ACL
		confACL := inputClusterACL(reader)

		// input cluster AdminServer address
		confAdminServer := inputAdminServerAddress(reader)

		// input cluster AdminServer command root URL
		confCommandURL := inputAdminServerCommandURL(reader)

		conf := confServer + confACL + confAdminServer + confCommandURL + "\n"
		saveConfigFile(conf)

		fmt.Println("########################################")
		fmt.Println("##### Config path:", getConfigFilePath())
		fmt.Println("##### Config data:")
		fmt.Println(conf)
	},
}

var configCatCmd = &cobra.Command{
	Use:   "cat",
	Short: "cat zkcmd config",
	Run: func(cmd *cobra.Command, args []string) {
		cfgPath := getConfigFilePath()
		f, err := os.ReadFile(cfgPath)
		checkError(err)

		fmt.Println("########################################")
		fmt.Println("##### Config path:", getConfigFilePath())
		fmt.Println("##### Config data:")
		fmt.Println(string(f))
	},
}

func saveConfigFile(s string) {
	cfgFilePath := getConfigFilePath()

	f, err := os.OpenFile(cfgFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	checkError(err)
	defer f.Close()

	_, err = f.WriteString(s)
	checkError(err)
}

func getConfigFilePath() string {
	home, err := homedir.Dir()
	checkError(errors.Wrap(err, "fail to get homedir"))

	return filepath.Join(home, ".zkcmd.yaml")
}

func inputClusterAddress(reader *bufio.Reader) string {
	fmt.Println("Please input zookeeper cluster addresses, multiple addresses with a comma. Defaults to 127.0.0.1:2181 > ")
	server, err := reader.ReadString('\n')
	checkError(err)

	serverTrim := strings.TrimSpace(server)

	confServer := `server:`
	if serverTrim == "" {
		confServer += `
  - 127.0.0.1:2181`
	} else {
		ss := strings.Split(serverTrim, ",")
		for _, s := range ss {
			confServer += fmt.Sprintf(`
  - %s`, s)
		}
	}

	return confServer
}

func inputClusterACL(reader *bufio.Reader) string {
	confACL := `
acl:`
	fmt.Println(`Please input zookeeper cluster ACL, multiple ACL with a comma. EX: "user:password" > `)
	acl, err := reader.ReadString('\n')
	checkError(err)

	aclTrim := strings.TrimSpace(acl)

	if aclTrim != "" {
		as := strings.Split(aclTrim, ",")
		for _, s := range as {
			ss := strings.Split(s, ":")
			if len(ss) < 2 {
				fmt.Printf("Invalid ACL input: %s, EX: \"user:password\"\n", s)
				os.Exit(1)
			}

			confACL += fmt.Sprintf(`
  - %s`, aclTrim)
		}
	}

	return confACL
}

func inputAdminServerAddress(reader *bufio.Reader) string {
	fmt.Println("Please input zookeeper AdminServer addresses, multiple addresses with a comma. Defaults to 127.0.0.1:8080 > ")
	server, err := reader.ReadString('\n')
	checkError(err)

	serverTrim := strings.TrimSpace(server)

	confAdminServer := `
adminServer:`
	if serverTrim == "" {
		confAdminServer += `
  - 127.0.0.1:8080`
	} else {
		ss := strings.Split(serverTrim, ",")
		for _, s := range ss {
			confAdminServer += fmt.Sprintf(`
  - %s`, s)
		}
	}

	return confAdminServer
}

func inputAdminServerCommandURL(reader *bufio.Reader) string {
	fmt.Println("Please input zookeeper AdminServer commandURL. Defaults to /commands > ")
	commandURL, err := reader.ReadString('\n')
	checkError(err)

	commandURLTrim := strings.TrimSpace(commandURL)

	confCommandURL := `
adminCommandURL: /commands`
	if commandURLTrim != "" {
		confCommandURL = fmt.Sprintf(`
adminCommandURL: %s`, commandURLTrim)
	}

	return confCommandURL
}

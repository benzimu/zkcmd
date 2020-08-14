package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Init zkcmd config",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		fmt.Println("Please input zookeeper cluster addresses, multiple addresses with a comma (default is 127.0.0.1:2181) > ")
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

		confACL := `
acl:`
		fmt.Println("Please input zookeeper cluster ACL, EX: \"digest user:password\" > ")
	aclScan:
		acl, err := reader.ReadString('\n')
		checkError(err)
		aclTrim := strings.TrimSpace(acl)

		if aclTrim != "" {
			as := strings.Split(aclTrim, " ")
			if len(as) < 2 {
				fmt.Println("Invalid ACL input, EX: \"digest user:password\"")
				os.Exit(1)
			}

			ss := strings.Split(as[1], ":")
			if len(ss) < 2 {
				fmt.Println("Invalid ACL input, EX: \"digest user:password\"")
				os.Exit(1)
			}

			confACL += fmt.Sprintf(`
  - %s`, aclTrim)
			goto aclScan
		}

		conf := confServer + confACL + "\n"
		saveConfigFile(conf)
	},
}

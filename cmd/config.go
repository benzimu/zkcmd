package cmd

import (
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
		fmt.Println("Please input zookeeper cluster addresses, multiple addresses with a comma (default is 127.0.0.1:2181) > ")
		var server string
		fmt.Scanln(&server)

		confServer := `
server:
`

		if server == "" {
			confServer += `
  - 127.0.0.1:2181
`
		} else {
			ss := strings.Split(server, ",")
			for _, s := range ss {
				confServer += fmt.Sprintf(`
  - %s
`, s)
			}
		}
		saveConfigFile(confServer)

		acls := make([]string, 0)
		fmt.Println("Please input zookeeper cluster ACL, EX: \"digest root:root\" > ")
	aclScan:
		var acl string
		fmt.Scanln(&acl)
		if acl != "" {
			as := strings.Split(acl, " ")
			if len(as) < 2 {
				fmt.Println("Invalid ACL")
				os.Exit(1)
			}

			ss := strings.Split(as[1], ":")
			if len(ss) < 2 {
				fmt.Println("Invalid ACL")
				os.Exit(1)
			}
			acls = append(acls, acl)
			goto aclScan
		}
	},
}

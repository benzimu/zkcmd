package cmd

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/spf13/cobra"
)

func newCmd4lw() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "4lw [flags] 4lwcmd",
		Short: `Zookeeper the four letter word commands, 4lwcmd like: stat, ruok, conf, isro`,
		Example: `  zkcmd 4lw stat
	  zkcmd 4lw conf

	  For more the four letter word commands, see:
		https://zookeeper.apache.org/doc/current/zookeeperAdmin.html#sc_4lw`,
		ValidArgs: []string{"conf", "cons", "crst", "dump", "envi", "ruok",
			"srst", "srvr", "stat", "wchs", "wchc", "dirs", "wchp", "mntr", "isro",
			"hash", "gtmk", "stmk", "icfg", "lsnp", "lead", "orst", "obsr", "sysp"},
		Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		Run:  cmdRun4lw,
	}

	return cmd
}

func cmdRun4lw(cmd *cobra.Command, args []string) {
	fourlwcmd := args[0]

	for _, srvAddr := range zkcmdConf.Server {
		fmt.Printf("############### Server: %s ###############\n", srvAddr)

		conn, err := net.DialTimeout("tcp", srvAddr, 3*time.Second)
		if err != nil {
			fmt.Println(err)
			continue
		}

		err = conn.SetDeadline(time.Now().Add(8 * time.Second))
		if err != nil {
			fmt.Println(err)
			continue
		}

		_, err = conn.Write([]byte(fourlwcmd))
		if err != nil {
			fmt.Println(err)
			continue
		}

		resData, err := io.ReadAll(conn)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println(string(resData))
	}
}

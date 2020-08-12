package zookeeper

import "fmt"

type logger struct {
	enable bool
}

func (l logger) Printf(format string, a ...interface{}) {
	if l.enable {
		fmt.Printf(format, a...)
		fmt.Printf("\n")
	}
}

package zookeeper

import "fmt"

type logger struct {
	enable bool
}

func (l logger) Printf(format string, a ...interface{}) {
	if l.enable {
		fmt.Printf("go-zookeeper: "+format+"\n", a...)
	}
}

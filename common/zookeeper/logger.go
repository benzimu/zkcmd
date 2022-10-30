package zookeeper

import "log"

type logger struct {
	enable bool
}

func (l logger) Printf(format string, a ...interface{}) {
	if l.enable {
		log.Printf("go-zookeeper: "+format+"\n", a...)
	}
}

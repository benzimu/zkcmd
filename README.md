# ZKCmd

ZKCmd: ZooKeeper command line tool, written in Go. This application can connect zookeeper server and list/create/delete/set znode, and use The Four Letter Words command, see: [4lw](https://zookeeper.apache.org/doc/current/zookeeperAdmin.html#sc_4lw)

## Installation

- Install from source:

```bash
$> go install github.com/benzimu/zkcmd@latest
```

- Install from Makefile

```bash
$> git clone https://github.com/benzimu/zkcmd.git
$> cd zkcmd
$> make
```

## Usage

```bash
$> zkcmd -h
zkcmd is a command tool for zookeeper cluster management.
  This application can connect zookeeper server and list/create/delete/set znode.
  And use The Four Letter Words command, see: https://zookeeper.apache.org/doc/current/zookeeperAdmin.html#sc_4lw

Usage:
  zkcmd [command]

Available Commands:
  4lw         Zookeeper the four letter word commands, 4lwcmd like: stat, ruok, conf, isro
  acl         Znode ACL command
  adminsrv    Zookeeper AdminServer, see: https://zookeeper.apache.org/doc/current/zookeeperAdmin.html#sc_adminserver
  completion  Generate the autocompletion script for the specified shell
  config      zkcmd config init and cat
  help        Help about any command
  version     Print version information of zkcmd and quit
  znode       Znode command

Flags:
      --acl strings      zookeeper cluster ACL, multiple ACL with a comma. EX: "user:password"
      --config string    config file. (default "$HOME/.zkcmd.yaml")
  -h, --help             help for zkcmd
      --server strings   zookeeper server address, multiple addresses with a comma. (default [127.0.0.1:2181])
  -V, --verbose          whether to print verbose log

Use "zkcmd [command] --help" for more information about a command.
```

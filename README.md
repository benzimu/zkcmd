# ZKCmd

ZKCmd: ZooKeeper command line tool, written in Go.

## Installation

- Install from source:

```bash
$> go install github.com/benzimu/zkcmd@latest
```

## Usage

```bash
$> zkcmd -h
zkcmd is a command tool for zookeeper cluster management.
This application can connect zookeeper server and create/delete/set node.

Usage:
  zkcmd [command]

Available Commands:
  acl         Znode ACL command
  completion  Generate the autocompletion script for the specified shell
  config      Zkcmd config
  help        Help about any command
  version     Print version information of zkcmd and quit
  znode       Znode command

Flags:
      --config string    config file (default is $HOME/.zkcmd.yaml)
  -h, --help             help for zkcmd
      --server strings   zookeeper server address, multiple addresses with a comma (default is 127.0.0.1:2181)
  -V, --verbose          whether to print verbose log

Use "zkcmd [command] --help" for more information about a command.
```

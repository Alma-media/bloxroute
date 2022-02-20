# bloxroute
Test project
## Localstack
```bash
[host@user cmd]$ docker-compose up
```
## Server
### Options
```
Usage:
  bloxroute daemon [flags]

Aliases:
  daemon, d, server, start

Flags:
  -h, --help            help for daemon
  -o, --output string   server output
```
### Example
```bash
[host@user cmd]$ go run main.go start -o /home/alma/bloxroute.txt
```
## Client
```
Usage:
  bloxroute client [command]

Aliases:
  client, cli

Available Commands:
  add         add new item
  all         get all items
  del         remove item
  get         get all items

Flags:
  -h, --help   help for client
```
### Examples
- add item:
```bash
[host@user cmd]$ go run main.go cli add -k foo -v bar
```
- get item:
```bash
[host@user cmd]$ go run main.go cli get -k foo
```
- get all items:
```bash
[host@user cmd]$ go run main.go cli all
```
- remove item:
```bash
[host@user cmd]$ go run main.go cli del -k foo
```
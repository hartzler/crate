package main

import (
	"github.com/codegangsta/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Version = "0.1"
	app.Name = "2ndroute"
	app.Usage = "create virtual wiring for cloud components"
	app.Commands = []cli.Command{
		SetupCommand,
		CreateCommand,
	}
	app.Run(os.Args)
}

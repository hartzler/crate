package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/opencontainers/runc/libcontainer/configs"
)

var netnsCommand = cli.Command{
	Name:        "netns",
	Usage:       "list the path of the container's net namespace",
	Description: "args: <id>",
	Flags: []cli.Flag{
		cli.StringFlag{Name: "id", Usage: "ID for the container"},
	},
	Action: func(context *cli.Context) {
		args := context.Args()
		if len(args) != 1 {
			fatal(fmt.Errorf("expected 1 arguments <id>: %d", len(args)))
		}
		id := args[0]

		container, err := getContainer(context, id)
		if err != nil {
			fatal(err)
		}
		status, err := container.Status()
		if err != nil {
			fatal(err)
		}
		fmt.Println(status.NamespacePaths[configs.NEWNET])
	},
}

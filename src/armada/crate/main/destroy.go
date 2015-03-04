package main

import (
	"github.com/codegangsta/cli"
)

var destroyCommand = cli.Command{
	Name:  "destroy",
	Usage: "destroy the container",
	Flags: []cli.Flag{
		cli.StringFlag{Name: "id", Usage: "ID of the container"},
	},
	Action: func(context *cli.Context) {
		container, err := getContainer(context)
		if err != nil {
			fatal(err)
		}
		if err = container.Destroy(); err != nil {
			fatal(err)
		}
	},
}

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
		if err := fromContext(context).Destroy(context.String("id")); err != nil {
			fatal(err)
		}
	},
}

package main

import (
	"github.com/codegangsta/cli"
)

var pauseCommand = cli.Command{
	Name:  "pause",
	Usage: "pause the container's processes",
	Flags: []cli.Flag{
		cli.StringFlag{Name: "id", Usage: "ID for the container"},
	},
	Action: func(context *cli.Context) {
		container, err := getContainer(context)
		if err != nil {
			fatal(err)
		}
		if err = container.Pause(); err != nil {
			fatal(err)
		}
	},
}

var unpauseCommand = cli.Command{
	Name:  "resume",
	Usage: "resume the container's processes",
	Flags: []cli.Flag{
		cli.StringFlag{Name: "id", Value: "nsinit", Usage: "specify the ID for a container"},
	},
	Action: func(context *cli.Context) {
		container, err := getContainer(context)
		if err != nil {
			fatal(err)
		}
		if err = container.Resume(); err != nil {
			fatal(err)
		}
	},
}

package main

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

var statusCommand = cli.Command{
	Name:        "status",
	Usage:       "show the status of a container",
	Description: "args: <id>",
	Flags: []cli.Flag{
		cli.StringFlag{Name: "id", Usage: "ID of the container"},
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

		if err = json.NewEncoder(log.StandardLogger().Out).Encode(&status); err != nil {
			fatal(err)
		}
	},
}

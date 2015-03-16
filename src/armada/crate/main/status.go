package main

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

var statusCommand = cli.Command{
	Name:  "status",
	Usage: "show the status of a container",
	Flags: []cli.Flag{
		cli.StringFlag{Name: "id", Usage: "ID of the container"},
	},
	Action: func(context *cli.Context) {
		container, err := getContainer(context)
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

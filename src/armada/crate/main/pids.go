package main

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

var pidsCommand = cli.Command{
	Name:  "pids",
	Usage: "list the pids of a container",
	Flags: []cli.Flag{
		cli.StringFlag{Name: "id", Usage: "ID for the container"},
	},
	Action: func(context *cli.Context) {
		container, err := getContainer(context)
		if err != nil {
			fatal(err)
		}
		pids, err := container.Processes()
		if err != nil {
			fatal(err)
		}

		if err = json.NewEncoder(log.StandardLogger().Out).Encode(&pids); err != nil {
			fatal(err)
		}
	},
}

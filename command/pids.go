package command

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"os"
)

var pidsCommand = cli.Command{
	Name:        "pids",
	Usage:       "list the pids of a container",
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
		pids, err := container.Processes()
		if err != nil {
			fatal(err)
		}

		if err = json.NewEncoder(os.Stdout).Encode(&pids); err != nil {
			fatal(err)
		}
	},
}

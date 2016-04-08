package command

import (
	"fmt"
	"github.com/codegangsta/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:        "destroy",
		Usage:       "destroys a container",
		Description: "args: <id>",
		Action: func(context *cli.Context) {
			args := context.Args()
			if len(args) != 1 {
				fatal(fmt.Errorf("expected 1 arguments <id>: %d", len(args)))
			}
			id := args[0]

			if err := fromContext(context).Destroy(id); err != nil {
				fatal(err)
			}
		},
	})
}

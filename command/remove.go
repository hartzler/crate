package command

import (
	"fmt"
	"github.com/codegangsta/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:        "remove",
		Usage:       "uninstalls a crate",
		Description: "args: <id>",
		Action: func(context *cli.Context) {
			args := context.Args()
			if len(args) != 1 {
				fatal(fmt.Errorf("expected 1 arguments <id>: %d", len(args)))
			}

			if err := FromContext(context).Remove(args[0]); err != nil {
				fatal(err)
			}
		},
	})
}

package command

import (
	"fmt"
	"github.com/codegangsta/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:        "pause",
		Usage:       "pause the container's processes",
		Description: "args: <id>",
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
			if err = container.Pause(); err != nil {
				fatal(err)
			}
		},
	})

	Commands = append(Commands, cli.Command{
		Name:        "resume",
		Usage:       "resume the container's processes",
		Description: "args: <id>",
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
			if err = container.Resume(); err != nil {
				fatal(err)
			}
		},
	})
}

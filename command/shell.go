package command

import (
	"fmt"
	"github.com/codegangsta/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:        "shell",
		Usage:       "run a shell inside a container",
		Description: "args: <id>",
		Action:      runShell,
		Flags: []cli.Flag{
			cli.StringFlag{Name: "user,u", Value: "root", Usage: "set the user, uid, and/or gid for the process"},
			cli.StringFlag{Name: "cwd", Value: "", Usage: "set the current working dir"},
			cli.StringSliceFlag{Name: "env", Value: &cliEnv, Usage: "set environment variables for the process"},
		},
	})
}

func runShell(context *cli.Context) {
	args := context.Args()
	if len(args) < 1 {
		fatal(fmt.Errorf("expected 1 arguments <id>"))
	}
	id := args[0]

	err := fromContext(context).Shell(id)
	if err != nil {
		fatal(err)
	}
}

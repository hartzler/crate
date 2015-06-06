package main

import (
	"armada/crate"
	"fmt"
	"github.com/codegangsta/cli"
)

var cliEnv = cli.StringSlice(crate.StandardEnvironment)

var runCommand = cli.Command{
	Name:        "run",
	Usage:       "start a process inside a container",
	Description: "args: <id> <path> [arguments...]",
	Action:      runAction,
	Flags: []cli.Flag{
		cli.StringFlag{Name: "user,u", Value: "root", Usage: "set the user, uid, and/or gid for the process"},
		cli.StringFlag{Name: "cwd", Value: "", Usage: "set the current working dir"},
		cli.StringSliceFlag{Name: "env", Value: &cliEnv, Usage: "set environment variables for the process"},
	},
}

func runAction(context *cli.Context) {
	args := context.Args()
	if len(args) < 2 {
		fatal(fmt.Errorf("expected 2+ arguments <id> <path> [arguments...]"))
	}
	id := args[0]

	err := fromContext(context).Run(id, crate.RunArgs{
		Args: args[1:],
		Env:  context.StringSlice("env"),
		User: context.String("user"),
		Cwd:  context.String("cwd"),
	})
	if err != nil {
		fatal(err)
	}
}

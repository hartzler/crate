package main

import (
	"armada/crate"
	"github.com/codegangsta/cli"
)

var cliEnv = cli.StringSlice(crate.StandardEnvironment)

var runCommand = cli.Command{
	Name:   "run",
	Usage:  "start a process inside a container",
	Action: runAction,
	Flags: []cli.Flag{
		cli.StringFlag{Name: "id", Usage: "specify the ID for a container"},
		cli.BoolFlag{Name: "tty,t", Usage: "allocate a TTY to the container"},
		cli.StringFlag{Name: "user,u", Value: "root", Usage: "set the user, uid, and/or gid for the process"},
		cli.StringFlag{Name: "cwd", Value: "", Usage: "set the current working dir"},
		cli.StringSliceFlag{Name: "env", Value: &cliEnv, Usage: "set environment variables for the process"},
	},
}

func runAction(context *cli.Context) {
	err := fromContext(context).Run(context.String("id"), crate.RunArgs{
		Args: context.Args(),
		Env:  context.StringSlice("env"),
		User: context.String("user"),
		Cwd:  context.String("cwd"),
	})
	if err != nil {
		fatal(err)
	}
}

package main

import (
	"github.com/codegangsta/cli"
	"github.com/docker/libcontainer"
	"os"
)

var cliEnv = cli.StringSlice(standardEnvironment)

var runCommand = cli.Command{
	Name:   "run",
	Usage:  "start a process inside a container",
	Action: runAction,
	Flags: append([]cli.Flag{
		cli.BoolFlag{Name: "tty,t", Usage: "allocate a TTY to the container"},
		cli.StringFlag{Name: "user,u", Value: "root", Usage: "set the user, uid, and/or gid for the process"},
		cli.StringFlag{Name: "cwd", Value: "", Usage: "set the current working dir"},
		cli.StringSliceFlag{Name: "env", Value: &cliEnv, Usage: "set environment variables for the process"},
	}, createFlags...),
}

func runAction(context *cli.Context) {
	container, err := getContainer(context)
	if err != nil {
		fatal(err)
	}

	process := &libcontainer.Process{
		Args:   context.Args(),
		Env:    context.StringSlice("env"),
		User:   context.String("user"),
		Cwd:    context.String("cwd"),
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	err = container.Start(process)
	if err != nil {
		fatal(err)
	}
}

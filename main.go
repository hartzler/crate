package main

import (
	"github.com/armada-io/crate/command"
	"github.com/codegangsta/cli"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Version = "0.1"
	app.Name = "crate"
	app.Usage = "simple read only chroot based package manager"
	app.Author = "Matt Hartzler"
	app.Email = "matt@armada.io"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "root", Value: "/var/lib/crate", Usage: "root directory for crate state"},
		cli.StringFlag{Name: "log-file", Value: "", Usage: "set the log file to output logs to"},
		cli.BoolFlag{Name: "debug", Usage: "enable debug output in the logs"},
	}
	app.Commands = command.Commands
	app.Before = func(context *cli.Context) error {
		if context.GlobalBool("debug") {
			// TODO
		}
		if path := context.GlobalString("log-file"); path != "" {
			f, err := os.Create(path)
			if err != nil {
				return err
			}
			log.SetOutput(f)
		}
		// make sure we have our root setup
		return command.FromContext(context).SetupRoot()
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

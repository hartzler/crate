package main

import (
	"armada/crate"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/libcontainer"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Version = "0.1"
	app.Name = "crate"
	app.Usage = "manage containers and connections"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "root", Value: "/var/lib/crate/containers", Usage: "root directory for containers"},
		cli.StringFlag{Name: "log-file", Value: "crate.log", Usage: "set the log file to output logs to"},
		cli.BoolFlag{Name: "debug", Usage: "enable debug output in the logs"},
	}
	app.Commands = []cli.Command{
		SetupCommand,
		CreateCommand,
		runCommand,
		destroyCommand,
		pauseCommand,
		unpauseCommand,
		pidsCommand,
		InitCommand,
	}
	app.Before = func(context *cli.Context) error {
		if context.GlobalBool("debug") {
			log.SetLevel(log.DebugLevel)
		}
		if path := context.GlobalString("log-file"); path != "" {
			f, err := os.Create(path)
			if err != nil {
				return err
			}
			log.SetOutput(f)
		}
		// make sure we have our root setup
		return fromContext(context).SetupRoot()
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func fromContext(context *cli.Context) *crate.Crate {
	return crate.New(context.GlobalString("root"))
}

func getContainer(context *cli.Context) (*crate.Container, error) {
	return fromContext(context).Load(context.String("id"))
}

func fatal(err error) {
	if lerr, ok := err.(libcontainer.Error); ok {
		lerr.Detail(os.Stderr)
		os.Exit(1)
	}
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func fatalf(t string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, t, v...)
	os.Exit(1)
}

package main

import (
	"armada/crate"
	"armada/crate/pid1"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/opencontainers/runc/libcontainer"
	"os"
	"path/filepath"
	"runtime"
)

func init() {
	if len(os.Args) > 1 && os.Args[1] == "init" {
		runtime.GOMAXPROCS(1)
		runtime.LockOSThread()
		factory, _ := libcontainer.New("")
		if err := factory.StartInitialization(); err != nil {
			fatal(err)
		}
		panic("--this line should have never been executed, congratulations--")
	}
}

func main() {
	if filepath.Base(os.Args[0]) == crate.CRATE_INIT {
		startPid1()
	} else {
		runCrate()
	}
}

// the container PID 1
func startPid1() {
	pid1.Start()
}

func runCrate() {
	app := cli.NewApp()
	app.Version = "0.1"
	app.Name = "crate"
	app.Usage = "manage containers and connections"
	app.Author = "Matt Hartzler"
	app.Email = "matt@armada.io"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "root", Value: "/var/lib/crate", Usage: "root directory for crate state"},
		cli.StringFlag{Name: "log-file", Value: "", Usage: "set the log file to output logs to"},
		cli.BoolFlag{Name: "debug", Usage: "enable debug output in the logs"},
	}
	app.Commands = []cli.Command{
		newCommand,
		runCommand,
		destroyCommand,
		pauseCommand,
		unpauseCommand,
		pidsCommand,
		statusCommand,
		shellCommand,
		netnsCommand,
		// agentCommand,
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

func getContainer(context *cli.Context, id string) (*crate.Container, error) {
	return fromContext(context).Load(id)
}

func fatal(err error) {
	if lerr, ok := err.(libcontainer.Error); ok {
		lerr.Detail(os.Stderr)
		os.Exit(1)
	}
	fmt.Fprintln(os.Stderr, "ERROR:", err)
	os.Exit(1)
}

func fatalf(t string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, t, v...)
	os.Exit(1)
}

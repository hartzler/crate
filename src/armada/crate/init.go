package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/docker/libcontainer"
	"os"
	"runtime"
)

var InitCommand = cli.Command{
	Name:  "init",
	Usage: "(internal) runs the init process inside the namespace",
	Action: func(c *cli.Context) {
		runtime.GOMAXPROCS(1)
		runtime.LockOSThread()
		factory, err := libcontainer.New("")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if err := factory.StartInitialization(3); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		panic("This line should never been executed")
	},
}

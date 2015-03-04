package main

import (
	"fmt"
	"github.com/codegangsta/cli"
)

var CreateCommand = cli.Command{
	Name:  "create",
	Usage: "creates a new container",
	Flags: createFlags,
	Action: func(context *cli.Context) {
		//comp := c.Args()[0]
		//rootfs := "/home/vagrant/busybox"

		id := context.String("id")
		config := getTemplate(id)
		modify(config, context)

		fmt.Println(config)

		crate := fromContext(context)
		if _, err := crate.Create(id, config); err != nil {
			fatal(err)
		}
	},
}

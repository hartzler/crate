package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/docker/libcontainer/configs"
)

var CreateCommand = cli.Command{
	Name:  "create",
	Usage: "creates a new container",
	Flags: createFlags,
	Action: func(c *cli.Context) {

		config := getTemplate()
		modify(config, c)

		//comp := c.Args()[0]
		//rootfs := "/home/vagrant/busybox"

		if err := create(config); err != nil {
			fmt.Println(err)
		}
	},
}

func create(config *configs.Config) error {
  factory, err := loadFactory(context)
  if err != nil {
    fatal(err)
  }

  container, err := factory.Load(context.String("id"))
  if err != nil {
    created = true
    if container, err = factory.Create(context.String("id"), config); err != nil {
      fatal(err)
    }
  }
}

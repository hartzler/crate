package main

import (
	"fmt"
	"github.com/codegangsta/cli"
)

var SetupCommand = cli.Command{
	Name:  "setup",
	Usage: "2ndroute bridge [bridge-name]",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "bridge, b",
			Value: "armada0",
			Usage: "the bridge name to use",
		},
		cli.StringFlag{
			Name:  "bip",
			Value: "10.4.0.1/8",
			Usage: "the cidr to use for the bridge",
		},
	},
	Action: func(c *cli.Context) {
		name := c.String("bridge")
		bridgeip := c.String("bip")
		if err := setup(name, bridgeip); err != nil {
			fmt.Println(err)
		}
	},
}

func setup(name, bridgeip string) error {
	fmt.Println("# setup bridge")
	cmd, err := ipLink("add name", name, "type bridge")
	if err != nil {
		return err
	}
	fmt.Println(cmd)

	cmd, err = ip("addr add", bridgeip, "dev", name)
	if err != nil {
		return err
	}
	fmt.Println(cmd)

	cmd, err = ipLink("set", name, "up")
	if err != nil {
		return err
	}
	fmt.Println(cmd)

	return nil
}

package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	//"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Version = "0.1"
	app.Name = "2ndroute"
	app.Usage = "create virtual wiring for cloud components"
	app.Commands = []cli.Command{
		{
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
		},
		{
			Name:  "create",
			Usage: "2ndroute create <component-name> <component-ip>",
			Action: func(c *cli.Context) {
				if len(c.Args()) != 2 {
					fmt.Println("invalid args...")
					return
				}

				comp := c.Args()[0]
				compip := c.Args()[1]

				if err := create(comp, compip); err != nil {
					fmt.Println(err)
				}
			},
		},
	}
	app.Run(os.Args)
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

func create(comp, ipaddr string) error {
	fmt.Println("# create netns")
	cmd, err := ip("netns", "add", comp)
	if err != nil {
		return err
	}
	fmt.Println(cmd)

	fmt.Println("# create veth pair and attach one to bridge and move other into netns")
	name := fmt.Sprintf("veth0.%s", comp)
	peer := fmt.Sprintf("veth1.%s", comp)
	cmd, err = ipLink("add name", name, "type veth peer name", peer)
	if err != nil {
		return err
	}
	fmt.Println(cmd)

	// attach to bridge
	cmd, err = ipLink("set", name, "master", "armada0")
	if err != nil {
		return err
	}
	fmt.Println(cmd)

	// move to ns
	cmd, err = ipLink("set", peer, "netns", comp)
	if err != nil {
		return err
	}
	fmt.Println(cmd)

	fmt.Println("# setup container interfaces")
	cmd, err = ip("netns exec", comp, "ip link set", peer, "name eth0")
	if err != nil {
		return err
	}
	fmt.Println(cmd)

	// bring up lo
	cmd, err = ip("netns exec", comp, "ip lo up")
	if err != nil {
		return err
	}
	fmt.Println(cmd)

	// assign eth0
	cmd, err = ip("netns exec", comp, "ip addr add", ipaddr, "dev eth0")
	if err != nil {
		return err
	}
	fmt.Println(cmd)

	// up eth0
	cmd, err = ip("netns exec", comp, "ip eth0 up")
	if err != nil {
		return err
	}
	fmt.Println(cmd)

	// setup default gw
	cmd, err = ip("netns exec", comp, "ip route add default via 10.4.0.1")
	if err != nil {
		return err
	}
	fmt.Println(cmd)

	return nil
}

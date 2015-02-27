package main

import (
	"fmt"
	"github.com/codegangsta/cli"
)

var CreateCommand = cli.Command{
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

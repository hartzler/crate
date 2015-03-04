package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/docker/libcontainer/netlink"
	"net"
	"os"
)

var SetupCommand = cli.Command{
	Name:  "setup",
	Usage: "create the network bridge [bridge-name]",
	Flags: []cli.Flag{
		cli.StringFlag{Name: "bip", Value: "10.4.0.255/16", Usage: "ID for the container"},
		cli.StringFlag{Name: "bridge", Value: "armada0", Usage: "name for the armada bridge"},
	},
	Action: func(c *cli.Context) {
		name := c.String("bridge")
		bridgeip := c.String("bip")
		if err := setup(name, bridgeip); err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	},
}

func setup(bridgename, bridgeip string) error {
	if err := netlink.CreateBridge(bridgename, true); err != nil {
		return err
	}

	iface, err := net.InterfaceByName(bridgename)
	if err != nil {
		return err
	}

	ip, cidr, err := net.ParseCIDR(bridgeip)
	if err != nil {
		return err
	}

	if err = netlink.NetworkLinkAddIp(iface, ip, cidr); err != nil {
		return err
	}

	if err = netlink.NetworkLinkUp(iface); err != nil {
		return err
	}

	return nil
}

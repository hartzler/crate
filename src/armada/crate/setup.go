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
	Usage: "2ndroute bridge [bridge-name]",
	Flags: createFlags,
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

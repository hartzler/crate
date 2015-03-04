package main

import (
	"bytes"
	"fmt"
	"github.com/docker/libcontainer/netlink"
	"github.com/docker/libcontainer/utils"
	"net"
	"os/exec"
	"strings"
	"text/template"
)

type Network struct {
	Bridge            string
	Address           string
	NamespaceFd       uintptr
	TxQueueLen        int
	Mtu               int
	HostInterfaceName string
	TempVethPeerName  string
}

func (self *Network) netlinkCreate() (err error) {
	if self.HostInterfaceName == "" {
		tmpName, err := self.generateTempPeerName()
		if err != nil {
			return err
		}
		self.HostInterfaceName = tmpName
	}
	tmpName, err := self.generateTempPeerName()
	if err != nil {
		return err
	}
	self.TempVethPeerName = tmpName
	defer func() {
		if err != nil {
			netlink.NetworkLinkDel(self.HostInterfaceName)
			netlink.NetworkLinkDel(self.TempVethPeerName)
		}
	}()

	fmt.Println("MY NETWORK: settings: ", self)

	if self.Bridge == "" {
		return fmt.Errorf("bridge is not specified")
	}
	fmt.Println("MY NETWORK: start")
	bridge, err := net.InterfaceByName(self.Bridge)
	if err != nil {
		return err
	}
	fmt.Println("MY NETWORK: got bridge")
	if err := netlink.NetworkCreateVethPair(self.HostInterfaceName, self.TempVethPeerName, self.TxQueueLen); err != nil {
		return err
	}
	fmt.Println("MY NETWORK: created veth pair")
	host, err := net.InterfaceByName(self.HostInterfaceName)
	if err != nil {
		return err
	}
	fmt.Println("MY NETWORK: got host veth")
	if err := netlink.AddToBridge(host, bridge); err != nil {
		return err
	}
	fmt.Println("MY NETWORK: added host veth to bridge")
	if err := netlink.NetworkSetMTU(host, self.Mtu); err != nil {
		return err
	}
	fmt.Println("MY NETWORK: set MTU")
	if err := netlink.NetworkLinkUp(host); err != nil {
		return err
	}
	fmt.Println("MY NETWORK: up'd host veth")
	child, err := net.InterfaceByName(self.TempVethPeerName)
	if err != nil {
		return err
	}
	fmt.Println("MY NETWORK: got temp veth")
	err = netlink.NetworkSetNsFd(child, int(self.NamespaceFd))
	if err != nil {
		return err
	}
	fmt.Println("MY NETWORK: moved temp veth to namespace fd: ", self.NamespaceFd)
	return nil
}

func (self *Network) generateTempPeerName() (string, error) {
	return utils.GenerateRandomName("armada", 7)
}

type ipContext struct {
	Name        string
	AddressCidr string
	Bridge      string
	BridgeIp    string
}

var ipCommands = []string{
	// create netns
	"ip netns add {{.Name}}",
	// create veth pair and attach one to bridge and move other into netns
	"ip link add name armada.{{.Name}} type veth peer name veth.{{.Name}}",
	"ip link set armada.{{.Name}} master {{.Bridge}}",
	"ip link set armada.{{.Name}} up",
	"ip link set veth.{{.Name}} netns {{.Name}}",
	// setup container interfaces
	"ip netns exec {{.Name}} ip link set lo up",
	"ip netns exec {{.Name}} ip link set veth.{{.Name}} name eth0",
	"ip netns exec {{.Name}} ip link set eth0 up",
	"ip netns exec {{.Name}} ip addr add {{.AddressCidr}} dev eth0",
	"ip netns exec {{.Name}} ip route add default via {{.BridgeIp}} dev eth0",
}

func (self *Network) ipCreate(context ipContext) error {
	fmt.Println("IPCREATE: START")
	for _, cmd := range ipCommands {
		// template it
		var buf bytes.Buffer
		tmpl, err := template.New("ip").Parse(cmd)
		if err != nil {
			return err
		}
		err = tmpl.Execute(&buf, context)
		if err != nil {
			return err
		}

		// run it
		fmt.Println("IPCREATE: TEMPLATE RESULT: ", buf.String())
		slice := strings.Split(buf.String(), " ")
		cmd := exec.Command(slice[0], slice[1:]...)
		err = cmd.Run()
		if err != nil {
			fmt.Println("IPCREATE: command error: ", err)
		}
	}
	return nil
}

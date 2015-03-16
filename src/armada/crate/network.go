package crate

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"text/template"
)

type Network struct {
	Name              string
	AddressCidr       string
	Bridge            string
	BridgeIp          string
	TxQueueLen        int
	Mtu               int
	HostInterfaceName string
	TempVethPeerName  string
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
	"ip netns exec {{.Name}} ip addr add {{.AddressCidr}} dev eth0",
	"ip netns exec {{.Name}} ip link set eth0 up",
	"ip netns exec {{.Name}} ip route add default via {{.BridgeIp}} dev eth0",
}

// TODO: handle error and clean up after ourselves
func (self *Network) create() error {
	fmt.Println("IPCREATE: START")
	for _, cmd := range ipCommands {
		// template it
		var buf bytes.Buffer
		tmpl, err := template.New("ip").Parse(cmd)
		if err != nil {
			return err
		}
		err = tmpl.Execute(&buf, self)
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

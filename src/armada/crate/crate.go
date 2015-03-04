package main

import (
	"encoding/json"
	"fmt"
	//log "github.com/Sirupsen/logrus"
	"github.com/docker/libcontainer"
	"github.com/docker/libcontainer/configs"
	"io/ioutil"
	"os"
	"path/filepath"
)

var standardEnvironment = []string{
	"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
	"HOSTNAME=nsinit",
	"TERM=xterm",
}

type Crate struct {
	Root string
}

func New(root string) *Crate {
	return &Crate{root}
}

func (self *Crate) Create(id string, config *configs.Config) (libcontainer.Container, error) {

	// create
	factory, err := self.Factory()
	if err != nil {
		return nil, err
	}

	container, err := factory.Create(id, config)
	if err != nil {
		return nil, err
	}

	/*
		// lookup network namespace fd
		fd, err := self.LookupNetworkNsFd(&config.Namespaces)
		if err != nil {
			return nil, err
		}

		// setup network
		nc := config.Networks[0]
		network := &Network{
			NamespaceFd:       fd,
			Address:           nc.Address,
			Bridge:            nc.Bridge,
			Mtu:               nc.Mtu,
			TxQueueLen:        nc.TxQueueLen,
			HostInterfaceName: nc.HostInterfaceName,
		}
		err = network.netlinkCreate()
		if err != nil {
			return nil, err
		}
	*/

	nc := config.Networks[0]
	network := &Network{}
	err = network.ipCreate(ipContext{
		Name:        id,
		Bridge:      nc.Bridge,
		BridgeIp:    "10.4.0.255",
		AddressCidr: nc.Address,
	})
	if err != nil {
		return nil, err
	}

	config.Networks = nil

	// start a dummy process to force libcontainer to actually create shit
	process := &libcontainer.Process{
		Args:   []string{"/sbin/ip", "link"},
		Env:    standardEnvironment,
		User:   "",
		Cwd:    "",
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	err = container.Start(process)
	if err != nil {
		return container, err
	}

	// write out config data
	data, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	err = ioutil.WriteFile(filepath.Join(self.Root, id, "container.json"), data, 0655)
	return container, err
}

func (self *Crate) Load(id string) (libcontainer.Container, error) {
	factory, err := self.Factory()
	if err != nil {
		return nil, err
	}

	return factory.Load(id)
}

func (self *Crate) Factory() (libcontainer.Factory, error) {
	return libcontainer.New(self.Root, libcontainer.Cgroupfs)
}

func (self *Crate) LookupNetworkNsFd(n *configs.Namespaces) (uintptr, error) {
	for _, ns := range *n {
		if ns.Type == configs.NEWNET {
			file, err := os.Open(ns.GetPath(0))
			if err != nil {
				return 0, err
			}
			return file.Fd(), nil
		}
	}
	return 0, fmt.Errorf("Network namespace not found")
}

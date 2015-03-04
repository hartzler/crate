package crate

import (
	"encoding/json"
	"github.com/docker/libcontainer"
	"github.com/docker/libcontainer/configs"
	"io/ioutil"
	"os"
	"path/filepath"
)

var StandardEnvironment = []string{
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

func (self *Crate) SetupRoot() error {
	if _, err := os.Stat(self.Root); os.IsNotExist(err) {
		return os.MkdirAll(self.Root, 0700)
	}
	return nil
}

func (self *Crate) Factory() (libcontainer.Factory, error) {
	return libcontainer.New(self.Root, libcontainer.Cgroupfs)
}

func (self *Crate) Create(id string, config *configs.Config) (*Container, error) {

	// create
	factory, err := self.Factory()
	if err != nil {
		return nil, err
	}

	container, err := factory.Create(id, config)
	if err != nil {
		return nil, err
	}

	// setup our network namespace, links, and routes
	nc := config.Networks[0]
	network := &Network{
		Name:        id,
		Bridge:      nc.Bridge,
		BridgeIp:    nc.Gateway,
		AddressCidr: nc.Address,
	}
	if err = network.create(); err != nil {
		return nil, err
	}

	// clear config so libcontainer doesnt blow it
	config.Networks = nil

	// TODO: handle uncreated libcontainer
	// start a dummy process to force libcontainer to actually create shit
	process := &libcontainer.Process{
		Args:   []string{"/sbin/ip", "link"},
		Env:    StandardEnvironment,
		User:   "",
		Cwd:    "",
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	err = container.Start(process)
	if err != nil {
		return nil, err
	}

	// write out config data
	data, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	err = ioutil.WriteFile(filepath.Join(self.Root, id, "container.json"), data, 0655)
	return &Container{
		id:        id,
		crate:     self,
		container: container,
	}, err
}

func (self *Crate) Load(id string) (*Container, error) {
	factory, err := self.Factory()
	if err != nil {
		return nil, err
	}
	container, err := factory.Load(id)
	if err != nil {
		return nil, err
	}

	return &Container{
		id:        id,
		crate:     self,
		container: container,
	}, nil
}

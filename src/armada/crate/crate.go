package crate

import (
	"encoding/json"
	"github.com/docker/libcontainer"
	"github.com/docker/libcontainer/configs"
	"io/ioutil"
	"os"
	"path/filepath"
)

const CRATE_INIT = "crate-init"

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

	// setup our container/rootfs dir
	containerDir := filepath.Join(self.Root, id)
	/*
		err := os.MkdirAll(containerDir, 0700)
		if err != nil {
			return nil, err
		}
	*/

	// libcontainer create
	factory, err := self.Factory()
	if err != nil {
		return nil, err
	}

	// validate config
	container, err := factory.Create(id, config)
	if err != nil {
		return nil, err
	}

	// write out config data
	data, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	err = ioutil.WriteFile(filepath.Join(containerDir, "container.json"), data, 0644)

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

	// TODO switch to our own config which includes a libcontainer config
	// clear config as we handle outside of libcontainer
	config.Networks = nil

	// copy self to /crate-init in the container rootfs (non-portable hack)
	exePath, err := os.Readlink("/proc/self/exe")
	if err != nil {
		return nil, err
	}
	if _, err = copyFile(exePath, filepath.Join(config.Rootfs, CRATE_INIT)); err != nil {
		return nil, err
	}

	// start crate-init
	process := &libcontainer.Process{
		Args:   []string{"/" + CRATE_INIT, id},
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

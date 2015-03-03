package main

import (
	"encoding/json"
	//"fmt"
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

	// start a dummy process to force libcontainer to actually create
	process := &libcontainer.Process{
		Args:   []string{"/bin/date"},
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

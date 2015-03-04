package crate

import (
	"github.com/docker/libcontainer"
)

// small wrapper around libcontainer.Container
type Container struct {
	id        string
	crate     *Crate                 // our manager
	container libcontainer.Container // underlying container
}

func (self *Container) Start(process *libcontainer.Process) error {
	return self.container.Start(process)
}

func (self *Container) Destroy() error {
	err := self.container.Destroy()
	return err
}

func (self *Container) Pause() error {
	return self.container.Pause()
}

func (self *Container) Resume() error {
	return self.container.Resume()
}

func (self *Container) Processes() ([]int, error) {
	return self.container.Processes()
}

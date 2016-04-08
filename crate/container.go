package crate

import (
	"github.com/opencontainers/runc/libcontainer"
	"os"
)

// small wrapper around libcontainer.Container
type Container struct {
	id        string
	crate     *Crate                 // our manager
	container libcontainer.Container // underlying container
}

type Status struct {
	libcontainer.Status
	*libcontainer.State
	Pretty string
}

func (self *Container) Start(process *libcontainer.Process) error {
	return self.container.Start(process)
}

func (self *Container) Status() (*Status, error) {
	state, err := self.container.State()
	if err != nil {
		return nil, err
	}
	status, err := self.container.Status()
	if err != nil {
		return nil, err
	}
	pretty := "unknown"
	switch status {
	case libcontainer.Running:
		pretty = "running"
	case libcontainer.Pausing:
		pretty = "pausing"
	case libcontainer.Paused:
		pretty = "paused"
	case libcontainer.Destroyed:
		pretty = "destroyed"
	}

	return &Status{status, state, pretty}, nil
}

func (self *Container) Destroy() error {
	// can't call this as the state is wrong... sigh
	//return self.container.Destroy()

	// kill all processes in this container
	pids, err := self.Processes()
	if err != nil {
		return err
	}

	for _, pid := range pids {
		p, err := os.FindProcess(pid)
		if err != nil {
			return err
		}
		err = p.Kill()
		if err != nil {
			return err
		}
	}

	return nil
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

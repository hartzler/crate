package crate

import (
	"github.com/docker/libcontainer"
	"os"
)

type RunArgs struct {
	Id   string
	Args []string
	Env  []string
	User string
	Cwd  string
}

func (self *Crate) Run(args RunArgs) error {
	// lookup container
	container, err := self.Load(args.Id)
	if err != nil {
		return err
	}

	// start process
	return container.Start(&libcontainer.Process{
		Args:   args.Args,
		Env:    args.Env,
		User:   args.User,
		Cwd:    args.Cwd,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})
}

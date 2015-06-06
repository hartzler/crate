package crate

import (
	"armada/pkg/fd"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
)

type Process struct {
	Id    string
	Args  []string
	Env   []string
	User  string
	Cwd   string
	Shell bool
}

func (self *Crate) Run(id string, process Process) error {
	conn, err := net.Dial("unix", self.controlSocket(id))
	if err != nil {
		return err
	}
	defer conn.Close()

	// send process
	if err := json.NewEncoder(conn).Encode(&process); err != nil {
		return err
	}

	// send our fd's
	if process.Shell {
		return fd.Send(conn.(*net.UnixConn), os.Stdin, os.Stdout, os.Stderr)
	} else {
		// stdout/err
		stdout, err := os.Create(filepath.Join(self.path(id), process.Id+".stdout"))
		if err != nil {
			return err
		}
		stderr, err := os.Create(filepath.Join(self.path(id), process.Id+".stderr"))
		if err != nil {
			return err
		}
		return fd.Send(conn.(*net.UnixConn), stdout, stderr)
	}
}

func (self *Crate) Shell(id string) error {
	fmt.Println("Shell: loading...")

	if err := self.Run(id, Process{
		Args:  []string{"/bin/sh"},
		Shell: true,
	}); err != nil {
		return err
	}

	// TODO: coordinate exit lol
	// for now, just have to ctrl-c to exit.
	select {}

	// never reached...
	return nil
}

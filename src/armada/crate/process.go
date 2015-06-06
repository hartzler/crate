package crate

import (
	"armada/pkg/fd"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type RunArgs struct {
	Args  []string
	Env   []string
	User  string
	Cwd   string
	Stdio bool // flag if we are sending stdin/out/err after the json
}

func (self *Crate) Run(id string, args RunArgs) error {
	conn, err := net.Dial("unix", self.controlSocket(id))
	if err != nil {
		return err
	}
	defer conn.Close()

	return json.NewEncoder(conn).Encode(&args)
}

func (self *Crate) Shell(id string) error {
	fmt.Println("Shell: loading...")

	args := RunArgs{
		Args:  []string{"/bin/sh"},
		Stdio: true,
	}

	conn, err := net.Dial("unix", self.controlSocket(id))
	if err != nil {
		return err
	}
	defer conn.Close()

	// send run args
	if err := json.NewEncoder(conn).Encode(&args); err != nil {
		return err
	}

	// send our stdio
	if err := fd.Send(conn.(*net.UnixConn), os.Stdin, os.Stdout, os.Stderr); err != nil {
		return err
	}

	// TODO: coordinate exit lol
	// for now, just have to ctrl-c to exit.
	select {}

	// never reached...
	return nil
}

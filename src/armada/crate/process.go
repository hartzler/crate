package crate

import (
	"fmt"
	"net"
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
	/*
		container, err := self.Load(args.Id)
		if err != nil {
			return err
		}
	*/

	fmt.Println("OUTSIDE: Dialing...")
	conn, err := net.Dial("unix", "/home/vagrant/busybox/crate.socket")
	if err != nil {
		return err
	}

	fmt.Println("OUTSIDE: Writing...")
	if _, err := conn.Write([]byte("hello initer!\n")); err != nil {
		return err
	}

	fmt.Println("OUTSIDE: Hanging up...")
	if err = conn.Close(); err != nil {
		return err
	}

	fmt.Println("OUTSIDE: Finished.")
	return nil

}

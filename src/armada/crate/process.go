package crate

import (
	"encoding/json"
	"fmt"
	"net"
	"path/filepath"
)

type RunArgs struct {
	Args []string
	Env  []string
	User string
	Cwd  string
}

func (self *Crate) Run(id string, args RunArgs) error {
	fmt.Println("OUTSIDE: Dialing...")
	conn, err := net.Dial("unix", filepath.Join(self.containersRoot(), id, "rootfs", "crate.socket"))
	if err != nil {
		return err
	}

	fmt.Println("OUTSIDE: Writing...")

	if err := json.NewEncoder(conn).Encode(&args); err != nil {
		return err
	}

	fmt.Println("OUTSIDE: Hanging up...")
	if err = conn.Close(); err != nil {
		return err
	}

	fmt.Println("OUTSIDE: Finished.")
	return nil

}

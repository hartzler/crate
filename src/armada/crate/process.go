package crate

import (
	"encoding/json"
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

	conn, err := net.Dial("unix", filepath.Join(self.containersRoot(), id, "rootfs", "crate.socket"))
	if err != nil {
		return err
	}

	if err := json.NewEncoder(conn).Encode(&args); err != nil {
		return err
	}

	if err = conn.Close(); err != nil {
		return err
	}

	return nil

}

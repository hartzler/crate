package crate

import (
	"encoding/json"
	"github.com/armada-io/crate/pkg/fd"
	"io"
	"net"
	"os"
	"path/filepath"
)

const (
	RestartPolicyNever = iota
	RestartPolicyAlways
	RestartPolicyOnFailure
)

type Process struct {
	Id            string
	Args          []string
	Env           []string
	User          string
	Cwd           string
	RestartPolicy int
}

type ProcessEvent struct {
	Timestamp int64
	Type      string // running, terminated
	Pid       int
	ExitCode  int
	Error     string
}

type ProcessCmd struct {
	Type    string // run, stat, stop, kill
	Process Process
}

func (self *Crate) sendCmd(id string, kind string, process Process) (net.Conn, error) {
	conn, err := net.Dial("unix", self.controlSocket(id))
	if err != nil {
		return nil, err
	}
	// send process
	cmd := ProcessCmd{
		Type:    kind,
		Process: process,
	}
	if err := json.NewEncoder(conn).Encode(&cmd); err != nil {
		return nil, err
	}
	return conn, nil
}

func (self *Crate) Run(id string, process Process) error {
	conn, err := self.sendCmd(id, "run", process)
	if err != nil {
		return err
	}
	defer conn.Close()

	// stdout/err
	stdout, err := os.Create(filepath.Join(self.path(id), process.Id+".stdout"))
	if err != nil {
		return err
	}
	stderr, err := os.Create(filepath.Join(self.path(id), process.Id+".stderr"))
	if err != nil {
		return err
	}
	if err := fd.Send(conn.(*net.UnixConn), stdout, stderr); err != nil {
		return err
	}

	// read response
	bytes := make([]byte, 1024)
	if _, err = conn.Read(bytes); err != nil {
		if err != io.EOF {
			return err
		}
	}
	return nil
}

func (self *Crate) Shell(cid string) error {
	conn, err := self.sendCmd(cid, "shell", Process{
		Args: []string{"/bin/sh"},
	})
	if err != nil {
		return err
	}
	defer conn.Close()

	// send our fd's
	if err := fd.Send(conn.(*net.UnixConn), os.Stdin, os.Stdout, os.Stderr); err != nil {
		return err
	}
	// read response
	bytes := make([]byte, 1024)
	if _, err = conn.Read(bytes); err != nil {
		if err != io.EOF {
			return err
		}
	}
	return nil
}

func (self *Crate) Stat(cid, pid string) (string, error) {
	conn, err := self.sendCmd(cid, "stat", Process{Id: pid})
	if err != nil {
		return "", err
	}
	defer conn.Close()
	return "", nil
}

func (self *Crate) Stop(cid, pid string) error {
	conn, err := self.sendCmd(cid, "stop", Process{Id: pid})
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}

func (self *Crate) Kill(cid, pid string) error {
	conn, err := self.sendCmd(cid, "kill", Process{Id: pid})
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}

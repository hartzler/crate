// small process that is made to run inside a crate container as PID 1
// and control container processes.
package pid1

import (
	"armada/crate"
	"armada/pkg/fd"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/opencontainers/runc/libcontainer/nsenter"
	"net"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"
)

const MAX_SLEEP = 2 * time.Second

type pinfo struct {
	process crate.Process
	pid     int
}

var processMap = make(map[string]*pinfo)
var processMapLock sync.RWMutex

func setProcessMap(inf *pinfo) {
	processMapLock.Lock()
	processMap[inf.process.Id] = inf
	processMapLock.Unlock()
}

func getProcessMap(id string) *pinfo {
	processMapLock.RLock()
	defer processMapLock.RUnlock()
	return processMap[id]
}

func Start() {

	// clear if exists
	fmt.Println("[crate-init] starting...")
	fmt.Println("[crate-init] opening passed socket...")
	l, err := net.FileListener(os.NewFile(3, "listener"))
	if err != nil {
		fmt.Println("[crate-init] err:", err)
		panic(err)
	}
	defer l.Close()

	fmt.Println("[crate-init] ready...")
	for {
		conn, err := l.Accept()
		if err != nil {
			fatal(err)
			continue // try again...
		}
		go func() {
			if err := handle(conn.(*net.UnixConn)); err != nil {
				fatal(err)
			}
			fmt.Println("[crate-init] finished command.")
		}()
	}
	fmt.Println("[crate-init] exiting...")
}

func handle(conn *net.UnixConn) error {
	defer conn.Close()
	var cmd crate.ProcessCmd
	if err := json.NewDecoder(conn).Decode(&cmd); err != nil {
		return err
	}
	fmt.Println("[crate-init] DECODED:", cmd)

	// handle type of command
	switch cmd.Type {
	case "run":
		return run(conn, cmd.Process)
	case "shell":
		return shell(conn, cmd.Process)
	case "stat":
		return stat(conn, cmd.Process.Id)
	case "stop":
		return stop(conn, cmd.Process.Id)
	case "kill":
		return kill(conn, cmd.Process.Id)
	}
	// unknown
	return errors.New("Unknown process command: " + cmd.Type)
}

// run process
func run(conn *net.UnixConn, process crate.Process) error {
	// read stdio fd's from socket
	var stdin, stdout, stderr *os.File
	files, err := fd.Receive(conn, 2, []string{"/dev/stdout", "/dev/stdin"})
	if err != nil {
		return fmt.Errorf("[crate-init] Error reading 2 fd's from socket: %s", err)
	}
	stdout = files[0]
	stderr = files[1]
	fmt.Println("[crate-init] closing connection.")
	conn.Close()

	exp := time.Millisecond
	for {
		cmd := exec.Command(process.Args[0], process.Args[1:]...)
		cmd.Env = process.Env
		cmd.Stdin = stdin
		cmd.Stdout = stdout
		cmd.Stderr = stderr
		if err := cmd.Start(); err != nil {
			return err
		}
		setProcessMap(&pinfo{
			process: process,
			pid:     cmd.Process.Pid,
		})
		err := cmd.Wait()
		switch process.RestartPolicy {
		case crate.RestartPolicyNever:
			if err != nil {
				return fmt.Errorf("Error running command: %s %s", cmd, err)
			}
			return nil
		case crate.RestartPolicyOnFailure:
			if err == nil {
				return nil
			}
			// error, so retry
		case crate.RestartPolicyAlways:
			// always retry
		}
		if err != nil {
			fmt.Println("[crate-init] Process exited with error: %s %s", cmd, err)
		}
		fmt.Println("[crate-init] Restarting process: %s", cmd)
		// exponential backoff
		exp = exp << 2
		if exp > MAX_SLEEP {
			exp = MAX_SLEEP
		}
		time.Sleep(exp)
	}
	panic("UNREACHABLE CODE!")
	return nil
}

// use parent passed std in/out/err and don't close conn until cmd is done
func shell(conn *net.UnixConn, process crate.Process) error {
	defer conn.Close()
	files, err := fd.Receive(conn, 3, []string{"/dev/stdin", "/dev/stdout", "/dev/stdin"})
	if err != nil {
		return fmt.Errorf("[crate-init] Error reading 3 fd's from socket: %s", err)
	}
	cmd := exec.Command(process.Args[0], process.Args[1:]...)
	cmd.Env = process.Env
	cmd.Stdin = files[0]
	cmd.Stdout = files[1]
	cmd.Stderr = files[2]
	return cmd.Run()
}

// query process state
// if no error, its running, else see error
func stat(conn *net.UnixConn, id string) error {
	pinf := processMap[id]
	if pinf == nil {
		// unknown
		return errors.New("Invalid id")
	}
	proc, err := os.FindProcess(pinf.pid)
	if err != nil {
		return err
	}
	return proc.Signal(syscall.Signal(0))
}

// stop process (send SIGINT)
func stop(conn *net.UnixConn, id string) error {
	pinf := processMap[id]
	if pinf == nil {
		// unknown
		return errors.New("Invalid id")
	}
	proc, err := os.FindProcess(pinf.pid)
	if err != nil {
		return err
	}
	if err := proc.Signal(os.Interrupt); err != nil {
		return err
	}
	return nil
}

// kill process
func kill(conn *net.UnixConn, id string) error {
	pinf := processMap[id]
	if pinf == nil {
		// unknown
		return errors.New("Invalid id")
	}
	proc, err := os.FindProcess(pinf.pid)
	if err != nil {
		return err
	}
	if err := proc.Kill(); err != nil {

	}
	return nil
}

func fatal(err error) {
	fmt.Println("[crate-init] [ERROR]", err)
}

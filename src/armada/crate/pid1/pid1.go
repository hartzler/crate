// small process that is made to run inside a crate container as PID 1
// and control container processes.
package pid1

import (
	"armada/crate"
	"armada/pkg/fd"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"time"
)

const MAX_SLEEP = 2 * time.Second

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
			defer conn.Close()
			if err := handle(conn.(*net.UnixConn)); err != nil {
				fatal(err)
			}
			fmt.Println("[crate-init] finished command.")
		}()
	}
	fmt.Println("[crate-init] exiting...")
}

func handle(conn *net.UnixConn) error {
	var process crate.Process
	if err := json.NewDecoder(conn).Decode(&process); err != nil {
		return err
	}
	fmt.Println("[crate-init] DECODED:", process)

	// read stdio fd's from socket
	var stdin, stdout, stderr *os.File
	if process.Shell {
		files, err := fd.Receive(conn, 3, []string{"/dev/stdin", "/dev/stdout", "/dev/stdin"})
		if err != nil {
			return fmt.Errorf("Error reading 3 fd's from socket: %s", err)
		}
		stdin = files[0]
		stdout = files[1]
		stderr = files[2]
	} else {
		files, err := fd.Receive(conn, 2, []string{"/dev/stdout", "/dev/stdin"})
		if err != nil {
			return fmt.Errorf("Error reading 2 fd's from socket: %s", err)
		}
		stdout = files[0]
		stderr = files[1]
	}

	exp := time.Millisecond
	for {
		cmd := exec.Command(process.Args[0], process.Args[1:]...)
		cmd.Env = process.Env
		cmd.Stdin = stdin
		cmd.Stdout = stdout
		cmd.Stderr = stderr
		err := cmd.Run()
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

func fatal(err error) {
	fmt.Println("[crate-init] [ERROR]", err)
}

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
)

func Start() {

	// clear if exists
	fmt.Println("PID1: starting...")
	fmt.Println("PID1: opening passed socket...")
	l, err := net.FileListener(os.NewFile(3, "listener"))
	if err != nil {
		fmt.Println("PID1: err:", err)
		panic(err)
	}
	defer l.Close()

	fmt.Println("PID1: ready...")
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
			fmt.Println("PID1: finished command.")
		}()
	}
	fmt.Println("PID1: exiting...")
}

func handle(conn *net.UnixConn) error {
	var process crate.Process
	if err := json.NewDecoder(conn).Decode(&process); err != nil {
		return err
	}
	fmt.Printf("PID1: DECODED: %s\n", process)

	cmd := exec.Command(process.Args[0], process.Args[1:]...)
	cmd.Env = process.Env

	// read stdio fd's from socket
	if process.Shell {
		files, err := fd.Receive(conn, 3, []string{"/dev/stdin", "/dev/stdout", "/dev/stdin"})
		if err != nil {
			return err
		}
		cmd.Stdin = files[0]
		cmd.Stdout = files[1]
		cmd.Stderr = files[2]
	} else {
		files, err := fd.Receive(conn, 2, []string{"/dev/stdout", "/dev/stdin"})
		if err != nil {
			return err
		}
		cmd.Stdout = files[0]
		cmd.Stderr = files[1]
	}

	return cmd.Run()
}

func fatal(err error) {
	fmt.Println("PID1: [ERROR]", err)
}

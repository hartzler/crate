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
	var args crate.RunArgs
	if err := json.NewDecoder(conn).Decode(&args); err != nil {
		return err
	}
	fmt.Printf("PID1: DECODED: %s\n", args)

	cmd := exec.Command(args.Args[0], args.Args[1:]...)
	cmd.Env = args.Env

	// check if stdin/out/err are being sent
	if args.Stdio {
		files, err := fd.Receive(conn, 3, []string{"/dev/stdin", "/dev/stdout", "/dev/stdin"})
		if err != nil {
			return err
		}
		cmd.Stdin = files[0]
		cmd.Stdout = files[1]
		cmd.Stderr = files[2]
	} else {
		// just use our own outputs
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	return cmd.Run()
}

func fatal(err error) {
	fmt.Println("PID1: [ERROR]", err)
}

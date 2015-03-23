// small process that is made to run inside a crate container as PID 1
// and control container processes.
package initer

import (
	"armada/crate"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"syscall"
)

var SocketFile = "crate.socket"

func CrateInit() {

	// clear if exists
	os.Remove("/" + SocketFile)

	// set permission on listen to 0600...
	oldmask := syscall.Umask(0177)
	l, err := net.Listen("unix", "/"+SocketFile)
	if err != nil {
		panic(err)
	}
	defer l.Close()

	// restore
	syscall.Umask(oldmask)

	fmt.Println("CRATE-INIT: ready...")
	for {
		conn, err := l.Accept()
		if err != nil {
			fatal(err)
			continue // try again...
		}
		go func() {
			defer conn.Close()
			if err := handle(conn); err != nil {
				fatal(err)
			}
		}()
	}
}

func handle(conn net.Conn) error {
	var args crate.RunArgs
	if err := json.NewDecoder(conn).Decode(&args); err != nil {
		return err
	}
	fmt.Printf("CRATE-INIT: DECODED: %s\n", args)

	// just inherit out/err for now...
	cmd := exec.Command(args.Args[0], args.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func fatal(err error) {
	fmt.Println("CRATE-INIT: [ERROR]", err)
}

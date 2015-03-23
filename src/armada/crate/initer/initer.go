// small process that is made to run inside a crate container and control
// container processes.
package initer

import (
	"armada/crate"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
)

var SocketFile = "crate.socket"

func CrateInit() {

	// clear if exists
	os.Remove("/" + SocketFile)
	l, err := net.Listen("unix", "/"+SocketFile)
	if err != nil {
		panic(err)
	}
	defer l.Close()

	fmt.Println("CRATE-INIT: hello from init!")

	for {
		conn, err := l.Accept()
		if err != nil {
			fatal(err)
			continue
		}
		var buf [1024]byte
		n, err := conn.Read(buf[:])
		if err != nil {
			fatal(err)
			continue
		}
		fmt.Printf("CRATE-INIT: %s\n", string(buf[:n]))

		var args crate.RunArgs
		if err = json.Unmarshal(buf[:n], &args); err != nil {
			fatal(err)
			continue
		}
		fmt.Printf("CRATE-INIT: DECODED: %s\n", args)

		cmd := exec.Command(args.Args[0], args.Args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		if err := cmd.Run(); err != nil {
			fatal(err)
			continue
		}
		//fmt.Println("CRATE-INIT: [RUN]: ", string(bytes))

		conn.Close()
	}
}

func fatal(err error) {
	fmt.Println("CRATE-INIT: [ERROR]", err)
}

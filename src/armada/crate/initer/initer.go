// small process that is made to run inside a crate container and control
// container processes.
package initer

import (
	"fmt"
	"net"
	"os"
	//"os/exec"
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

	fmt.Println("hello from init!")

	for {
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}
		var buf [1024]byte
		n, err := conn.Read(buf[:])
		if err != nil {
			panic(err)
		}
		fmt.Printf("INSIDE: %s\n", string(buf[:n]))
		conn.Close()
	}
}

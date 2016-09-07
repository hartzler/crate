// cli commands
package command

import (
	"fmt"
	"github.com/armada-io/crate/crate"
	"github.com/codegangsta/cli"
	"os"
)

var Commands []cli.Command

func FromContext(context *cli.Context) *crate.Crate {
	return crate.New(context.GlobalString("root"))
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "ERROR:", err)
	os.Exit(1)
}

func fatalf(t string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, t, v...)
	os.Exit(1)
}

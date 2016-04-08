// cli commands
package command

import (
	"fmt"
	"github.com/armada-io/crate/crate"
	"github.com/codegangsta/cli"
	"github.com/opencontainers/runc/libcontainer"
	"os"
)

var Commands []cli.Command

func FromContext(context *cli.Context) *crate.Crate {
	return crate.New(context.GlobalString("root"))
}

// hack TODO: rename all from -> From
func fromContext(context *cli.Context) *crate.Crate {
	return crate.New(context.GlobalString("root"))
}

func getContainer(context *cli.Context, id string) (*crate.Container, error) {
	return fromContext(context).Load(id)
}

func fatal(err error) {
	if lerr, ok := err.(libcontainer.Error); ok {
		lerr.Detail(os.Stderr)
		os.Exit(1)
	}
	fmt.Fprintln(os.Stderr, "ERROR:", err)
	os.Exit(1)
}

func fatalf(t string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, t, v...)
	os.Exit(1)
}

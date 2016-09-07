package command

import (
	"fmt"
	"github.com/codegangsta/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:        "install",
		Usage:       "Installs a crate",
		Description: "args: <id>",
		Flags: []cli.Flag{
			cli.StringSliceFlag{Name: "cargo", Value: &cli.StringSlice{}, Usage: "urls of cargo files that make up the crate"},
		},
		Action: func(context *cli.Context) {
			args := context.Args()
			if len(args) != 1 {
				fatal(fmt.Errorf("expected argument <id>"))
			}

			id := args[0]
			crate := FromContext(context)
			if err := crate.Install(id, []string(context.StringSlice("cargo"))); err != nil {
				fatal(err)
			}
		},
	})
}

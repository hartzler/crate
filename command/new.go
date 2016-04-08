package command

import (
	"fmt"
	"github.com/codegangsta/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:        "new",
		Usage:       "creates a container",
		Description: "args: <id>",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "read-only", Usage: "set the container's rootfs as read-only"},
			cli.StringSliceFlag{Name: "cargo", Value: &cli.StringSlice{}, Usage: "urls of cargo files to unpack in rootfs"},
			cli.StringSliceFlag{Name: "bind", Value: &cli.StringSlice{}, Usage: "add bind mounts to the container"},
			cli.StringSliceFlag{Name: "tmpfs", Value: &cli.StringSlice{}, Usage: "add tmpfs mounts to the container"},
			cli.IntFlag{Name: "cpushares", Usage: "set the cpushares for the container"},
			cli.IntFlag{Name: "memory-limit", Usage: "set the memory limit for the container"},
			cli.IntFlag{Name: "memory-swap", Usage: "set the memory swap limit for the container"},
			cli.StringFlag{Name: "cpuset-cpus", Usage: "set the cpuset cpus"},
			cli.StringFlag{Name: "cpuset-mems", Usage: "set the cpuset mems"},
			cli.IntFlag{Name: "userns-root-uid", Usage: "set the user namespace root uid"},
			cli.StringFlag{Name: "hostname", Value: "crate", Usage: "hostname value for the container"},
			cli.StringFlag{Name: "bridge", Usage: "name of bridge interface"},
			cli.StringFlag{Name: "address", Usage: "ip/cidr address"},
			cli.StringFlag{Name: "gateway", Value: "10.4.0.255", Usage: "container gateway address"},
			cli.IntFlag{Name: "mtu", Value: 1500, Usage: "veth mtu"},
			cli.IntFlag{Name: "txq", Value: 200, Usage: "veth tx queue length"},
		},
		Action: func(context *cli.Context) {

			args := context.Args()
			if len(args) != 1 {
				fatal(fmt.Errorf("expected argument <id>"))
			}

			// id
			id := args[0]

			// libconfig
			libconfig := getTemplate(id)
			modify(libconfig, context)
			fmt.Println(libconfig)

			crate := fromContext(context)
			container, err := crate.Create(id, []string(context.StringSlice("cargo")), libconfig)
			if err != nil {
				fatal(err)
			}
			pids, err := container.Processes()
			if err != nil {
				fatal(err)
			}
			// return PID1 host pid
			fmt.Println(pids[0])
		},
	})
}

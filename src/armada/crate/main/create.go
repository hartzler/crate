package main

import (
	"fmt"
	"github.com/codegangsta/cli"
)

var CreateCommand = cli.Command{
	Name:  "create",
	Usage: "creates a new container",
	Flags: []cli.Flag{
		cli.StringFlag{Name: "id", Usage: "specify the ID for a container"},
		cli.BoolFlag{Name: "read-only", Usage: "set the container's rootfs as read-only"},
		cli.StringSliceFlag{Name: "bind", Value: &cli.StringSlice{}, Usage: "add bind mounts to the container"},
		cli.StringSliceFlag{Name: "tmpfs", Value: &cli.StringSlice{}, Usage: "add tmpfs mounts to the container"},
		cli.IntFlag{Name: "cpushares", Usage: "set the cpushares for the container"},
		cli.IntFlag{Name: "memory-limit", Usage: "set the memory limit for the container"},
		cli.IntFlag{Name: "memory-swap", Usage: "set the memory swap limit for the container"},
		cli.StringFlag{Name: "cpuset-cpus", Usage: "set the cpuset cpus"},
		cli.StringFlag{Name: "cpuset-mems", Usage: "set the cpuset mems"},
		cli.StringFlag{Name: "apparmor-profile", Usage: "set the apparmor profile"},
		cli.StringFlag{Name: "process-label", Usage: "set the process label"},
		cli.StringFlag{Name: "mount-label", Usage: "set the mount label"},
		cli.StringFlag{Name: "rootfs", Usage: "set the rootfs"},
		cli.IntFlag{Name: "userns-root-uid", Usage: "set the user namespace root uid"},
		cli.StringFlag{Name: "hostname", Value: "crate", Usage: "hostname value for the container"},
		cli.StringFlag{Name: "bridge", Value: "armada0", Usage: "name of bridge interface"},
		cli.StringFlag{Name: "address", Usage: "ip/cidr address"},
		cli.StringFlag{Name: "gateway", Value: "10.4.0.255", Usage: "container gateway address"},
		cli.IntFlag{Name: "mtu", Value: 1500, Usage: "veth mtu"},
		cli.IntFlag{Name: "txq", Value: 200, Usage: "veth tx queue length"},
	},
	Action: func(context *cli.Context) {
		id := context.String("id")
		config := getTemplate(id)
		modify(config, context)

		fmt.Println(config)

		crate := fromContext(context)
		if _, err := crate.Create(id, config); err != nil {
			fatal(err)
		}
	},
}

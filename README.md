Crate is a libcontainer based container manager.  It is intended to be used in conjunction with an overlay network to make containers from different hosts act as if they are on the same virtual network.  Container IPs are intended to be stable and unique across the network, so containers can be moved around by only adjusting routing entires.

For more details of the networking see: network.md

## Usage

    NAME:
       crate - manage containers and connections

    USAGE:
       crate [global options] command [command options] [arguments...]

    VERSION:
       0.1

    AUTHOR:
      Author - <matt@armada.io>

    COMMANDS:
       setup	create the network bridge [bridge-name]
       create	creates a new container
       run		start a process inside a container
       destroy	destroy the container
       pause	pause the container's processes
       resume	resume the container's processes
       pids		list the pids of a container
       status	show the status of a container
       help, h	Shows a list of commands or help for one command

    GLOBAL OPTIONS:
       --root "/var/lib/crate/containers"	root directory for containers
       --log-file 				set the log file to output logs to
       --debug				enable debug output in the logs
       --help, -h				show help
       --version, -v			print the version

## Crate Image

Applications are packaged into a "crate" file.

## Crate Container (Runtime Environment)

The environment that the crate processes run in, or the "container".

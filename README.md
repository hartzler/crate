Crate is simple chroot based package manager.

## Usage

    NAME:
       crate - simple read only chroot based package manager

    USAGE:
       crate [global options] command [command options] [arguments...]

    VERSION:
       0.1

    AUTHOR(S):
       Matt Hartzler <matt@armada.io>

    COMMANDS:
         install  Installs a crate
         remove   uninstalls a crate
         help, h  Shows a list of commands or help for one command

    GLOBAL OPTIONS:
       --root value      root directory for crate state (default: "/var/lib/crate")
       --log-file value  set the log file to output logs to
       --debug           enable debug output in the logs
       --help, -h        show help
       --version, -v     print the version


## Hacking

To get started, use vagrant and build / run crate as root:

    vagrant up
    vagrant ssh
    sudo -i
    cd /vagrant
    ./build.sh

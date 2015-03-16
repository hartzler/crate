crate manages containers.  a container is a directory tree containing image dirs
and the files needed to run and manage it.


# outside PID# for PID 1 of this container
# if file is missing or pid doesn't exist container is not running
/init.pid

# the unix domain socket for communicating with PID1
/init.socket

# rootfs for the container
/rootfs

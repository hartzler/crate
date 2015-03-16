#!/bin/sh
NAME=crate

# provision
apt-add-repository ppa:zfs-native/daily
apt-get update
apt-get install -y ubuntu-zfs
apt-get install -y lxc

# pool
mkdir -p /var/lib/$NAME
fallocate -l 10G /var/lib/$NAME/zfs.image
zpool create $NAME /var/lib/$NAME/zfs.image
zfs set dedup=on $NAME
zfs set compression=on $NAME
zpool set listsnapshots=on $NAME

# lxc
zfs create $NAME/lxc
zfs create $NAME/lxc/containers
cat >/etc/lxc/lxc.conf <<EOF
lxc.lxcpath = /$NAME/lxc/containers
lxc.bdev.zfs.root = $NAME/lxc/containers
EOF


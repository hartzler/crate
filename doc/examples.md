crate new --cargo=75e4e3de4cbbe5305c88cd7beff5ff82c40a0ba4cedb1a3c1150bcb896aba832 --cargo=10854cdf93f16e94d9d1b52ac1b89da8bbcb2d30d7bb2e2f6d64de65a658fdeb es-5

mount -t overlay -o lowerdir=/var/lib/crate/cargo/75e4e3de4cbbe5305c88cd7beff5ff82c40a0ba4cedb1a3c1150bcb896aba832:/var/lib/crate/cargo/10854cdf93f16e94d9d1b52ac1b89da8bbcb2d30d7bb2e2f6d64de65a658fdeb,upperdir=/var/lib/crate/containers/es-5/mounts/content,workdir=/var/lib/crate/containers/es-5/mounts/work overlay /var/lib/crate/containers/es-5/rootfs

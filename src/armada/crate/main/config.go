package main

import (
	"fmt"
	"math"
	"strings"
	"syscall"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/libcontainer/configs"
	"github.com/docker/libcontainer/utils"
)

const defaultMountFlags = syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV

func modify(config *configs.Config, context *cli.Context) {
	id := context.String("id")
	rootfs := context.String("rootfs")

	//config.ParentDeathSignal = context.Int("parent-death-signal")
	config.Readonlyfs = context.Bool("read-only")

	config.Cgroups = &configs.Cgroup{
		Name:            id,
		Parent:          "crate",
		AllowAllDevices: false,
		AllowedDevices:  configs.DefaultAllowedDevices,
		CpusetCpus:      context.String("cpuset-cpus"),
		CpusetMems:      context.String("cpuset-mems"),
		CpuShares:       int64(context.Int("cpushares")),
		Memory:          int64(context.Int("memory-limit")),
		MemorySwap:      int64(context.Int("memory-swap")),
	}

	config.AppArmorProfile = context.String("apparmor-profile")
	config.ProcessLabel = context.String("process-label")
	config.MountLabel = context.String("mount-label")

	config.Rootfs = rootfs

	userns_uid := context.Int("userns-root-uid")
	if userns_uid != 0 {
		//config.Namespaces.Add(configs.NEWUSER, fmt.Sprintf("/crate/%s/user", id))
		config.Namespaces.Add(configs.NEWUSER, "")
		config.UidMappings = []configs.IDMap{
			{ContainerID: 0, HostID: userns_uid, Size: 1},
			{ContainerID: 1, HostID: 1, Size: userns_uid - 1},
			{ContainerID: userns_uid + 1, HostID: userns_uid + 1, Size: math.MaxInt32 - userns_uid},
		}
		config.GidMappings = []configs.IDMap{
			{ContainerID: 0, HostID: userns_uid, Size: 1},
			{ContainerID: 1, HostID: 1, Size: userns_uid - 1},
			{ContainerID: userns_uid + 1, HostID: userns_uid + 1, Size: math.MaxInt32 - userns_uid},
		}
		for _, node := range config.Devices {
			node.Uid = uint32(userns_uid)
			node.Gid = uint32(userns_uid)
		}
	}
	for _, rawBind := range context.StringSlice("bind") {
		mount := &configs.Mount{
			Device: "bind",
			Flags:  syscall.MS_BIND | syscall.MS_REC,
		}
		parts := strings.SplitN(rawBind, ":", 3)
		switch len(parts) {
		default:
			logrus.Fatalf("invalid bind mount %s", rawBind)
		case 2:
			mount.Source, mount.Destination = parts[0], parts[1]
		case 3:
			mount.Source, mount.Destination = parts[0], parts[1]
			switch parts[2] {
			case "ro":
				mount.Flags |= syscall.MS_RDONLY
			case "rw":
			default:
				logrus.Fatalf("invalid bind mount mode %s", parts[2])
			}
		}
		config.Mounts = append(config.Mounts, mount)
	}
	for _, tmpfs := range context.StringSlice("tmpfs") {
		config.Mounts = append(config.Mounts, &configs.Mount{
			Device:      "tmpfs",
			Destination: tmpfs,
			Flags:       syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV,
		})
	}
	// todo: bind mount other namespaces into container dir
	config.Namespaces = configs.Namespaces([]configs.Namespace{
		{Type: configs.NEWNS},
		{Type: configs.NEWUTS},
		{Type: configs.NEWIPC},
		{Type: configs.NEWPID},
		{Type: configs.NEWNET},
		//{Type: configs.NEWNET, Path: fmt.Sprintf("/var/run/netns/%s", id)},
	})
	if bridge := context.String("bridge"); bridge != "" {
		hostName, err := utils.GenerateRandomName("armada", 7)
		if err != nil {
			logrus.Fatal(err)
		}
		network := &configs.Network{
			Type:              "veth",
			Name:              "eth0",
			Bridge:            bridge,
			Address:           context.String("address"),
			Gateway:           context.String("gateway"),
			Mtu:               context.Int("mtu"),
			TxQueueLen:        context.Int("txq"),
			HostInterfaceName: hostName,
		}
		config.Networks = append(config.Networks, network)
		fmt.Println(config.Networks[0])
	}
}

func getTemplate(id string) *configs.Config {
	return &configs.Config{
		ParentDeathSignal: int(syscall.SIGKILL),
		Capabilities: []string{
			"CHOWN",
			"DAC_OVERRIDE",
			"FSETID",
			"FOWNER",
			"MKNOD",
			"NET_RAW",
			"SETGID",
			"SETUID",
			"SETFCAP",
			"SETPCAP",
			"NET_BIND_SERVICE",
			"SYS_CHROOT",
			"KILL",
			"AUDIT_WRITE",
		},
		Devices: configs.DefaultAutoCreatedDevices,
		MaskPaths: []string{
			"/proc/kcore",
		},
		ReadonlyPaths: []string{
			"/proc/sys", "/proc/sysrq-trigger", "/proc/irq", "/proc/bus",
		},
		Mounts: []*configs.Mount{
			{
				Device:      "tmpfs",
				Source:      "shm",
				Destination: "/dev/shm",
				Data:        "mode=1777,size=65536k",
				Flags:       defaultMountFlags,
			},
			{
				Source:      "mqueue",
				Destination: "/dev/mqueue",
				Device:      "mqueue",
				Flags:       defaultMountFlags,
			},
			{
				Source:      "sysfs",
				Destination: "/sys",
				Device:      "sysfs",
				Flags:       defaultMountFlags | syscall.MS_RDONLY,
			},
		},
		Rlimits: []configs.Rlimit{
			{
				Type: syscall.RLIMIT_NOFILE,
				Hard: 1024,
				Soft: 1024,
			},
		},
	}

}

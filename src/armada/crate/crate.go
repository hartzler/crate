package crate

import (
	"encoding/json"
	"fmt"
	"github.com/opencontainers/runc/libcontainer"
	"github.com/opencontainers/runc/libcontainer/configs"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

const (
	CRATE_INIT       = "crate-init"
	PID_FILE         = "init.pid"
	LIBCONTAINER_DIR = "libcontainer"
)

var StandardEnvironment = []string{
	"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
	"HOSTNAME=nsinit",
	"TERM=xterm",
}

type Crate struct {
	Root string
}

func New(root string) *Crate {
	return &Crate{
		Root: root,
	}
}

func (self *Crate) containersRoot() string {
	return filepath.Join(self.Root, "containers")
}

func (self *Crate) controlSocket(id string) string {
	return filepath.Join(self.containersRoot(), id, "crate.socket")
}

func (self *Crate) SetupRoot() error {
	containersDir := self.containersRoot()
	if _, err := os.Stat(containersDir); os.IsNotExist(err) {
		return os.MkdirAll(containersDir, 0700)
	}
	return nil
}

func (self *Crate) factory(id string) (libcontainer.Factory, error) {
	containerDir := filepath.Join(self.containersRoot(), id)
	return libcontainer.New(containerDir, libcontainer.Cgroupfs)
}

func (self *Crate) path(id string) string {
	return filepath.Join(self.containersRoot(), id)
}

func (self *Crate) mounts(id string) string {
	return filepath.Join(self.path(id), "mounts")
}

func (self *Crate) rootfs(id string) string {
	return filepath.Join(self.path(id), "rootfs")
}

func (self *Crate) pidfile(id string) string {
	return filepath.Join(self.path(id), PID_FILE)
}

func (self *Crate) pid(id string) (int, error) {
	file, err := os.Open(self.pidfile(id))
	if err != nil {
		return 0, err
	}
	defer file.Close()
	pid, err := ioutil.ReadAll(file)
	if err != nil {
		return 0, err
	}
	fmt.Println("[crate] pid: ", string(pid))
	return strconv.Atoi(string(pid))
}

// TODO switch to our own config which includes options, cargo, and a libcontainer configs.Config
func (self *Crate) Create(id string, cargo []string, libconfig *configs.Config) (*Container, error) {

	// setup our container/rootfs dir
	containerDir := self.path(id)
	rootfs := self.rootfs(id)
	libconfig.Rootfs = rootfs
	if err := self.setupRootfs(id, rootfs, cargo); err != nil {
		return nil, err
	}

	// libcontainer create
	container, err := setupLibcontainer(id, containerDir, libconfig)
	if err != nil {
		return nil, err
	}

	// start crate-init
	fmt.Println("[crate] starting init...")
	if _, err := self.startInit(id, containerDir, container); err != nil {
		return nil, err
	}

	return &Container{
		id:        id,
		crate:     self,
		container: container,
	}, nil

}

func setupLibcontainer(id, containerDir string, libconfig *configs.Config) (libcontainer.Container, error) {
	// write out config data
	data, err := json.Marshal(libconfig)
	if err != nil {
		return nil, err
	}
	if err := ioutil.WriteFile(filepath.Join(containerDir, "config.json"), data, 0600); err != nil {
		return nil, err
	}

	// create container directory (will die if it already exists... lame)
	lcDir := filepath.Join(containerDir, LIBCONTAINER_DIR)
	var b bool
	if b, err = exists(lcDir); b {
		err = os.RemoveAll(lcDir)
	}
	if err != nil {
		return nil, err
	}

	// get absolute path our ourselves
	exePath, err := os.Readlink("/proc/self/exe")
	if err != nil {
		return nil, err
	}
	// ugly hack to get libcontainer to use absolute path...
	factory, err := libcontainer.New(containerDir, libcontainer.InitPath(exePath, exePath, "init"))
	if err != nil {
		return nil, err
	}
	container, err := factory.Create(LIBCONTAINER_DIR, libconfig)
	if err != nil {
		return nil, err
	}

	return container, nil
}

// TODO: cleanup all mounts/dirs on error
func (self *Crate) setupRootfs(id, rootfs string, cargo []string) error {
	// setup merge path for final rootfs
	if err := os.MkdirAll(rootfs, 0700); err != nil {
		return err
	}

	// setup cargo mounts dir
	mountsDir := self.mounts(id)
	fmt.Println("[crate] creating mounts dir: ", mountsDir)
	if err := os.MkdirAll(mountsDir, 0700); err != nil {
		return err
	}

	// create work dir
	workDir := filepath.Join(mountsDir, "work")
	if err := os.MkdirAll(workDir, 0700); err != nil {
		return err
	}

	// create write dir
	writeDir := filepath.Join(mountsDir, "content")
	if err := os.MkdirAll(writeDir, 0700); err != nil {
		return err
	}

	// get list of cargo mounts
	lowerDirs := []string{}
	for _, path := range cargo {
		lowerDirs = append(lowerDirs, filepath.Join(self.Root, "cargo", path))
	}

	// mount overlay
	opts := fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", strings.Join(lowerDirs, ":"), writeDir, workDir)
	args := []string{"-t", "overlay", "-o", opts, "overlay", rootfs}
	fmt.Println("[crate] mounting merged rootfs: mount", strings.Join(args, " "))
	if out, err := exec.Command("mount", args...).CombinedOutput(); err != nil {
		fmt.Println(string(out))
		return err
	}
	return nil
}
func (self *Crate) startInit(id, containerDir string, container libcontainer.Container) (int, error) {
	// copy self to /crate-init in the container rootfs (non-portable hack?)
	fmt.Println("[crate] DEBUG copying crate-init...")
	exePath, err := os.Readlink("/proc/self/exe")
	if err != nil {
		return -1, err
	}
	if _, err = copyFile(exePath, filepath.Join(containerDir, "rootfs", CRATE_INIT)); err != nil {
		return -2, err
	}

	// stdout/err
	fmt.Println("[crate] DEBUG creating stdio files...")
	stdout, err := os.Create(filepath.Join(containerDir, "stdout"))
	if err != nil {
		return -3, err
	}
	stderr, err := os.Create(filepath.Join(containerDir, "stderr"))
	if err != nil {
		return -4, err
	}

	// set permission on listen to 0600...
	fmt.Println("[crate] DEBUG syscall umask 0177...")
	oldmask := syscall.Umask(0177)

	// ctrl socket file to pass to our crate-init
	fmt.Println("[crate] DEBUG opening control socket...")
	socketPath := self.controlSocket(id)
	if err := os.Remove(socketPath); !os.IsNotExist(err) {
		return -1, err
	}
	fmt.Println("[crate] DEBUG open control socket listener...")
	socket, err := net.ListenUnix("unix", &net.UnixAddr{socketPath, "unix"})
	if err != nil {
		return -1, err
	}
	fmt.Println("[crate] DEBUG get fd to listener...")
	socketFile, err := socket.File()
	if err != nil {
		return -1, err
	}

	// restore
	fmt.Println("[crate] DEBUG restore syscall umask...")
	syscall.Umask(oldmask)

	// create init process
	process := &libcontainer.Process{
		Args:       []string{"/" + CRATE_INIT, id},
		Env:        StandardEnvironment,
		User:       "",
		Cwd:        "/",
		Stdin:      nil,
		Stdout:     stdout,
		Stderr:     stderr,
		ExtraFiles: []*os.File{socketFile},
	}
	fmt.Println("[crate] DEBUG starting container crate-init...")
	if err := container.Start(process); err != nil {
		return -1, err
	}

	// drop pid file
	fmt.Println("[crate] DEBUG dropping pid file...")
	pidfile := filepath.Join(containerDir, PID_FILE)
	pid, err := process.Pid()
	fmt.Println("[crate] DEBUG writing pid file...", pidfile, pid)
	if err != nil {
		return -1, err
	}
	if err := ioutil.WriteFile(pidfile, []byte(fmt.Sprintf("%d", pid)), 0600); err != nil {
		return -1, err
	}
	return pid, err
}

func (self *Crate) Load(id string) (*Container, error) {
	factory, err := self.factory(id)
	if err != nil {
		return nil, err
	}
	container, err := factory.Load(LIBCONTAINER_DIR)
	if err != nil {
		return nil, err
	}

	return &Container{
		id:        id,
		crate:     self,
		container: container,
	}, nil
}

func (self *Crate) Destroy(id string) error {
	// unmount rootfs
	fmt.Println("[crate] unmounting rootfs:", self.rootfs(id))
	if out, err := exec.Command("umount", self.rootfs(id)).CombinedOutput(); err != nil {
		fmt.Println("[crate] Error unmounting rootfs: ", string(out), err)
		return err
	}
	fmt.Println("[crate] loading container...")
	c, err := self.Load(id)
	if err != nil {
		return err
	}
	fmt.Println("[crate] destroying container...")
	err = c.Destroy()
	if err != nil {
		return err
	}
	return os.RemoveAll(self.path(id))
}

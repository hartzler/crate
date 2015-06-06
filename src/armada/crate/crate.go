package crate

import (
	"armada/pkg/tar"
	"encoding/json"
	"fmt"
	"github.com/docker/libcontainer"
	"github.com/docker/libcontainer/configs"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
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
	Root  string
	Cargo *Cargo
}

func New(root string) *Crate {
	return &Crate{
		Root:  root,
		Cargo: NewCargo(filepath.Join(root, "cargo")),
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

// TODO switch to our own config which includes options, cargo, and a libcontainer configs.Config
func (self *Crate) Create(id string, cargo []string, libconfig *configs.Config) (*Container, error) {

	// setup our container/rootfs dir
	containerDir := self.path(id)
	rootfs := filepath.Join(containerDir, "rootfs")
	libconfig.Rootfs = rootfs
	if err := self.setupRootfs(rootfs, cargo); err != nil {
		return nil, err
	}

	// libcontainer create
	container, err := setupLibcontainer(id, containerDir, libconfig)
	if err != nil {
		return nil, err
	}

	// start crate-init
	fmt.Println("PARENT: starting init...")
	if err := self.startInit(id, containerDir, container); err != nil {
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

func (self *Crate) setupRootfs(rootfs string, cargo []string) error {
	if err := os.MkdirAll(rootfs, 0700); err != nil {
		return err
	}

	// extract cargo!
	for _, path := range cargo {
		// make sure we got it!
		fmt.Println("Fetching: ", path)
		hash, err := self.Cargo.fetch(path)
		if err != nil {
			return err
		}
		reader, err := self.Cargo.load(hash)
		if err != nil {
			return err
		}
		defer reader.Close()

		fmt.Println("Extracting: ", hash)
		if err = tar.Extract(reader, rootfs); err != nil {
			return err
		}
	}
	return nil
}
func (self *Crate) startInit(id, containerDir string, container libcontainer.Container) error {
	// copy self to /crate-init in the container rootfs (non-portable hack?)
	fmt.Println("PARENT: copying crate-init...")
	exePath, err := os.Readlink("/proc/self/exe")
	if err != nil {
		return err
	}
	if _, err = copyFile(exePath, filepath.Join(containerDir, "rootfs", CRATE_INIT)); err != nil {
		return err
	}

	// stdout/err
	stdout, err := os.Create(filepath.Join(containerDir, "stdout"))
	if err != nil {
		return err
	}
	stderr, err := os.Create(filepath.Join(containerDir, "stderr"))
	if err != nil {
		return err
	}

	// set permission on listen to 0600...
	oldmask := syscall.Umask(0177)

	// ctrl socket file to pass to our crate-init
	socketPath := self.controlSocket(id)
	if err := os.Remove(socketPath); !os.IsNotExist(err) {
		return err
	}
	socket, err := net.ListenUnix("unix", &net.UnixAddr{socketPath, "unix"})
	if err != nil {
		return err
	}
	socketFile, err := socket.File()
	if err != nil {
		return err
	}

	// restore
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
	if err := container.Start(process); err != nil {
		return err
	}

	// drop pid file
	pidfile := filepath.Join(containerDir, PID_FILE)
	pid, err := process.Pid()
	fmt.Println("PARENT: writing pid file...", pidfile, pid)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(pidfile, []byte(fmt.Sprintf("%d", pid)), 0600)
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
	c, err := self.Load(id)
	if err != nil {
		return err
	}
	err = c.Destroy()
	if err != nil {
		return err
	}
	return os.RemoveAll(self.path(id))
}

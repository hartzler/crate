package crate

import (
	"armada/pkg/tar"
	"encoding/json"
	"fmt"
	"github.com/docker/libcontainer"
	"github.com/docker/libcontainer/configs"
	"io/ioutil"
	"os"
	"path/filepath"
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
	return &Crate{root}
}

func (self *Crate) SetupRoot() error {
	if _, err := os.Stat(self.Root); os.IsNotExist(err) {
		return os.MkdirAll(self.Root, 0700)
	}
	return nil
}

func (self *Crate) factory(id string) (libcontainer.Factory, error) {
	containerDir := filepath.Join(self.Root, id)
	return libcontainer.New(containerDir, libcontainer.Cgroupfs)
}

func (self *Crate) path(id string) string {
	return filepath.Join(self.Root, id)
}

// TODO switch to our own config which includes options, a Dotcrate, and a libcontainer configs.Config
func (self *Crate) Create(id string, dotcrate *Dotcrate, libconfig *configs.Config) (*Container, error) {

	// setup our container/rootfs dir
	containerDir := self.path(id)
	rootfs := filepath.Join(containerDir, "rootfs")
	libconfig.Rootfs = rootfs
	if err := setupRootfs(rootfs, dotcrate); err != nil {
		return nil, err
	}

	// libcontainer create
	container, err := setupLibcontainer(id, containerDir, libconfig)
	if err != nil {
		return nil, err
	}

	// start crate-init
	fmt.Println("PARENT: starting init...")
	if err := startInit(id, containerDir, container); err != nil {
		return nil, err
	}

	return &Container{
		id:        id,
		crate:     self,
		container: container,
	}, err

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

	factory, err := libcontainer.New(containerDir, libcontainer.Cgroupfs)
	if err != nil {
		return nil, err
	}
	container, err := factory.Create(LIBCONTAINER_DIR, libconfig)
	if err != nil {
		return nil, err
	}

	return container, nil
}

func setupRootfs(rootfs string, dotcrate *Dotcrate) error {
	if err := os.MkdirAll(rootfs, 0700); err != nil {
		return err
	}

	// extract cargo!
	for _, path := range dotcrate.Cargo {
		//pieces := strings.Split(path, "/")
		//tarfile := pieces[len(pieces)-1]
		//err = os.Stat(tarfile);
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		if err = tar.Extract(file, rootfs); err != nil {
			return err
		}
	}
	return nil
}
func startInit(id, containerDir string, container libcontainer.Container) error {
	// copy self to /crate-init in the container rootfs (non-portable hack?)
	fmt.Println("PARENT: copying crate-init...")
	exePath, err := os.Readlink("/proc/self/exe")
	if err != nil {
		return err
	}
	if _, err = copyFile(exePath, filepath.Join(containerDir, "rootfs", CRATE_INIT)); err != nil {
		return err
	}

	process := &libcontainer.Process{
		Args:   []string{"/" + CRATE_INIT, id},
		Env:    StandardEnvironment,
		User:   "",
		Cwd:    "/",
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
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

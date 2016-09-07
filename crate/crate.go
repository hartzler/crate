package crate

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Crate struct {
	Root string
}

func New(root string) *Crate {
	return &Crate{
		Root: root,
	}
}

func (self *Crate) cargoPath() string {
	return filepath.Join(self.Root, ".cargo")
}

func (self *Crate) manifestPath() string {
	return filepath.Join(self.Root, ".manifest")
}

func (self *Crate) cratePath(id string) string {
	return filepath.Join(self.Root, id)
}

func (self *Crate) SetupRoot() error {
	if _, err := os.Stat(self.cargoPath()); os.IsNotExist(err) {
		return os.MkdirAll(self.cargoPath(), 0755)
	}
	if _, err := os.Stat(self.manifestPath()); os.IsNotExist(err) {
		return os.MkdirAll(self.manifestPath(), 0755)
	}
	return nil
}

func (self *Crate) Install(id string, cargo []string) error {
	crateDir := self.cratePath(id)

	// get list of cargo mounts
	lowerDirs := []string{}
	for _, path := range cargo {
		// strip off cargo
		lowerDirs = append(lowerDirs, filepath.Join(self.cargoPath(), strings.Replace(path, ".cargo", "", -1)))
	}

	if len(cargo) > 1 {
		// setup crate dir
		fmt.Println("[crate] creating crate dir: ", crateDir)
		if err := os.MkdirAll(crateDir, 0700); err != nil {
			return err
		}

		// mount overlay
		opts := fmt.Sprintf("lowerdir=%s", strings.Join(lowerDirs, ":"))
		args := []string{"-t", "overlay", "-o", opts, "overlay", crateDir}
		fmt.Println("[crate] mounting merged rootfs: mount", strings.Join(args, " "))
		if out, err := exec.Command("mount", args...).CombinedOutput(); err != nil {
			fmt.Println(string(out))
			return err
		}
	} else {
		// just symlink
		if out, err := exec.Command("ln", "-s", lowerDirs[0], crateDir).CombinedOutput(); err != nil {
			fmt.Println(string(out))
			return err
		}
	}
	return nil
}

func (self *Crate) Remove(id string) error {
	// unmount crate dir
	fmt.Println("[crate] unmounting crate dir:", self.cratePath(id))
	if out, err := exec.Command("umount", self.cratePath(id)).CombinedOutput(); err != nil {
		fmt.Println("[crate] Error unmounting crate dir: ", string(out), err)
		return err
	}
	return os.RemoveAll(self.cratePath(id))
}

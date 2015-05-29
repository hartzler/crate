package crate

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Dotcrate struct {
	Cargo []string
}

func LoadDot(path string) (*Dotcrate, error) {
	var dotcrate Dotcrate
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(bytes, &dotcrate); err != nil {
		return nil, err
	}
	return &dotcrate, nil
}

func (self *Dotcrate) Store(path string) error {
	bytes, err := yaml.Marshal(self)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, bytes, 0600)
}

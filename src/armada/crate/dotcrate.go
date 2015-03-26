package crate

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Connector struct {
	Name        string
	Proto       string
	Port        string
	Cardinatity string
}

type Connectors struct {
	In  []*Connector
	Out []*Connector
}

type Dotcrate struct {
	Name       string
	Connectors Connectors
	Hooks      map[string]string
	Cargo      []string
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

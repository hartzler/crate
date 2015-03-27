package crate

import (
	"github.com/peterbourgon/diskv"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
)

type Cargo struct {
	kv *diskv.Diskv
}

func NewCargo(base string) *Cargo {
	return &Cargo{
		kv: diskv.New(diskv.Options{
			BasePath:     base,
			Transform:    transform,
			CacheSizeMax: 0, // no cache...
		}),
	}
}

func transform(key string) []string {
	return []string{
		key[0:2],
	}
}

func (self *Cargo) fetch(cargoUrl string) (string, error) {
	// check if we have it first
	u, err := url.Parse(cargoUrl)
	if err != nil {
		return "", err
	}
	filename := filepath.Base(u.Path)
	if self.kv.Has(filename) {
		return filename, nil
	}

	// lets download it
	res, err := http.Get(cargoUrl)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	// write and fsync
	return filename, self.kv.WriteStream(filename, res.Body, true)
}

func (self *Cargo) load(hash string) (io.ReadCloser, error) {
	return self.kv.ReadStream(hash, true)
}

func (self *Cargo) remove(hash string) error {
	return self.kv.Erase(hash)
}

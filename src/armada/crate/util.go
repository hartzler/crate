package crate

import (
	"fmt"
	"io"
	"os"
)

func copyFile(src, dst string) (int64, error) {
	sf, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer sf.Close()
	fmt.Println("OPEN:", dst, "0700")
	df, err := os.OpenFile(dst, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0700)
	if err != nil {
		return 0, err
	}
	defer df.Close()
	fmt.Println("COPY:", src, dst)
	return io.Copy(df, sf)
}

// exists returns whether the given file or directory exists or not
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

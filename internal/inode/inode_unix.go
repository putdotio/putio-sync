// +build linux darwin

package inode

import (
	"os"
	"syscall"
)

func Get(fi os.FileInfo) (uint64, error) {
	stat := fi.Sys().(*syscall.Stat_t)
	return stat.Ino, nil
}

func GetPath(path string) (uint64, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return Get(fi)
}

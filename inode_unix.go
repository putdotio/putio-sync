// +build linux darwin

package main

import (
	"os"
	"syscall"
)

func GetInode(fi os.FileInfo) (uint64, error) {
	stat := fi.Sys().(*syscall.Stat_t)
	return stat.Ino, nil
}

func GetInodePath(path string) (uint64, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return GetInode(fi)
}

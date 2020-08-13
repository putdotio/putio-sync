package main

import (
	"os"
	"syscall"
)

func GetInode(fi os.FileInfo) (uint64, error) {
	return GetInodePath(fi.Name())
}

func GetInodePath(path string) (uint64, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	var d syscall.ByHandleFileInformation
	err = syscall.GetFileInformationByHandle(syscall.Handle(f.Fd()), &d)
	if err != nil {
		return 0, err
	}
	return (uint64(d.FileIndexHigh) << 32) | uint64(d.FileIndexLow), nil
}

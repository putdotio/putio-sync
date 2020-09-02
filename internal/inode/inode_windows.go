package inode

import (
	"os"
	"syscall"
)

func Get(path string, _ os.FileInfo) (uint64, error) {
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

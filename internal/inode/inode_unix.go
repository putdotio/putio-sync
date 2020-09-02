// +build linux darwin

package inode

import (
	"os"
	"syscall"
)

func Get(path string, fi os.FileInfo) (uint64, error) {
	if fi == nil {
		var err error
		fi, err = os.Stat(path)
		if err != nil {
			return 0, err
		}
	}
	return fi.Sys().(*syscall.Stat_t).Ino, nil
}

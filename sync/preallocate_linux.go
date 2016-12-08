package sync

import (
	"os"
	"syscall"
)

func fallocate(f *os.File, size int64) error {
	err := syscall.Fallocate(int(f.Fd()), 0, 0, size)
	if err == syscall.ENOTSUP {
		return f.Truncate(size)
	}
	return err
}

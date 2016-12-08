package sync

import (
	"os"
	"syscall"
)

func preallocate(f *os.File, size int64) error {
	return syscall.Fallocate(int(f.Fd()), 0, 0, size)
}

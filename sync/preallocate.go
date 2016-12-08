// +build !linux

package sync

import "os"

func fallocate(f *os.File, size int64) error {
	return f.Truncate(size)
}

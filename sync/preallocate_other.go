// +build !linux,!darwin

package sync

import (
	"fmt"
	"os"
)

func preallocate(f *os.File, size int64) error {
	return fmt.Errorf("Operation not supported on this platform")
}

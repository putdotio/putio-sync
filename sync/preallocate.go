package sync

import "os"

// Preallocate tries to allocate continuous disk space for the given file. If
// system specific pre-allocation fails, it fills the file with zeroes.
func Preallocate(f *os.File, size int64) error {
	err := preallocate(f, size)
	if err == nil {
		return nil
	}

	// use default truncation if platform specific allocation fails
	return f.Truncate(size)
}

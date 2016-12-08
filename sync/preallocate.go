package sync

import "os"

func Preallocate(f *os.File, size int64) error {
	err := preallocate(f, size)
	if err == nil {
		return nil
	}

	// use default truncation if platform specific allocation fails
	return f.Truncate(size)
}

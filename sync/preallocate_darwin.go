package sync

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

// preallocate tries to allocate continuous disk space for the given file. On a
// succesive operation, ftruncate() call is needed on macOS to have the file
// size reported correctly. Without the call, space is being reserved but file
// size reported is zero.
func preallocate(f *os.File, size int64) error {
	fstore := &syscall.Fstore_t{
		Flags:   syscall.F_ALLOCATECONTIG,
		Posmode: syscall.F_PEOFPOSMODE,
		Offset:  int64(0),
		Length:  size,
	}

	// Try to get a continuous chunk of disk space
	_, _, errno := syscall.Syscall(syscall.SYS_FCNTL, f.Fd(), uintptr(syscall.F_PREALLOCATE), uintptr(unsafe.Pointer(fstore)))
	if errno == syscall.ENOTSUP {
		return fmt.Errorf("Operation not supported")
	}

	if errno == 0 {
		return f.Truncate(size)
	}

	// OK, perhaps we are too fragmented, allocate non-continuous space
	fstore.Flags = syscall.F_ALLOCATEALL
	_, _, errno = syscall.Syscall(syscall.SYS_FCNTL, f.Fd(), uintptr(syscall.F_PREALLOCATE), uintptr(unsafe.Pointer(fstore)))
	if errno == syscall.ENOTSUP {
		return fmt.Errorf("Operation not supported")
	}

	if errno != 0 {
		return fmt.Errorf("error: %v", errno)
	}

	return f.Truncate(size)
}

package tmpdir

import (
	"os"
	"path/filepath"
	"syscall"
)

const Name = "putio-sync-tmp"

func Create(dir string) (string, error) {
	filename := filepath.Join(dir, Name)
	err := os.MkdirAll(filename, 0777)
	if err != nil {
		return "", err
	}
	filenameW, err := syscall.UTF16PtrFromString(filename)
	if err != nil {
		return "", err
	}
	err = syscall.SetFileAttributes(filenameW, syscall.FILE_ATTRIBUTE_HIDDEN)
	if err != nil {
		return "", err
	}
	return filename, nil
}

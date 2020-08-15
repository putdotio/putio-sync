// +build linux darwin

package main

import (
	"os"
	"path/filepath"
)

const tempDirName = ".putio-sync-tmp"

func CreateTempDir() (string, error) {
	filename := filepath.Join(localPath, tempDirName)
	err := os.MkdirAll(filename, 0777)
	if err != nil {
		return "", err
	}
	return filename, nil
}

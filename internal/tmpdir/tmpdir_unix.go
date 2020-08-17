// +build linux darwin

package tmpdir

import (
	"os"
	"path/filepath"
)

const Name = ".putio-sync-tmp"

func Create(dir string) (string, error) {
	filename := filepath.Join(dir, Name)
	err := os.MkdirAll(filename, 0777)
	if err != nil {
		return "", err
	}
	return filename, nil
}

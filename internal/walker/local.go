package walker

import (
	"os"
	"path/filepath"
)

type localWalker struct {
	root string
}

func (w *localWalker) Walk(walkFn walkFunc) error {
	return filepath.Walk(w.root, func(path string, info os.FileInfo, err error) error {
		relpath, err2 := filepath.Rel(w.root, path)
		if err2 != nil {
			panic(err2)
		}
		return walkFn(newLocalFile(info, relpath, w.root), err)
	})
}

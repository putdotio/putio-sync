package putiosync

import "github.com/syncthing/syncthing/lib/fs"

func shouldIgnoreName(name string) bool {
	return fs.WindowsInvalidFilename(name) != nil
}

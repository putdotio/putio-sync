package walker

import (
	"os"
	"time"

	"github.com/putdotio/go-putio"
)

type fileInfo struct {
	putio.File
}

var _ os.FileInfo = (*fileInfo)(nil)

func newFileInfo(f putio.File) *fileInfo { return &fileInfo{File: f} }
func (fi *fileInfo) Name() string        { return fi.File.Name }
func (fi *fileInfo) Size() int64         { return fi.File.Size }
func (fi *fileInfo) ModTime() time.Time  { return fi.File.CreatedAt.Time }
func (fi *fileInfo) Sys() interface{}    { return nil }
func (fi *fileInfo) IsDir() bool         { return fi.File.IsDir() }
func (fi *fileInfo) Mode() os.FileMode {
	if fi.IsDir() {
		return 755 | os.ModeDir
	}
	return 644
}

package main

import (
	"os"
	"time"

	"github.com/putdotio/go-putio"
)

type FileInfo struct {
	putio.File
}

var _ os.FileInfo = (*FileInfo)(nil)

func newFileInfo(f putio.File) *FileInfo { return &FileInfo{File: f} }
func (fi *FileInfo) Name() string        { return fi.File.Name }
func (fi *FileInfo) Size() int64         { return fi.File.Size }
func (fi *FileInfo) ModTime() time.Time  { return fi.File.CreatedAt.Time }
func (fi *FileInfo) Sys() interface{}    { return nil }
func (fi *FileInfo) IsDir() bool         { return fi.File.IsDir() }
func (fi *FileInfo) Mode() os.FileMode {
	if fi.IsDir() {
		return 755 | os.ModeDir
	}
	return 644
}

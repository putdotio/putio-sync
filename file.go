package main

import (
	"os"
)

type File struct {
	os.FileInfo
	path string
}

func newFile(path string, fi os.FileInfo) File {
	return File{
		FileInfo: fi,
		path:     path,
	}
}

func (f File) Path() string { return f.path }

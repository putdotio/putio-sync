package main

import (
	"os"
	"path/filepath"

	"github.com/putdotio/go-putio"
)

type File interface {
	Info() os.FileInfo
	// RelPath is for comparing with local/remote path. Must always return a path with slash as path seperator.
	RelPath() string
}

type LocalFile struct {
	info    os.FileInfo
	root    string
	relpath string
}

func NewLocalFile(info os.FileInfo, root, relpath string) *LocalFile {
	return &LocalFile{
		info:    info,
		root:    root,
		relpath: filepath.ToSlash(relpath),
	}
}

func (f *LocalFile) Info() os.FileInfo { return f.info }
func (f *LocalFile) RelPath() string   { return f.relpath }
func (f *LocalFile) AbsPath() string   { return filepath.Join(f.root, f.relpath) }

type RemoteFile struct {
	putioFile putio.File
	relpath   string
}

func NewRemoteFile(pf putio.File, relpath string) *RemoteFile {
	return &RemoteFile{
		putioFile: pf,
		relpath:   relpath,
	}
}

func (f *RemoteFile) Info() os.FileInfo { return newFileInfo(f.putioFile) }
func (f *RemoteFile) RelPath() string   { return f.relpath }

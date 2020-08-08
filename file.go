package main

import (
	"os"
	"path/filepath"

	"github.com/putdotio/go-putio"
)

type File interface {
	Info() os.FileInfo
	// RelPath is for comparing with local/remote path. Path seperator must be slash.
	RelPath() string
}

type LocalFile struct {
	info    os.FileInfo
	relpath string
}

func NewLocalFile(info os.FileInfo, relpath string) *LocalFile {
	return &LocalFile{
		info:    info,
		relpath: filepath.ToSlash(relpath),
	}
}

func (f *LocalFile) Info() os.FileInfo { return f.info }
func (f *LocalFile) RelPath() string   { return f.relpath }

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

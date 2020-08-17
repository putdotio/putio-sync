package walker

import (
	"os"
	"path/filepath"

	"github.com/putdotio/go-putio"
)

type file interface {
	Info() os.FileInfo
	// RelPath is for comparing with local/remote path. Path separator must be slash.
	RelPath() string
}

type LocalFile struct {
	info    os.FileInfo
	relpath string
	root    string
}

func newLocalFile(info os.FileInfo, relpath, root string) *LocalFile {
	return &LocalFile{
		info:    info,
		relpath: filepath.ToSlash(relpath),
		root:    root,
	}
}

func (f *LocalFile) Info() os.FileInfo { return f.info }
func (f *LocalFile) RelPath() string   { return f.relpath }
func (f *LocalFile) FullPath() string  { return filepath.Join(f.root, filepath.FromSlash(f.relpath)) }

type RemoteFile struct {
	putioFile putio.File
	relpath   string
}

func newRemoteFile(pf putio.File, relpath string) *RemoteFile {
	return &RemoteFile{
		putioFile: pf,
		relpath:   relpath,
	}
}

func (f *RemoteFile) Info() os.FileInfo      { return newFileInfo(f.putioFile) }
func (f *RemoteFile) RelPath() string        { return f.relpath }
func (f *RemoteFile) PutioFile() *putio.File { return &f.putioFile }

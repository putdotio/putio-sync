package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/putdotio/go-putio"
)

type File interface {
	Info() os.FileInfo
	// RelPath is for comparing with local/remote path. Must always return a path with slash as path seperator.
	RelPath() string
	OpenForRead(offset int64) (io.ReadCloser, error)
	OpenForWrite(offset int64) (io.WriteCloser, error)
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

func CreateLocalFile(relpath string) (*LocalFile, error) {
	of, err := os.OpenFile(filepath.Join(localPath, filepath.FromSlash(relpath)), os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer of.Close()
	fi, err := of.Stat()
	if err != nil {
		return nil, err
	}
	err = of.Close()
	if err != nil {
		return nil, err
	}
	return NewLocalFile(fi, relpath), nil
}

func (f *LocalFile) Info() os.FileInfo { return f.info }
func (f *LocalFile) RelPath() string   { return f.relpath }

func (f *LocalFile) OpenForRead(offset int64) (rc io.ReadCloser, err error) {
	panic("LocalFile.OpenForRead not implemented")
}

func (f *LocalFile) OpenForWrite(offset int64) (wc io.WriteCloser, err error) {
	of, err := os.OpenFile(filepath.Join(localPath, filepath.FromSlash(f.relpath)), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return
	}
	wc = of
	_, err = of.Seek(offset, io.SeekStart)
	return
}

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

func (f *RemoteFile) OpenForRead(offset int64) (rc io.ReadCloser, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	u, err := client.Files.URL(ctx, f.putioFile.ID, true)
	if err != nil {
		return
	}
	// TODO do not use default http client
	// TODO use proper timeouts on client
	// TODO retry failed operations
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return
	}
	req.Header.Set("range", fmt.Sprintf("bytes=%d-", offset))
	resp, err := http.Get(u)
	if err != nil {
		return
	}
	rc = resp.Body
	return
}

func (f *RemoteFile) OpenForWrite(offset int64) (wc io.WriteCloser, err error) {
	panic("RemoteFile.OpenForWrite not implemented")
}

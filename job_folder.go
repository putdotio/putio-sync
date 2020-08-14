package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

type CreateLocalFolder struct {
	relpath  string
	remoteID int64
}

func (j *CreateLocalFolder) String() string {
	return "Creating local folder " + j.relpath
}

func (j *CreateLocalFolder) Run(ctx context.Context) error {
	dirPath := filepath.Join(localPath, filepath.FromSlash(j.relpath))
	err := os.MkdirAll(dirPath, 0777)
	if err != nil {
		return err
	}
	s := State{
		Status:   StatusSynced,
		IsDir:    true,
		RemoteID: j.remoteID,
		relpath:  j.relpath,
	}
	return s.Write()
}

type CreateRemoteFolder struct {
	relpath string
}

func (j *CreateRemoteFolder) String() string {
	return fmt.Sprintf("Creating remote folder %q", j.relpath)
}

func (j *CreateRemoteFolder) Run(ctx context.Context) error {
	remoteID, err := dirCache.Mkdirp(j.relpath)
	if err != nil {
		return err
	}
	s := State{
		Status:   StatusSynced,
		IsDir:    true,
		RemoteID: remoteID,
		relpath:  j.relpath,
	}
	return s.Write()
}

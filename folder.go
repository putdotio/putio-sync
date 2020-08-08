package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type CreateLocalFolder struct {
	relpath string
	state   *State
}

func (j *CreateLocalFolder) String() string {
	return "Creating local folder " + j.relpath
}

func (j *CreateLocalFolder) Run() error {
	return os.MkdirAll(filepath.Join(localPath, filepath.FromSlash(j.relpath)), 0777)
}

type CreateRemoteFolder struct {
	relpath string
	state   *State
}

func (j *CreateRemoteFolder) String() string {
	return fmt.Sprintf("Creating remote folder %q", j.relpath)
}

func (j *CreateRemoteFolder) Run() error {
	_, err := dirCache.Mkdirp(j.relpath)
	return err
}

package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
)

type MoveLocalFile struct {
	localFile LocalFile
	toRelpath string
	state     State
}

func (j *MoveLocalFile) String() string {
	return fmt.Sprintf("Moving local file from %q to %q", j.state.relpath, j.toRelpath)
}

func (j *MoveLocalFile) Run(ctx context.Context) error {
	oldPath := filepath.Join(localPath, filepath.FromSlash(j.localFile.relpath))
	newPath := filepath.Join(localPath, filepath.FromSlash(j.toRelpath))
	exists, err := j.exists(newPath)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("file already exists at move target")
	}
	err = os.Rename(oldPath, newPath)
	if err != nil {
		return err
	}
	return j.state.Move(j.toRelpath)
}

func (j *MoveLocalFile) exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

type MoveRemoteFile struct {
	remoteFile RemoteFile
	toRelpath  string
	state      State
}

func (j *MoveRemoteFile) String() string {
	return fmt.Sprintf("Moving remote file from %q to %q", j.state.relpath, j.toRelpath)
}

func (j *MoveRemoteFile) Run(ctx context.Context) error {
	dir := path.Dir(j.toRelpath)
	parentID, err := dirCache.Mkdirp(dir)
	if err != nil {
		return err
	}
	// TODO check if there is any file at target path
	// exists, err := j.exists(newPath)
	// if err != nil {
	// 	return err
	// }
	// if exists {
	// 	return errors.New("file already exists at move target")
	// }
	err = client.Files.Move(ctx, parentID, j.remoteFile.putioFile.ID)
	if err != nil {
		return err
	}
	return j.state.Move(j.toRelpath)
}

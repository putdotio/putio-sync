package main

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
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
	err = os.MkdirAll(filepath.Dir(newPath), 0777)
	if err != nil {
		return err
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
	dir, name := path.Split(j.toRelpath)
	parentID, err := dirCache.Mkdirp(ctx, dir)
	if err != nil {
		return err
	}
	err = moveRemoteFile(ctx, parentID, j.remoteFile.putioFile.ID, name)
	if err != nil {
		return err
	}
	return j.state.Move(j.toRelpath)
}

func moveRemoteFile(ctx context.Context, parentID, fileID int64, name string) error {
	params := url.Values{}
	params.Set("file_id", strconv.FormatInt(fileID, 10))
	params.Set("parent_id", strconv.FormatInt(parentID, 10))
	params.Set("name", name)

	req, err := client.NewRequest(ctx, "POST", "/v2/files/move", strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	_, err = client.Do(req, nil)
	return err
}

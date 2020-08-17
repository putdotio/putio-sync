package putiosync

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

type moveLocalFileJob struct {
	localFile iLocalFile
	toRelpath string
	state     stateType
}

func (j *moveLocalFileJob) String() string {
	return fmt.Sprintf("Moving local file from %q to %q", j.state.relpath, j.toRelpath)
}

func (j *moveLocalFileJob) Run(ctx context.Context) error {
	oldPath := j.localFile.FullPath()
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

func (j *moveLocalFileJob) exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

type moveRemoteFileJob struct {
	remoteFile iRemoteFile
	toRelpath  string
	state      stateType
}

func (j *moveRemoteFileJob) String() string {
	return fmt.Sprintf("Moving remote file from %q to %q", j.state.relpath, j.toRelpath)
}

func (j *moveRemoteFileJob) Run(ctx context.Context) error {
	dir, name := path.Split(j.toRelpath)
	parentID, err := dirCache.Mkdirp(ctx, dir)
	if err != nil {
		return err
	}
	err = moveRemoteFile(ctx, parentID, j.remoteFile.PutioFile().ID, name)
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

package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

type DeleteLocalFile struct {
	localFile LocalFile
	state     State
}

func (j *DeleteLocalFile) String() string {
	return fmt.Sprintf("Deleting local file %q", j.state.relpath)
}

func (j *DeleteLocalFile) Run(ctx context.Context) error {
	var err error
	removePath := filepath.Join(localPath, filepath.FromSlash(j.localFile.relpath))
	if j.localFile.info.IsDir() {
		// TODO Check if dir is empty before delete
		// var files []File
		// files, err = WalkOnFolder(ctx, LocalWalker{root: removePath})
		// if err != nil {
		// 	return err
		// }
		// if len(files) > 0 {
		// 	return errors.New("folder not empty: " + removePath)
		// }
		err = os.RemoveAll(removePath)
	} else {
		err = os.Remove(removePath)
	}
	if os.IsNotExist(err) {
		err = nil
	}
	if err != nil {
		return err
	}
	return j.state.Delete()
}

type DeleteRemoteFile struct {
	remoteFile RemoteFile
	state      State
}

func (j *DeleteRemoteFile) String() string {
	return fmt.Sprintf("Deleting remote file %q", j.state.relpath)
}

func (j *DeleteRemoteFile) Run(ctx context.Context) error {
	err := client.Files.Delete(ctx, j.remoteFile.putioFile.ID)
	if err != nil {
		return err
	}
	return j.state.Delete()
}

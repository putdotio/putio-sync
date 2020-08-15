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
	err := os.RemoveAll(filepath.Join(localPath, filepath.FromSlash(j.localFile.relpath)))
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

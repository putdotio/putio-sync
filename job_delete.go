package putiosync

import (
	"context"
	"fmt"
	"os"
)

type deleteLocalFileJob struct {
	localFile iLocalFile
	state     stateType
}

func (j *deleteLocalFileJob) String() string {
	return fmt.Sprintf("Deleting local file %q", j.state.relpath)
}

func (j *deleteLocalFileJob) Run(ctx context.Context) error {
	err := os.RemoveAll(j.localFile.FullPath())
	if err != nil {
		return err
	}
	return j.state.Delete()
}

type deleteRemoteFileJob struct {
	remoteFile iRemoteFile
	state      stateType
}

func (j *deleteRemoteFileJob) String() string {
	return fmt.Sprintf("Deleting remote file %q", j.state.relpath)
}

func (j *deleteRemoteFileJob) Run(ctx context.Context) error {
	err := client.Files.Delete(ctx, j.remoteFile.PutioFile().ID)
	if err != nil {
		return err
	}
	return j.state.Delete()
}

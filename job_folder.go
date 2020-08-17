package putiosync

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

type createLocalFolderJob struct {
	relpath  string
	remoteID int64
}

func (j *createLocalFolderJob) String() string {
	return "Creating local folder " + j.relpath
}

func (j *createLocalFolderJob) Run(ctx context.Context) error {
	err := os.MkdirAll(filepath.Join(localPath, filepath.FromSlash(j.relpath)), 0777)
	if err != nil {
		return err
	}
	s := stateType{
		Status:   statusSynced,
		IsDir:    true,
		RemoteID: j.remoteID,
		relpath:  j.relpath,
	}
	return s.Write()
}

type createRemoteFolderJob struct {
	relpath string
}

func (j *createRemoteFolderJob) String() string {
	return fmt.Sprintf("Creating remote folder %q", j.relpath)
}

func (j *createRemoteFolderJob) Run(ctx context.Context) error {
	remoteID, err := dirCache.Mkdirp(ctx, j.relpath)
	if err != nil {
		return err
	}
	s := stateType{
		Status:   statusSynced,
		IsDir:    true,
		RemoteID: remoteID,
		relpath:  j.relpath,
	}
	return s.Write()
}

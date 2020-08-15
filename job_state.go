package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cenkalti/log"
)

type DeleteState struct {
	state State
}

func (j *DeleteState) String() string {
	return fmt.Sprintf("Deleting state %q", j.state.relpath)
}

func (j *DeleteState) Run(ctx context.Context) error {
	if j.state.DownloadTempName != "" {
		err := os.Remove(filepath.Join(localPath, tempDirName, j.state.DownloadTempName))
		if err != nil {
			log.Errorln("cannot remove temp download file:", err.Error())
		}
	}
	if j.state.UploadURL != "" {
		err := TerminateUpload(ctx, token, j.state.UploadURL)
		if err != nil {
			log.Errorln("cannot remove upload:", err.Error())
		}
	}
	return j.state.Delete()
}

type WriteFileState struct {
	localFile  LocalFile
	remoteFile RemoteFile
}

func (j *WriteFileState) String() string {
	return fmt.Sprintf("Saving file state %q", j.localFile.relpath)
}

func (j *WriteFileState) Run(ctx context.Context) error {
	inode, err := GetInode(j.localFile.info)
	if err != nil {
		return err
	}
	s := State{
		Status:     StatusSynced,
		LocalInode: inode,
		RemoteID:   j.remoteFile.putioFile.ID,
		Size:       j.remoteFile.putioFile.Size,
		CRC32:      j.remoteFile.putioFile.CRC32,
		relpath:    j.remoteFile.relpath,
	}
	return s.Write()
}

type WriteDirState struct {
	remoteID int64
	relpath  string
}

func (j *WriteDirState) String() string {
	return fmt.Sprintf("Saving folder state %q", j.relpath)
}

func (j *WriteDirState) Run(ctx context.Context) error {
	s := State{
		Status:   StatusSynced,
		IsDir:    true,
		RemoteID: j.remoteID,
		relpath:  j.relpath,
	}
	return s.Write()
}

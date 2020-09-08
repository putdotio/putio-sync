package putiosync

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cenkalti/log"
	"github.com/putdotio/putio-sync/v2/internal/inode"
)

type deleteStateJob struct {
	state stateType
}

func (j *deleteStateJob) String() string {
	return fmt.Sprintf("Deleting state %q", j.state.relpath)
}

func (j *deleteStateJob) Run(ctx context.Context) error {
	if j.state.DownloadTempName != "" {
		err := os.Remove(filepath.Join(tempDirPath, j.state.DownloadTempName))
		if err != nil {
			log.Errorln("cannot remove temp download file:", err.Error())
		}
	}
	if j.state.UploadURL != "" {
		err := client.Upload.TerminateUpload(ctx, j.state.UploadURL)
		if err != nil {
			log.Errorln("cannot remove upload:", err.Error())
		}
	}
	return j.state.Delete()
}

type writeFileStateJob struct {
	localFile  iLocalFile
	remoteFile iRemoteFile
}

func (j *writeFileStateJob) String() string {
	return fmt.Sprintf("Saving file state %q", j.localFile.RelPath())
}

func (j *writeFileStateJob) Run(ctx context.Context) error {
	in, err := inode.Get(j.localFile.FullPath(), j.localFile.Info())
	if err != nil {
		return err
	}
	s := stateType{
		Status:     statusSynced,
		LocalInode: in,
		RemoteID:   j.remoteFile.PutioFile().ID,
		Size:       j.remoteFile.PutioFile().Size,
		CRC32:      j.remoteFile.PutioFile().CRC32,
		relpath:    j.remoteFile.RelPath(),
	}
	return s.Write()
}

type writeDirStateJob struct {
	remoteID int64
	relpath  string
}

func (j *writeDirStateJob) String() string {
	return fmt.Sprintf("Saving folder state %q", j.relpath)
}

func (j *writeDirStateJob) Run(ctx context.Context) error {
	s := stateType{
		Status:   statusSynced,
		IsDir:    true,
		RemoteID: j.remoteID,
		relpath:  j.relpath,
	}
	return s.Write()
}

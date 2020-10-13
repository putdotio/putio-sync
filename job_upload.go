package putiosync

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/cenkalti/log"
	"github.com/putdotio/putio-sync/v2/internal/inode"
	"github.com/putdotio/putio-sync/v2/internal/progress"
	"github.com/putdotio/putio-sync/v2/internal/watcher"
)

type uploadJob struct {
	localFile iLocalFile
	state     *stateType
}

func (d *uploadJob) String() string {
	return fmt.Sprintf("Uploading %q", d.localFile.RelPath())
}

func (d *uploadJob) tryResume(ctx context.Context) bool {
	if d.state == nil {
		return false
	}
	if d.state.Status != statusUploading {
		return false
	}
	if d.state.UploadURL == "" {
		return false
	}
	if d.state.Size != d.localFile.Info().Size() {
		return false
	}
	in, _ := inode.Get(d.localFile.FullPath(), d.localFile.Info())
	if d.state.LocalInode != in {
		return false
	}
	offset, err := client.Upload.GetOffset(ctx, d.state.UploadURL)
	if err != nil {
		return false
	}
	d.state.Offset = offset
	return offset <= d.localFile.Info().Size()
}

func (d *uploadJob) Run(ctx context.Context) error {
	modwatch, err := watcher.WatchFileModification(ctx, d.localFile.FullPath())
	if err != nil {
		return err
	}
	defer modwatch.Stop()

	ok := d.tryResume(ctx)
	if !ok {
		in, err := inode.Get(d.localFile.FullPath(), d.localFile.Info())
		if err != nil {
			return err
		}
		dir, filename := path.Split(d.localFile.RelPath())
		parentID, err := dirCache.Mkdirp(ctx, dir)
		if err != nil {
			return err
		}
		location, err := client.Upload.CreateUpload(ctx, filename, parentID, d.localFile.Info().Size(), true)
		if err != nil {
			return err
		}
		d.state = &stateType{
			Status:     statusUploading,
			LocalInode: in,
			UploadURL:  location,
			Size:       d.localFile.Info().Size(),
			relpath:    d.localFile.RelPath(),
		}
		err = d.state.Write()
		if err != nil {
			return err
		}
	}
	f, err := os.Open(d.localFile.FullPath())
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Seek(d.state.Offset, io.SeekStart)
	if err != nil {
		return err
	}
	pr := progress.New(f, d.state.Offset, d.state.Size, d.String())
	pr.Start()
	fileID, crc32, err := client.Upload.SendFile(modwatch.Context(), pr, d.state.UploadURL, d.state.Offset)
	pr.Stop()
	modified := modwatch.Stop()
	if modified {
		log.Warningln("File modified while uploading")
		return nil
	}
	if err != nil {
		return err
	}
	d.state.Status = statusSynced
	d.state.RemoteID = fileID
	d.state.CRC32 = crc32
	err = d.state.Write()
	if err != nil {
		return err
	}
	return nil
}

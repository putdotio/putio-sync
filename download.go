package putiosync

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/cenkalti/log"
	"github.com/putdotio/putio-sync/v2/internal/inode"
	"github.com/putdotio/putio-sync/v2/internal/progress"
)

type downloadJob struct {
	remoteFile iRemoteFile
	state      *stateType
}

func (d *downloadJob) String() string {
	return fmt.Sprintf("Downloading %q", d.remoteFile.RelPath())
}

func (d *downloadJob) tryResume() io.WriteCloser {
	if d.state == nil {
		return nil
	}
	if d.state.Status != statusDownloading {
		return nil
	}
	if d.state.DownloadTempName == "" {
		return nil
	}
	f, err := os.OpenFile(filepath.Join(tempDirPath, d.state.DownloadTempName), os.O_WRONLY, 0)
	if err != nil {
		return nil
	}
	defer func() {
		if err != nil && f != nil {
			f.Close()
		}
	}()
	if d.state.Size != d.remoteFile.PutioFile().Size {
		return nil
	}
	if d.state.CRC32 != d.remoteFile.PutioFile().CRC32 {
		return nil
	}
	n, err := f.Seek(0, io.SeekEnd)
	if err != nil {
		return nil
	}
	if n < d.state.Offset {
		return nil
	}
	_, err = f.Seek(d.state.Offset, io.SeekStart)
	if err != nil {
		return nil
	}
	return f
}

func (d *downloadJob) Run(ctx context.Context) error {
	wc := d.tryResume()
	if wc == nil {
		f, err := ioutil.TempFile(tempDirPath, "download-")
		if err != nil {
			return err
		}
		defer f.Close()
		wc = f
		d.state = &stateType{
			Status:           statusDownloading,
			RemoteID:         d.remoteFile.PutioFile().ID,
			DownloadTempName: filepath.Base(f.Name()),
			Size:             d.remoteFile.PutioFile().Size,
			CRC32:            d.remoteFile.PutioFile().CRC32,
			relpath:          d.remoteFile.RelPath(),
		}
		err = d.state.Write()
		if err != nil {
			return err
		}
	}

	remaining := d.state.Size - d.state.Offset
	if remaining > 0 {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		rc, err := d.openRemote(ctx, d.state.Offset)
		if err != nil {
			return err
		}
		defer rc.Close()

		// Stop download if download speed is too slow.
		// Timer for cancelling the context will be reset after each successful read from stream.
		trw := &timerResetWriter{timer: time.AfterFunc(defaultTimeout, cancel)}
		tr := io.TeeReader(rc, trw)

		pr := progress.New(tr, d.state.Offset, d.state.Size, d.String())
		pr.Start()
		n, copyErr := io.CopyN(wc, pr, remaining)
		pr.Stop()

		err = wc.Close()
		if err != nil {
			return err
		}

		d.state.Offset += n
		err = d.state.Write()
		if err != nil {
			return err
		}

		if copyErr != nil {
			return copyErr
		}
	}

	oldPath := filepath.Join(tempDirPath, d.state.DownloadTempName)
	newPath := filepath.Join(localPath, filepath.FromSlash(d.state.relpath))
	err := os.MkdirAll(filepath.Dir(newPath), 0777)
	if err != nil {
		return err
	}
	err = os.Rename(oldPath, newPath)
	if err != nil {
		return err
	}

	in, err := inode.GetPath(newPath)
	if err != nil {
		return err
	}

	d.state.Status = statusSynced
	d.state.LocalInode = in
	return d.state.Write()
}

func (d *downloadJob) openRemote(baseCtx context.Context, offset int64) (rc io.ReadCloser, err error) {
	ctx, cancel := context.WithTimeout(baseCtx, defaultTimeout)
	defer cancel()
	u, err := client.Files.URL(ctx, d.remoteFile.PutioFile().ID, true)
	if err != nil {
		return
	}
	req, err := http.NewRequestWithContext(baseCtx, http.MethodGet, u, nil)
	if err != nil {
		return
	}
	req.Header.Set("range", fmt.Sprintf("bytes=%d-", offset))
	log.Debugln("range", req.Header.Get("range"))
	resp, err := httpClient.Do(req)
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusPartialContent {
		resp.Body.Close()
		err = fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		return
	}
	rc = resp.Body
	return
}

type timerResetWriter struct {
	timer *time.Timer
}

func (w *timerResetWriter) Write(p []byte) (int, error) {
	w.timer.Reset(defaultTimeout)
	return len(p), nil
}

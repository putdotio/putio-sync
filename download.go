package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

type Download struct {
	remoteFile *RemoteFile
	state      *State
}

func (d *Download) String() string {
	return fmt.Sprintf("Downloading %q", d.remoteFile.RelPath())
}

func (d *Download) resume() io.WriteCloser {
	if d.state == nil {
		return nil
	}
	if d.state.Status != StatusDownloading {
		return nil
	}
	if d.state.DownloadTempPath == "" {
		return nil
	}
	f, err := os.OpenFile(d.state.DownloadTempPath, os.O_WRONLY, 0)
	if err != nil {
		return nil
	}
	defer func() {
		if err != nil && f != nil {
			f.Close()
		}
	}()
	n, err := f.Seek(0, io.SeekEnd)
	if err != nil {
		return nil
	}
	if n < d.state.Offset {
		return nil
	}
	if d.state.Snapshot.ModTime != d.remoteFile.putioFile.UpdatedAt.Time {
		return nil
	}
	if d.state.Snapshot.Size != d.remoteFile.putioFile.Size {
		return nil
	}
	if d.state.Snapshot.CRC32 != d.remoteFile.putioFile.CRC32 {
		return nil
	}
	return f
}

func (d *Download) Run() error {
	wc := d.resume()
	if wc == nil {
		f, err := ioutil.TempFile("", "putio-sync-")
		if err != nil {
			return err
		}
		defer f.Close()
		wc = f
		d.state = &State{
			Status:           StatusDownloading,
			RemoteID:         d.remoteFile.putioFile.ID,
			DownloadTempPath: f.Name(),
			Snapshot: &Snapshot{
				Size:    d.remoteFile.putioFile.Size,
				ModTime: d.remoteFile.putioFile.UpdatedAt.Time,
				CRC32:   d.remoteFile.putioFile.CRC32,
			},
			relpath: d.remoteFile.relpath,
		}
		err = d.state.Write()
		if err != nil {
			return err
		}
	}
	rc, err := d.openRemote(d.state.Offset)
	if err != nil {
		return err
	}
	defer rc.Close()

	n, copyErr := io.Copy(wc, rc)

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
		return err
	}

	err = os.Rename(d.state.DownloadTempPath, filepath.Join(localPath, filepath.FromSlash(d.state.relpath)))
	if err != nil {
		return err
	}

	d.state.Status = StatusSynced
	err = d.state.Write()
	if err != nil {
		return err
	}

	return nil
}

func (d *Download) openRemote(offset int64) (rc io.ReadCloser, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	u, err := client.Files.URL(ctx, d.remoteFile.putioFile.ID, true)
	if err != nil {
		return
	}
	req, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, u, nil)
	if err != nil {
		return
	}
	req.Header.Set("range", fmt.Sprintf("bytes=%d-", offset))
	// TODO remove nolint directives
	resp, err := httpClient.Do(req) // nolint: bodyclose
	if err != nil {
		return
	}
	rc = resp.Body
	return
}

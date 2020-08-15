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

func (d *Download) tryResume() io.WriteCloser {
	if d.state == nil {
		return nil
	}
	if d.state.Status != StatusDownloading {
		return nil
	}
	if d.state.DownloadTempName == "" {
		return nil
	}
	f, err := os.OpenFile(filepath.Join(localPath, tempDirName, d.state.DownloadTempName), os.O_WRONLY, 0)
	if err != nil {
		return nil
	}
	defer func() {
		if err != nil && f != nil {
			f.Close()
		}
	}()
	if d.state.Size != d.remoteFile.putioFile.Size {
		return nil
	}
	if d.state.CRC32 != d.remoteFile.putioFile.CRC32 {
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

func (d *Download) Run(ctx context.Context) error {
	wc := d.tryResume()
	if wc == nil {
		tmpdir, err := TempDir()
		if err != nil {
			return err
		}
		f, err := ioutil.TempFile(tmpdir, "download-")
		if err != nil {
			return err
		}
		defer f.Close()
		wc = f
		d.state = &State{
			Status:           StatusDownloading,
			RemoteID:         d.remoteFile.putioFile.ID,
			DownloadTempName: filepath.Base(f.Name()),
			Size:             d.remoteFile.putioFile.Size,
			CRC32:            d.remoteFile.putioFile.CRC32,
			relpath:          d.remoteFile.relpath,
		}
		err = d.state.Write()
		if err != nil {
			return err
		}
	}
	rc, err := d.openRemote(ctx, d.state.Offset)
	if err != nil {
		return err
	}
	defer rc.Close()

	// TODO fail download if download stream is stuck

	remaining := d.state.Size - d.state.Offset
	n, copyErr := io.CopyN(wc, rc, remaining)

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

	oldPath := filepath.Join(localPath, tempDirName, d.state.DownloadTempName)
	newPath := filepath.Join(localPath, filepath.FromSlash(d.state.relpath))
	err = os.Rename(oldPath, newPath)
	if err != nil {
		return err
	}

	inode, err := GetInodePath(newPath)
	if err != nil {
		return err
	}

	d.state.Status = StatusSynced
	d.state.LocalInode = inode
	return d.state.Write()
}

func (d *Download) openRemote(baseCtx context.Context, offset int64) (rc io.ReadCloser, err error) {
	ctx, cancel := context.WithTimeout(baseCtx, defaultTimeout)
	defer cancel()
	u, err := client.Files.URL(ctx, d.remoteFile.putioFile.ID, true)
	if err != nil {
		return
	}
	req, err := http.NewRequestWithContext(baseCtx, http.MethodGet, u, nil)
	if err != nil {
		return
	}
	req.Header.Set("range", fmt.Sprintf("bytes=%d-", offset))
	resp, err := httpClient.Do(req)
	if err != nil {
		return
	}
	if resp.StatusCode != 206 {
		resp.Body.Close()
		err = fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		return
	}
	rc = resp.Body
	return
}

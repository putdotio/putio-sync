package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
)

type Upload struct {
	localFile *LocalFile
	state     *State
}

func (d *Upload) String() string {
	return fmt.Sprintf("Uploading %q", d.localFile.RelPath())
}

func (d *Upload) resume() bool {
	if d.state == nil {
		return false
	}
	if d.state.Status != StatusUploading {
		return false
	}
	if d.state.UploadURL == "" {
		return false
	}
	offset, err := GetUploadOffset(token, d.state.UploadURL)
	if err != nil {
		return false
	}
	d.state.Offset = offset
	if offset > d.localFile.info.Size() {
		return false
	}
	if d.state.Snapshot.ModTime != d.localFile.info.ModTime() {
		return false
	}
	if d.state.Snapshot.Size != d.localFile.info.Size() {
		return false
	}
	// TODO maybe check for CRC32
	return true
}

func (d *Upload) Run() error {
	ok := d.resume()
	if !ok {
		dir, filename := path.Split(d.localFile.RelPath())
		parentID, err := dirCache.Mkdirp(dir)
		if err != nil {
			return err
		}
		location, err := CreateUpload(context.TODO(), token, filename, parentID, d.localFile.info.Size())
		if err != nil {
			return err
		}
		d.state = &State{
			Status:    StatusUploading,
			UploadURL: location,
			Snapshot: &Snapshot{
				Size:    d.localFile.info.Size(),
				ModTime: d.localFile.info.ModTime(),
				// TODO maybe save and check inode number for local files
			},
			relpath: d.localFile.relpath,
		}
		err = d.state.Write()
		if err != nil {
			return err
		}
	}
	of, err := os.Open(filepath.Join(localPath, filepath.FromSlash(d.localFile.relpath)))
	if err != nil {
		return err
	}
	defer of.Close()
	_, err = of.Seek(d.state.Offset, io.SeekStart)
	if err != nil {
		return err
	}
	fileID, err := SendFile(token, of, d.state.UploadURL, d.state.Offset)
	if err != nil {
		return err
	}
	d.state.Status = StatusSynced
	d.state.RemoteID = fileID
	err = d.state.Write()
	if err != nil {
		return err
	}
	return nil
}
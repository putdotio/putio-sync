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

func (d *Upload) tryResume(ctx context.Context) bool {
	if d.state == nil {
		return false
	}
	if d.state.Status != StatusUploading {
		return false
	}
	if d.state.UploadURL == "" {
		return false
	}
	if d.state.Size != d.localFile.info.Size() {
		return false
	}
	inode, _ := GetInode(d.localFile.info)
	if d.state.LocalInode != inode {
		return false
	}
	offset, err := GetUploadOffset(ctx, token, d.state.UploadURL)
	if err != nil {
		return false
	}
	d.state.Offset = offset
	return offset <= d.localFile.info.Size()
}

func (d *Upload) Run(ctx context.Context) error {
	ok := d.tryResume(ctx)
	if !ok {
		inode, err := GetInode(d.localFile.info)
		if err != nil {
			return err
		}
		dir, filename := path.Split(d.localFile.RelPath())
		parentID, err := dirCache.Mkdirp(ctx, dir)
		if err != nil {
			return err
		}
		location, err := CreateUpload(ctx, token, filename, parentID, d.localFile.info.Size())
		if err != nil {
			return err
		}
		d.state = &State{
			Status:     StatusUploading,
			LocalInode: inode,
			UploadURL:  location,
			Size:       d.localFile.info.Size(),
			relpath:    d.localFile.relpath,
		}
		err = d.state.Write()
		if err != nil {
			return err
		}
	}
	f, err := os.Open(filepath.Join(localPath, filepath.FromSlash(d.localFile.relpath)))
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Seek(d.state.Offset, io.SeekStart)
	if err != nil {
		return err
	}
	fileID, crc32, err := SendFile(ctx, token, f, d.state.UploadURL, d.state.Offset)
	if err != nil {
		return err
	}
	d.state.Status = StatusSynced
	d.state.RemoteID = fileID
	d.state.CRC32 = crc32
	err = d.state.Write()
	if err != nil {
		return err
	}
	return nil
}

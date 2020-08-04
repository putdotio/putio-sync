package main

import (
	"context"
	"io"
	"net/http"

	"github.com/cenkalti/log"
)

type Download struct {
	remoteFile *RemoteFile
}

func NewDownload(rf File) *Download {
	return &Download{
		remoteFile: rf.(*RemoteFile),
	}

}

func (d *Download) String() string {
	return "download " + d.remoteFile.RelPath()
}

func (d *Download) Run() error {
	log.Infof("Downloading %q", d.remoteFile.RelPath())
	lf, err := CreateLocalFile(d.remoteFile.RelPath())
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	u, err := client.Files.URL(ctx, d.remoteFile.putioFile.ID, true)
	if err != nil {
		return err
	}
	// TODO do not use default http client
	// TODO use proper timeouts on client
	// TODO retry failed operations
	resp, err := http.Get(u)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// TODO get offset from database
	var offset int64

	rc, err := d.remoteFile.OpenForRead(offset)
	if err != nil {
		return err
	}
	defer rc.Close()

	wc, err := lf.OpenForWrite(offset)
	if err != nil {
		return err
	}
	defer wc.Close()

	_, err = io.Copy(wc, rc)
	if err != nil {
		return err
	}

	return wc.Close()
}

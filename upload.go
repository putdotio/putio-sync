package main

import (
	"github.com/cenkalti/log"
)

type Upload struct {
	localFile *LocalFile
}

func NewUpload(lf File) *Upload {
	return &Upload{
		localFile: lf.(*LocalFile),
	}

}

func (d *Upload) String() string {
	return "upload " + d.localFile.RelPath()
}

func (d *Upload) Run() error {
	log.Infof("Uploading %q", d.localFile.RelPath())
	// TODO implement upload
	return nil
}

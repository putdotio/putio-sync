package main

import (
	"fmt"
	"os"

	"github.com/cenkalti/log"
)

type Job interface {
	Run() error
	String() string
}

type DeleteState struct {
	state State
}

func (j *DeleteState) String() string {
	return fmt.Sprintf("Deleting state %q", j.state.relpath)
}

func (j *DeleteState) Run() error {
	if j.state.DownloadTempPath != "" {
		err := os.Remove(j.state.DownloadTempPath)
		if err != nil {
			log.Errorln("cannot remove temp download file:", err.Error())
		}
	}
	if j.state.UploadURL != "" {
		err := TerminateUpload(token, j.state.UploadURL)
		if err != nil {
			log.Errorln("cannot remove upload:", err.Error())
		}
	}
	return j.state.Delete()
}

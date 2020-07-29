package main

import (
	"github.com/cenkalti/log"
)

const folderName = "putio-sync"

func sync() error {
	log.Infof("Syncing https://put.io/files/%d with %q", remoteFolderID, localPath)
	err := CreateJobs()
	if err != nil {
		return err
	}
	// Print jobs for debugging
	for _, job := range jobs {
		log.Debugln("Job:", job.String())
	}
	// Run all jobs one by one
	for _, job := range jobs {
		err = job.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

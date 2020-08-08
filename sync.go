package main

import (
	"fmt"

	"github.com/cenkalti/log"
)

const folderName = "putio-sync"

func sync() error {
	remoteURL := fmt.Sprintf("https://put.io/files/%d", remoteFolderID)
	log.Infof("Syncing %q with %q", remoteURL, localPath)

	states, err := ReadAllStates()
	if err != nil {
		return err
	}
	// TODO walk on local and remote folders in parallel
	localFiles, err := WalkOnFolder(LocalWalker{})
	if err != nil {
		return err
	}
	remoteFiles, err := WalkOnFolder(RemoteWalker{})
	if err != nil {
		return err
	}
	for _, f := range remoteFiles {
		rf := f.(*RemoteFile)
		if rf.putioFile.IsDir() {
			dirCache.Set(rf.relpath, rf.putioFile.ID)
		}
	}
	syncFiles := GroupFiles(states, localFiles, remoteFiles)
	jobs := Reconciliation(syncFiles)
	// Print jobs for debugging
	for _, job := range jobs {
		log.Debugln("Job:", job.String())
	}
	dirCache.Debug()
	// Run all jobs one by one
	for _, job := range jobs {
		log.Infoln(job.String())
		err = job.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

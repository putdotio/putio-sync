package main

import (
	"fmt"
	"sync"

	"github.com/cenkalti/log"
)

const folderName = "putio-sync"

func syncRoots() error {
	remoteURL := fmt.Sprintf("https://put.io/files/%d", remoteFolderID)
	log.Infof("Syncing %q with %q", remoteURL, localPath)

	states, err := ReadAllStates()
	if err != nil {
		return err
	}

	// Walk on local and remote folders in parallel
	localFiles, remoteFiles, err := walkParallel()
	if err != nil {
		return err
	}

	for _, rf := range remoteFiles {
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

func walkParallel() ([]*LocalFile, []*RemoteFile, error) {
	var wg sync.WaitGroup
	localFilesC := make(chan []File, 1)
	remoteFilesC := make(chan []File, 1)
	errC := make(chan error, 2)
	wg.Add(2)
	go walkAsync(&wg, LocalWalker{}, "local", localFilesC, errC)
	go walkAsync(&wg, RemoteWalker{}, "remote", remoteFilesC, errC)
	wg.Wait()
	select {
	case err := <-errC:
		return nil, nil, err
	default:
	}
	files := <-localFilesC
	localFiles := make([]*LocalFile, 0, len(files))
	for _, f := range files {
		localFiles = append(localFiles, f.(*LocalFile))
	}
	files = <-remoteFilesC
	remoteFiles := make([]*RemoteFile, 0, len(files))
	for _, f := range files {
		remoteFiles = append(remoteFiles, f.(*RemoteFile))
	}
	return localFiles, remoteFiles, nil
}

func walkAsync(wg *sync.WaitGroup, walker Walker, side string, filesC chan []File, errC chan error) {
	files, err := WalkOnFolder(walker)
	if err != nil {
		errC <- err
	}
	log.Infof("Fetched %s filesystem tree", side)
	filesC <- files
	wg.Done()
}

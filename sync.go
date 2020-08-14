package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/cenkalti/log"
)

const folderName = "putio-sync"

func syncRoots(ctx context.Context) error {
	remoteURL := fmt.Sprintf("https://put.io/files/%d", remoteFolderID)
	log.Infof("Syncing %q with %q", remoteURL, localPath)

	// Read previous sync state from db.
	states, err := ReadAllStates()
	if err != nil {
		return err
	}

	// Walk on local and remote folders in parallel
	localFiles, remoteFiles, err := walkParallel(ctx)
	if err != nil {
		return err
	}

	// Set DirCache entries for existing remote folders
	for _, rf := range remoteFiles {
		if rf.putioFile.IsDir() {
			dirCache.Set(rf.relpath, rf.putioFile.ID)
		}
	}

	// Calculate what needs to be done
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
		err = job.Run(ctx)
		if err != nil {
			log.Error(err.Error())
		}
	}
	return nil
}

func walkParallel(ctx context.Context) ([]*LocalFile, []*RemoteFile, error) {
	var wg sync.WaitGroup
	localFilesC := make(chan []File, 1)
	remoteFilesC := make(chan []File, 1)
	errC := make(chan error, 2)
	wg.Add(2)
	go walkAsync(ctx, &wg, LocalWalker{}, "local", localFilesC, errC)
	go walkAsync(ctx, &wg, RemoteWalker{}, "remote", remoteFilesC, errC)
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

func walkAsync(ctx context.Context, wg *sync.WaitGroup, walker Walker, side string, filesC chan []File, errC chan error) {
	files, err := WalkOnFolder(ctx, walker)
	if err != nil {
		errC <- err
	}
	log.Infof("Fetched %s filesystem tree", side)
	filesC <- files
	wg.Done()
}

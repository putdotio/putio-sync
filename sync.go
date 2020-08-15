package main

import (
	"context"
	"fmt"

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
	if *dryrun {
		log.Noticeln("Command run in dry-run mode. No changes will be made.")
	}
	for _, job := range jobs {
		log.Infoln(job.String())
		if !*dryrun {
			err = job.Run(ctx)
			if err != nil {
				log.Error(err.Error())
			}
		}
	}
	return nil
}

func walkParallel(baseCtx context.Context) (localFiles []*LocalFile, remoteFiles []*RemoteFile, err error) {
	localFilesC := make(chan []File, 1)
	remoteFilesC := make(chan []File, 1)
	errC := make(chan error, 2)
	ctx, cancel := context.WithCancel(baseCtx)
	defer cancel()
	go walkAsync(ctx, LocalWalker{root: localPath}, localFilesC, errC)
	go walkAsync(ctx, RemoteWalker{root: remoteFolderID}, remoteFilesC, errC)
	for {
		if localFiles != nil && remoteFiles != nil {
			return localFiles, remoteFiles, nil
		}
		select {
		case files := <-localFilesC:
			log.Info("Fetched local filesystem tree")
			localFiles = make([]*LocalFile, 0, len(files))
			for _, f := range files {
				localFiles = append(localFiles, f.(*LocalFile))
			}
		case files := <-remoteFilesC:
			log.Info("Fetched remote filesystem tree")
			remoteFiles = make([]*RemoteFile, 0, len(files))
			for _, f := range files {
				remoteFiles = append(remoteFiles, f.(*RemoteFile))
			}
		case err = <-errC:
			// Cancel ongoing walk operation on first error
			cancel()
			return
		}
	}
}

func walkAsync(ctx context.Context, walker Walker, filesC chan []File, errC chan error) {
	files, err := WalkOnFolder(ctx, walker)
	if err != nil {
		errC <- err
		return
	}
	filesC <- files
}

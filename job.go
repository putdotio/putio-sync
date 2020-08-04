package main

import (
	"github.com/cenkalti/log"
)

type Job interface {
	Run() error
	String() string
}

func CreateJobs() error {
	jobs = nil

	// TODO walk on local and remote folders in parallel
	localFiles, err := WalkOnFolder(LocalWalker{})
	if err != nil {
		return err
	}
	printFiles(localFiles, "Local file: ")

	remoteFiles, err := WalkOnFolder(RemoteWalker{})
	if err != nil {
		return err
	}
	printFiles(remoteFiles, "Remote file:")

	localFilesByPath := mapFiles(localFiles)
	for _, rf := range remoteFiles {
		if lf, ok := localFilesByPath[rf.RelPath()]; ok {
			_ = lf
			_ = ok
			// if hasSameTimestamp(lf, rf) {
			// 	if hasSameHash(lf, rf) {
			// 		continue
			// 	}
			// }
		} else {
			jobs = append(jobs, NewDownload(rf))
		}
	}

	remoteFilesByPath := mapFiles(remoteFiles)
	for _, lf := range localFiles {
		if rf, ok := remoteFilesByPath[lf.RelPath()]; ok {
			_ = rf
			_ = ok
		} else {
			jobs = append(jobs, NewUpload(lf))
		}
	}
	return nil
}

func printFiles(files []File, prefix string) {
	for _, f := range files {
		log.Debugln(prefix, f.RelPath())
	}
}

func mapFiles(files []File) map[string]File {
	m := make(map[string]File)
	for _, f := range files {
		m[f.RelPath()] = f
	}
	return m
}

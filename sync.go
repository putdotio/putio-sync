package main

import (
	"fmt"

	"github.com/cenkalti/log"
)

const folderName = "putio-sync"

func sync() error {
	log.Infof("Syncing https://put.io/files/%d with %q", remoteFolderID, localPath)
	log.Infoln("printing local files")
	files, err := WalkOnFolder(LocalWalker{})
	if err != nil {
		return err
	}
	for _, file := range files {
		_ = printFile(file)
	}
	log.Infoln("printing remote files")
	files, err = WalkOnFolder(RemoteWalker{})
	if err != nil {
		return err
	}
	for _, file := range files {
		_ = printFile(file)
	}
	return nil
}

func printFile(info File) error {
	fmt.Println(info.Path(), info.Size())
	return nil
}

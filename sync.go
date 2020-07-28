package main

import (
	"fmt"
	"os"

	"github.com/cenkalti/log"
)

const folderName = "putio-sync"

func sync() error {
	log.Infof("Syncing https://put.io/files/%d with %q", remoteFolderID, localPath)
	log.Infoln("printing local files")
	err := WalkLocal(printFile)
	if err != nil {
		return err
	}
	log.Infoln("printing remote files")
	err = WalkRemote(printFile)
	if err != nil {
		return err
	}
	return nil
}

func printFile(path string, info os.FileInfo, err error) error {
	fmt.Println(path)
	return nil
}

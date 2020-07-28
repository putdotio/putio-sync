package main

import (
	"context"
	"fmt"

	"github.com/cenkalti/log"
)

const folderName = "putio-sync"

func sync() error {
	log.Infof("Syncing https://put.io/files/%d with %q", remoteFolderID, localPath)
	ai, err := client.Account.Info(context.TODO())
	if err != nil {
		return err
	}
	fmt.Printf("%#v\n", ai.UserID)
	return nil
}

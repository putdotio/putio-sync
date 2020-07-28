package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cenkalti/log"
	"github.com/putdotio/go-putio"
)

const folderName = "putio-sync"

func ensurePaths() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	localPath = filepath.Join(home, folderName)
	err = os.MkdirAll(localPath, 0750)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	folders, _, err := client.Files.List(ctx, 0)
	if err != nil {
		return err
	}
	found := false
	var f putio.File
	for _, f = range folders {
		if f.IsDir() && f.Name == folderName {
			found = true
			break
		}
	}
	if !found {
		ctx, cancel = context.WithTimeout(context.Background(), defaultTimeout)
		defer cancel()
		f, err = client.Files.CreateFolder(ctx, folderName, 0)
		if err != nil {
			return err
		}
	}
	remoteFolderID = f.ID
	return nil
}

func sync() error {
	log.Infof("Syncing https://put.io/files/%d with %q", remoteFolderID, localPath)
	ai, err := client.Account.Info(context.TODO())
	if err != nil {
		return err
	}
	fmt.Printf("%#v\n", ai.UserID)
	return nil
}

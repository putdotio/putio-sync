package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/putdotio/go-putio"
)

func ensureRoots(baseCtx context.Context) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	localPath = filepath.Join(home, folderName)
	err = os.MkdirAll(localPath, 0777)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(baseCtx, defaultTimeout)
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
		ctx, cancel = context.WithTimeout(baseCtx, defaultTimeout)
		defer cancel()
		f, err = client.Files.CreateFolder(ctx, folderName, 0)
		if err != nil {
			return err
		}
	}
	remoteFolderID = f.ID
	return nil
}

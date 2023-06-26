package putiosync

import (
	"context"
	"os"

	"github.com/putdotio/go-putio"
	"github.com/syncthing/syncthing/lib/fs"
)

const remoteFolderName = "putio-sync"

func ensureRoots(baseCtx context.Context) error {
	var err error
	localPath, err = fs.ExpandTilde(cfg.LocalDir)
	if err != nil {
		return err
	}
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
		if f.IsDir() && f.Name == remoteFolderName {
			found = true
			break
		}
	}
	if !found {
		ctx, cancel = context.WithTimeout(baseCtx, defaultTimeout)
		defer cancel()
		f, err = client.Files.CreateFolder(ctx, remoteFolderName, 0)
		if err != nil {
			return err
		}
	}
	remoteFolderID = f.ID
	return nil
}

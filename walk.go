package main

import (
	"context"
	"path"
	"path/filepath"

	"github.com/putdotio/go-putio"
)

func WalkLocal(walkFn filepath.WalkFunc) error {
	return filepath.Walk(localPath, walkFn)
}

func WalkRemote(walkFn filepath.WalkFunc) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	root, err := client.Files.Get(ctx, remoteFolderID)
	if err != nil {
		return walkFn("", nil, err)
	}
	return walk("/"+folderName, root, walkFn)
}

func walk(root string, dir putio.File, walkFn filepath.WalkFunc) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	children, parent, err := client.Files.List(ctx, remoteFolderID)
	if err != nil {
		return walkFn(root, nil, err)
	}
	err = walkFn(root, newFileInfo(parent), nil)
	if err != nil {
		return err
	}
	for _, child := range children {
		if !child.IsDir() {
			err = walkFn(root, newFileInfo(child), nil)
			if err != nil {
				return err
			}
		}
	}
	for _, child := range children {
		if child.IsDir() {
			err = walkFn(path.Join(root, child.Name), newFileInfo(child), nil)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

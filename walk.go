package main

import (
	"context"
	"os"
	"path"
	"path/filepath"
	"regexp"

	"github.com/putdotio/go-putio"
)

var ignoredFiles = regexp.MustCompile(`(?i)(^|/)(desktop\.ini|thumbs\.db|\.ds_store|icon\r)$`)

type Walker interface {
	Walk(walkFn filepath.WalkFunc) error
}

func WalkOnFolder(walker Walker) ([]File, error) {
	var l []File
	fn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if ignoredFiles.MatchString(info.Name()) {
			return nil
		}
		l = append(l, newFile(path, info))
		return nil
	}
	err := walker.Walk(fn)
	if err != nil {
		return nil, err
	}
	return l, nil
}

type LocalWalker struct{}

func (LocalWalker) Walk(walkFn filepath.WalkFunc) error {
	return filepath.Walk(localPath, walkFn)
}

type RemoteWalker struct{}

func (RemoteWalker) Walk(walkFn filepath.WalkFunc) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	dir, err := client.Files.Get(ctx, remoteFolderID)
	if err != nil {
		return walkFn("/", nil, err)
	}
	return walk("/", dir, walkFn)
}

func walk(root string, dir putio.File, walkFn filepath.WalkFunc) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	children, parent, err := client.Files.List(ctx, dir.ID)
	if err != nil {
		return walkFn(root, nil, err)
	}
	root = path.Join(root, parent.Name)
	err = walkFn(root, newFileInfo(parent), nil)
	if err != nil {
		return err
	}
	for _, child := range children {
		if !child.IsDir() {
			err = walkFn(path.Join(root, child.Name), newFileInfo(child), nil)
			if err != nil {
				return err
			}
		}
	}
	for _, child := range children {
		if child.IsDir() {
			err = walk(root, child, walkFn)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

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
	Walk(walkFn WalkFunc) error
}

type WalkFunc func(file File, err error) error

func WalkOnFolder(walker Walker) ([]File, error) {
	var l []File
	fn := func(file File, err error) error {
		if err != nil {
			return err
		}
		if file.RelPath() == "." {
			return nil
		}
		if ignoredFiles.MatchString(file.Info().Name()) {
			return nil
		}
		l = append(l, file)
		return nil
	}
	return l, walker.Walk(fn)
}

type LocalWalker struct{}

func (LocalWalker) Walk(walkFn WalkFunc) error {
	return filepath.Walk(localPath, func(path string, info os.FileInfo, err error) error {
		relpath, err2 := filepath.Rel(localPath, path)
		if err2 != nil {
			panic(err2)
		}
		return walkFn(NewLocalFile(info, relpath), err)
	})
}

type RemoteWalker struct{}

func (RemoteWalker) Walk(walkFn WalkFunc) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	// TODO list remote files sorted by datetime asc, this will make sure that latest uploaded version is always the most recent
	dir, err := client.Files.Get(ctx, remoteFolderID)
	if err != nil {
		return err
	}
	return walk(".", dir, walkFn)
}

func walk(relpath string, parent putio.File, walkFn WalkFunc) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	children, _, err := client.Files.List(ctx, parent.ID)
	if err != nil {
		return walkFn(nil, err)
	}
	err = walkFn(NewRemoteFile(parent, relpath), nil)
	if err != nil {
		return err
	}
	for _, child := range children {
		if !child.IsDir() {
			err = walkFn(NewRemoteFile(child, path.Join(relpath, child.Name)), nil)
			if err != nil {
				return err
			}
		}
	}
	for _, child := range children {
		if child.IsDir() {
			err = walk(path.Join(relpath, child.Name), child, walkFn)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

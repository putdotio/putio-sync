package main

import (
	"context"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"

	"github.com/putdotio/go-putio"
)

var ignoredFiles = regexp.MustCompile(`(?i)(^|/)(desktop\.ini|thumbs\.db|\.ds_store|icon\r)$`)

type Walker interface {
	Walk(walkFn WalkFunc) error
}

type WalkFunc func(file File, err error) error

func WalkOnFolder(ctx context.Context, walker Walker) ([]File, error) {
	var l []File
	fn := func(file File, err error) error {
		if err != nil {
			return err
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
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

type LocalWalker struct {
	root string
}

func (w LocalWalker) Walk(walkFn WalkFunc) error {
	return filepath.Walk(w.root, func(path string, info os.FileInfo, err error) error {
		relpath, err2 := filepath.Rel(w.root, path)
		if err2 != nil {
			panic(err2)
		}
		return walkFn(NewLocalFile(info, relpath), err)
	})
}

type RemoteWalker struct {
	root int64
}

func (w RemoteWalker) Walk(walkFn WalkFunc) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	dir, err := client.Files.Get(ctx, w.root)
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
	// List remote files sorted by ID ascending, this will make sure that latest uploaded version is always the most recent
	sort.Slice(children, func(i, j int) bool {
		return children[i].ID < children[j].ID
	})
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

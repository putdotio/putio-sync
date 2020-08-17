package walker

import (
	"context"
	"path"
	"sort"
	"time"

	"github.com/putdotio/go-putio"
)

type remoteWalker struct {
	root           int64
	client         *putio.Client
	requestTimeout time.Duration
}

func (w *remoteWalker) Walk(walkFn walkFunc) error {
	ctx, cancel := context.WithTimeout(context.Background(), w.requestTimeout)
	defer cancel()
	dir, err := w.client.Files.Get(ctx, w.root)
	if err != nil {
		return err
	}
	return w.walk(".", dir, walkFn)
}

func (w *remoteWalker) walk(relpath string, parent putio.File, walkFn walkFunc) error {
	ctx, cancel := context.WithTimeout(context.Background(), w.requestTimeout)
	defer cancel()
	children, _, err := w.client.Files.List(ctx, parent.ID)
	if err != nil {
		return walkFn(nil, err)
	}
	err = walkFn(newRemoteFile(parent, relpath), nil)
	if err != nil {
		return err
	}
	// List remote files sorted by ID ascending, this will make sure that latest uploaded version is always the most recent
	sort.Slice(children, func(i, j int) bool {
		return children[i].ID < children[j].ID
	})
	for _, child := range children {
		if !child.IsDir() {
			err = walkFn(newRemoteFile(child, path.Join(relpath, child.Name)), nil)
			if err != nil {
				return err
			}
		}
	}
	for _, child := range children {
		if child.IsDir() {
			err = w.walk(path.Join(relpath, child.Name), child, walkFn)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

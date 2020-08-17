package main

import (
	"context"
	"path"
	"strings"

	"github.com/cenkalti/log"
)

// DirCache holds a map for accessing IDs by path.
type DirCache struct {
	m map[string]int64
}

func NewDirCache() *DirCache {
	return &DirCache{
		m: make(map[string]int64),
	}
}

func (c *DirCache) Debug() {
	for k, v := range c.m {
		log.Debugln("DirCache", k, v)
	}
}

func (c *DirCache) Clear() {
	c.m = make(map[string]int64)
}

func (c *DirCache) Set(relpath string, id int64) {
	relpath = strings.TrimRight(relpath, "/")
	c.m[relpath] = id
}

func (c *DirCache) Mkdirp(ctx context.Context, relpath string) (int64, error) {
	relpath = strings.TrimRight(relpath, "/")
	log.Debugln("DirCache.Mkdirp", relpath)
	if relpath == "." || relpath == "" {
		return remoteFolderID, nil
	}
	if id, ok := c.m[relpath]; ok {
		return id, nil
	}
	dir, base := path.Split(relpath)
	dirID, err := c.Mkdirp(ctx, dir)
	if err != nil {
		return 0, err
	}
	log.Debugf("DirCache.Mkdirp Creating remote folder %q", relpath)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	f, err := client.Files.CreateFolder(ctx, base, dirID)
	if err != nil {
		return 0, err
	}
	c.m[relpath] = f.ID
	return f.ID, nil
}

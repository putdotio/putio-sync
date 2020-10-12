package watcher

import (
	"context"

	"github.com/cenkalti/log"
	"github.com/fsnotify/fsnotify"
)

type FileModificationWatcher struct {
	watcher  *fsnotify.Watcher
	ctx      context.Context
	cancel   func()
	modified bool
	done     chan interface{}
}

func WatchFileModification(ctx context.Context, file string) (*FileModificationWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	err = watcher.Add(file)
	if err != nil {
		watcher.Close()
		return nil, err
	}
	w := &FileModificationWatcher{
		watcher: watcher,
		done:    make(chan interface{}),
	}
	w.ctx, w.cancel = context.WithCancel(ctx)
	go w.run()
	return w, nil
}

func (w *FileModificationWatcher) Context() context.Context {
	return w.ctx
}

func (w *FileModificationWatcher) Stop() bool {
	w.cancel()
	w.watcher.Close()
	<-w.done
	return w.modified
}

func (w *FileModificationWatcher) run() {
	const mask = fsnotify.Write | fsnotify.Remove | fsnotify.Rename
	defer close(w.done)
	for {
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}
			if event.Op&mask != 0 {
				w.modified = true
				w.cancel()
				return
			}
		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			log.Errorln("File watch error:", err)
		case <-w.ctx.Done():
			return
		}
	}
}

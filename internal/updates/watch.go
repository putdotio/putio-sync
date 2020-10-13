package updates

import (
	"context"
)

type FileWatcher struct {
	Stop     func() bool
	fileID   int64
	ctx      context.Context
	cancel   func()
	modified bool
}

func newFileWatcher(ctx context.Context, fileID int64, stopFn func() bool) *FileWatcher {
	ctx, cancel := context.WithCancel(ctx)
	return &FileWatcher{
		Stop:   stopFn,
		fileID: fileID,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (w *FileWatcher) Context() context.Context {
	return w.ctx
}

func (w *FileWatcher) notify(id int64) {
	if w == nil {
		return
	}
	if id != w.fileID {
		return
	}
	w.modified = true
	w.cancel()
}

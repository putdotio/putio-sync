package watcher

import (
	"context"

	"github.com/cenkalti/log"
	"github.com/fsnotify/fsnotify"
)

const mask = fsnotify.Create | fsnotify.Write | fsnotify.Remove | fsnotify.Rename

func Watch(ctx context.Context, dir string) (chan struct{}, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	err = watcher.Add(dir)
	if err != nil {
		watcher.Close()
		return nil, err
	}

	ch := make(chan struct{}, 1)
	go processEvents(ctx, watcher, ch)
	return ch, nil
}

func processEvents(ctx context.Context, watcher *fsnotify.Watcher, ch chan struct{}) {
	defer log.Debugln("end process events")
	defer watcher.Close()
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&mask != 0 {
				logEvent(event)
				select {
				case ch <- struct{}{}:
				case <-ctx.Done():
					return
				default:
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Errorln("fsnotify error:", err)
		case <-ctx.Done():
			return
		}
	}
}

func logEvent(event fsnotify.Event) {
	log.Debugf("Event Name: %s Op: %s", event.Name, event.Op.String())
}

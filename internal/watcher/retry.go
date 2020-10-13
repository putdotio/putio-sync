package watcher

import (
	"context"
	"path/filepath"
	"strings"
	"time"

	"github.com/cenkalti/log"
	"github.com/putdotio/putio-sync/v2/internal/tmpdir"
	"github.com/putdotio/putio-sync/v2/internal/walker"
)

func retry(ctx context.Context, dir string, watchFn func(ctx context.Context, dir string) (chan string, error)) (chan string, error) {
	in, err := watchFn(ctx, dir)
	if err != nil {
		return nil, err
	}

	// watch started successfully. Wait for channel close event for errors and restart watching.
	out := make(chan string, 1)
	go func() {
		for {
			select {
			case event, ok := <-in:
				if !ok {
					in, err = watchFn(ctx, dir)
					if err != nil {
						log.Error(err)
						select {
						case <-time.After(time.Second):
						case <-ctx.Done():
							return
						}
					}
					break
				}

				// This is not the correct place for filtering path names,
				// but for now it is okay because this `retry` function is used in all implementations.
				if walker.Ignored(filepath.Base(event)) || strings.Contains(event, tmpdir.Name) {
					continue
				}

				// Forward the event to returned channel.
				select {
				case out <- event:
				case <-ctx.Done():
					return
				default:
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return out, nil
}

package watcher

import (
	"context"
	"time"

	"github.com/cenkalti/log"
)

func retry(ctx context.Context, dir string, watchFn func(ctx context.Context, dir string) (chan struct{}, error)) (chan struct{}, error) {
	in, err := watchFn(ctx, dir)
	if err != nil {
		return nil, err
	}

	// watch started successfully. Wait for channel close event for errors and restart watching.
	out := make(chan struct{}, 1)
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

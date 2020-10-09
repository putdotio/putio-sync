package watcher

import (
	"context"
	"time"

	"github.com/cenkalti/log"
	"github.com/fsnotify/fsevents"
)

func Watch(ctx context.Context, dir string) (chan struct{}, error) {
	return retry(ctx, dir, watch)
}

func watch(ctx context.Context, dir string) (chan struct{}, error) {
	ch := make(chan struct{}, 1)

	dev, err := fsevents.DeviceForPath(dir)
	if err != nil {
		return nil, err
	}
	log.Debugln("Device:", dev)
	log.Debugln("Event ID:", fsevents.EventIDForDeviceBeforeTime(dev, time.Now()))

	es := &fsevents.EventStream{
		Paths:   []string{dir},
		Latency: 500 * time.Millisecond,
		Device:  dev,
		Flags:   fsevents.FileEvents | fsevents.WatchRoot,
	}

	es.Start()
	go processEvents(ctx, es, ch)
	return ch, nil
}

const mask = fsevents.ItemCreated | fsevents.ItemRemoved | fsevents.ItemRenamed

func processEvents(ctx context.Context, es *fsevents.EventStream, ch chan struct{}) {
	defer log.Debugln("end process events")
	defer close(ch)
	defer es.Stop()
	events := es.Events
	for {
		select {
		case msg, ok := <-events:
			if !ok {
				return
			}
			for _, event := range msg {
				logEvent(event)
				if event.Flags&mask != 0 {
					select {
					case ch <- struct{}{}:
					case <-ctx.Done():
						return
					default:
					}
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

var noteDescription = map[fsevents.EventFlags]string{
	fsevents.MustScanSubDirs: "MustScanSubdirs",
	fsevents.UserDropped:     "UserDropped",
	fsevents.KernelDropped:   "KernelDropped",
	fsevents.EventIDsWrapped: "EventIDsWrapped",
	fsevents.HistoryDone:     "HistoryDone",
	fsevents.RootChanged:     "RootChanged",
	fsevents.Mount:           "Mount",
	fsevents.Unmount:         "Unmount",

	fsevents.ItemCreated:       "Created",
	fsevents.ItemRemoved:       "Removed",
	fsevents.ItemInodeMetaMod:  "InodeMetaMod",
	fsevents.ItemRenamed:       "Renamed",
	fsevents.ItemModified:      "Modified",
	fsevents.ItemFinderInfoMod: "FinderInfoMod",
	fsevents.ItemChangeOwner:   "ChangeOwner",
	fsevents.ItemXattrMod:      "XAttrMod",
	fsevents.ItemIsFile:        "IsFile",
	fsevents.ItemIsDir:         "IsDir",
	fsevents.ItemIsSymlink:     "IsSymLink",
}

func logEvent(event fsevents.Event) {
	note := ""
	for bit, description := range noteDescription {
		if event.Flags&bit == bit {
			note += description + " "
		}
	}
	log.Debugf("Event ID: %d Path: %s Flags: %s", event.ID, event.Path, note)
}

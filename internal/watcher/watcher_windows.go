package watcher

import (
	"context"
	"os"
	"syscall"
	"time"
	"unsafe"

	"github.com/cenkalti/log"
)

// Buffer size cannot exceed 64K
const bufferSize = 32 << 10

const mask = syscall.FILE_NOTIFY_CHANGE_SIZE | syscall.FILE_NOTIFY_CHANGE_FILE_NAME | syscall.FILE_NOTIFY_CHANGE_DIR_NAME

func Watch(ctx context.Context, dir string) (chan struct{}, error) {
	in, err := watch(ctx, dir)
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
					in, err = watch(ctx, dir)
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

func watch(ctx context.Context, dir string) (chan struct{}, error) {
	var err error
	var overlapped syscall.Overlapped
	buffer := make([]byte, bufferSize)

	pdir, err := syscall.UTF16PtrFromString(dir)
	if err != nil {
		return nil, os.NewSyscallError("UTF16PtrFromString", err)
	}

	dh, err := syscall.CreateFile(
		pdir,
		syscall.FILE_LIST_DIRECTORY,
		syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE|syscall.FILE_SHARE_DELETE,
		nil,
		syscall.OPEN_EXISTING,
		syscall.FILE_FLAG_BACKUP_SEMANTICS|syscall.FILE_FLAG_OVERLAPPED,
		0,
	)
	if err != nil {
		return nil, os.NewSyscallError("CreateFile", err)
	}
	defer func() {
		if err != nil {
			_ = syscall.CloseHandle(dh)
		}
	}()

	cph, err := syscall.CreateIoCompletionPort(dh, 0, 0, 0)
	if err != nil {
		return nil, os.NewSyscallError("CreateIoCompletionPort", err)
	}
	defer func() {
		if err != nil {
			_ = syscall.CloseHandle(cph)
		}
	}()

	err = readDirChanges(dh, buffer, &overlapped)
	if err != nil {
		return nil, os.NewSyscallError("ReadDirectoryChanges", err)
	}

	done := make(chan struct{}) // will be closed when processEvents ends
	go closeHandles(ctx, dh, cph, done)

	ch := make(chan struct{}, 1)
	go processEvents(ctx, dh, cph, buffer, &overlapped, ch, done)

	return ch, nil
}

func closeHandles(ctx context.Context, dh, cph syscall.Handle, done chan struct{}) {
	select {
	case <-done: // processEvents returned
	case <-ctx.Done():
	}
	_ = syscall.CloseHandle(dh)
	_ = syscall.CloseHandle(cph)
}

func processEvents(ctx context.Context, dh, cph syscall.Handle, buffer []byte, overlapped *syscall.Overlapped, ch, done chan struct{}) {
	defer log.Debugln("end process events")
	defer close(done)
	defer close(ch)

	var n, key uint32
	var ov *syscall.Overlapped

	for {
		err := syscall.GetQueuedCompletionStatus(cph, &n, &key, &ov, syscall.INFINITE)
		select {
		case <-ctx.Done():
			return
		default:
		}
		if err != nil {
			log.Errorln("GetQueuedCompletionPort", err)
			return
		}

		var offset uint32
		for {
			event := (*syscall.FileNotifyInformation)(unsafe.Pointer(&buffer[offset]))
			buf := (*[syscall.MAX_PATH]uint16)(unsafe.Pointer(&event.FileName))
			name := syscall.UTF16ToString(buf[:event.FileNameLength/2])

			logEvent(event, name)
			select {
			case ch <- struct{}{}:
			case <-ctx.Done():
				return
			default:
			}

			if event.NextEntryOffset == 0 {
				break
			}

			offset += event.NextEntryOffset
			if offset >= n {
				log.Error("Windows system assumed buffer larger than it is, events have likely been missed.")
				break
			}
		}

		err = readDirChanges(dh, buffer, overlapped)
		select {
		case <-ctx.Done():
			return
		default:
		}
		if err != nil {
			log.Error(os.NewSyscallError("ReadDirectoryChanges", err))
			return
		}
	}
}

func readDirChanges(h syscall.Handle, buf []byte, ov *syscall.Overlapped) error {
	return syscall.ReadDirectoryChanges(
		h,
		&buf[0],
		uint32(len(buf)),
		true, // bWatchSubtree
		mask,
		nil,
		(*syscall.Overlapped)(unsafe.Pointer(ov)),
		0,
	)
}

func logEvent(event *syscall.FileNotifyInformation, name string) {
	act := "OTHER"
	switch event.Action {
	case syscall.FILE_ACTION_ADDED:
		act = "ADDED"
	case syscall.FILE_ACTION_MODIFIED:
		act = "MODIFIED"
	case syscall.FILE_ACTION_REMOVED:
		act = "REMOVED"
	case syscall.FILE_ACTION_RENAMED_NEW_NAME | syscall.FILE_ACTION_RENAMED_OLD_NAME:
		act = "RENAMED"
	}
	log.Debugf("Event Path: %s Action: %s", name, act)
}

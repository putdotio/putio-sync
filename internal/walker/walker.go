package walker

import (
	"context"
	"regexp"
	"strings"
	"time"

	"github.com/cenkalti/log"
	"github.com/putdotio/go-putio"
)

var ignoredFiles = regexp.MustCompile(`(?i)(^|/)(desktop\.ini|thumbs\.db|\.ds_store|icon\r)$`)

type walker interface {
	Walk(walkFn walkFunc) error
}

type walkFunc func(file file, err error) error

type Walker struct {
	LocalPath      string
	RemoteFolderID int64
	TempDirName    string
	Client         *putio.Client
	RequestTimeout time.Duration
}

func (w *Walker) Walk(ctx context.Context) (localFiles []*LocalFile, remoteFiles []*RemoteFile, err error) {
	localFilesC := make(chan []file, 1)
	remoteFilesC := make(chan []file, 1)
	errC := make(chan error, 2)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go w.walkAsync(ctx, &localWalker{root: w.LocalPath}, localFilesC, errC)
	go w.walkAsync(ctx, &remoteWalker{root: w.RemoteFolderID, client: w.Client, requestTimeout: w.RequestTimeout}, remoteFilesC, errC)
	for {
		if localFiles != nil && remoteFiles != nil {
			return localFiles, remoteFiles, nil
		}
		select {
		case files := <-localFilesC:
			log.Debug("Fetched local filesystem tree")
			localFiles = make([]*LocalFile, 0, len(files))
			for _, f := range files {
				localFiles = append(localFiles, f.(*LocalFile))
			}
		case files := <-remoteFilesC:
			log.Debug("Fetched remote filesystem tree")
			remoteFiles = make([]*RemoteFile, 0, len(files))
			for _, f := range files {
				remoteFiles = append(remoteFiles, f.(*RemoteFile))
			}
		case err = <-errC:
			// Cancel ongoing walk operation on first error
			cancel()
			return
		}
	}
}

func (w *Walker) walkAsync(ctx context.Context, walker walker, filesC chan []file, errC chan error) {
	files, err := w.walkOnFolder(ctx, walker)
	if err != nil {
		errC <- err
		return
	}
	filesC <- files
}

func (w *Walker) walkOnFolder(ctx context.Context, walker walker) ([]file, error) {
	var l []file
	fn := func(file file, err error) error {
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
		if strings.HasPrefix(file.RelPath(), w.TempDirName) {
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

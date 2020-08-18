package putiosync

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/adrg/xdg"
	"github.com/cenkalti/log"
	"github.com/putdotio/go-putio"
	"github.com/putdotio/putio-sync/v2/internal/auth"
	"github.com/putdotio/putio-sync/v2/internal/dircache"
	"github.com/putdotio/putio-sync/v2/internal/tmpdir"
	"github.com/putdotio/putio-sync/v2/internal/tus"
	"github.com/putdotio/putio-sync/v2/internal/walker"
	"go.etcd.io/bbolt"
)

const (
	folderName     = "putio-sync"
	defaultTimeout = 10 * time.Second
)

// Variables that used by Sync function.
var (
	cfg            Config
	db             *bbolt.DB
	token          string
	client         *putio.Client
	localPath      string
	remoteFolderID int64
	dirCache       *dircache.DirCache
	tempDirPath    string
	uploader       *tus.Uploader
)

func Sync(ctx context.Context, config Config) error {
	if err := config.validate(); err != nil {
		return err
	}
	dbPath, err := xdg.DataFile(filepath.Join("putio-sync", "sync.db"))
	if err != nil {
		return err
	}
	log.Infof("Using database file %q", dbPath)
	db, err = bbolt.Open(dbPath, 0666, nil)
	if err != nil {
		return err
	}
	defer db.Close()
	cfg = config
	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketFiles)
		return err
	})
	if err != nil {
		return err
	}
	var srv *httpServer
	if cfg.Server != "" {
		srv = newServer(cfg.Server)
		srv.Start()
		defer srv.Close()
	}

REPEAT_LOOP:
	for {
		err = syncOnce(ctx)
		if err != nil {
			if cfg.Repeat == 0 {
				return err
			}
			log.Error(err)
		} else {
			log.Infoln("Sync finished successfully")
		}
		if cfg.Repeat == 0 {
			break
		}
		select {
		case <-time.After(cfg.Repeat):
		case <-ctx.Done():
			break REPEAT_LOOP
		}
	}
	if srv != nil {
		if err := srv.Shutdown(); err != nil {
			return err
		}
		log.Debug("Server has shutdown successfully")
	}
	return nil
}

func syncOnce(ctx context.Context) error {
	var err error
	token, client, err = auth.Authenticate(ctx, httpClient, defaultTimeout, cfg.Username, cfg.Password)
	if err != nil {
		return err
	}
	err = ensureRoots(ctx)
	if err != nil {
		return err
	}
	tempDirPath, err = tmpdir.Create(localPath)
	if err != nil {
		return err
	}
	dirCache = dircache.New(client, defaultTimeout, remoteFolderID)
	uploader = tus.NewUploader(httpClient, defaultTimeout, token)

	return syncRoots(ctx)
}

func syncRoots(ctx context.Context) error {
	remoteURL := fmt.Sprintf("https://put.io/files/%d", remoteFolderID)
	log.Infof("Syncing %q with %q", remoteURL, localPath)

	// Read previous sync state from db.
	states, err := readAllStates()
	if err != nil {
		return err
	}

	// Walk on local and remote folders in parallel
	w := walker.Walker{
		LocalPath:      localPath,
		RemoteFolderID: remoteFolderID,
		TempDirName:    tmpdir.Name,
		Client:         client,
		RequestTimeout: defaultTimeout,
	}
	localFiles, remoteFiles, err := w.Walk(ctx)
	if err != nil {
		return err
	}

	// Set DirCache entries for existing remote folders
	for _, rf := range remoteFiles {
		if rf.PutioFile().IsDir() {
			dirCache.Set(rf.RelPath(), rf.PutioFile().ID)
		}
	}

	// Calculate what needs to be done
	syncFiles := groupFiles(states, localFiles, remoteFiles)
	jobs := reconciliation(syncFiles)

	// Print jobs for debugging
	for _, job := range jobs {
		log.Debugln("Job:", job.String())
	}
	dirCache.Debug()

	// Run all jobs one by one
	if cfg.DryRun {
		log.Noticeln("Command run in dry-run mode, no changes will be made")
	}
	if len(jobs) == 0 {
		log.Infoln("No changes detected")
		return nil
	}
	for _, job := range jobs {
		log.Infoln(job.String())
		if cfg.DryRun {
			continue
		}
		err = job.Run(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

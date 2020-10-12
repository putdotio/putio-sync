package putiosync

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"time"

	"github.com/adrg/xdg"
	"github.com/cenkalti/log"
	"github.com/putdotio/go-putio"
	"github.com/putdotio/putio-sync/v2/internal/auth"
	"github.com/putdotio/putio-sync/v2/internal/dircache"
	"github.com/putdotio/putio-sync/v2/internal/tmpdir"
	"github.com/putdotio/putio-sync/v2/internal/updates"
	"github.com/putdotio/putio-sync/v2/internal/walker"
	"github.com/putdotio/putio-sync/v2/internal/watcher"
	"go.etcd.io/bbolt"
)

const (
	folderName     = "putio-sync"
	defaultTimeout = 10 * time.Second
)

var ErrInvalidCredentials = errors.New("invalid credentials")

// Variables that used by Sync function.
var (
	cfg            Config
	db             *bbolt.DB
	token          string
	client         *putio.Client
	notifier       = updates.NewNotifier("wss://socket.put.io/socket/sockjs/websocket", 10*time.Second, 5*time.Second)
	watcherUpdates chan string
	localPath      string
	remoteFolderID int64
	dirCache       *dircache.DirCache
	tempDirPath    string
	syncing        bool
	syncStatus     = "Starting sync..."
)

func Sync(ctx context.Context, config Config) error {
	if err := config.validate(); err != nil {
		return err
	}
	if config.Debug {
		log.SetLevel(log.DEBUG)
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

	for {
		err = syncOnce(ctx)
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return ErrInvalidCredentials
		}
		if err != nil {
			if cfg.Repeat == 0 {
				return err
			}
			log.Error(err)
		} else {
			syncStatus = "Sync finished successfully"
			log.Infoln(syncStatus)
		}
		ok := waitNextSync(ctx)
		if !ok {
			break
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
	if cfg.Repeat != 0 {
		notifier.SetToken(token)
		notifier.Start()
	}
	if watcherUpdates == nil {
		watcherUpdates, err = watcher.Watch(ctx, localPath)
		if err != nil {
			log.Error(err)
		}
	}
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
	// dirCache.Debug()

	// Run all jobs one by one
	if cfg.DryRun {
		log.Noticeln("Command run in dry-run mode, no changes will be made")
	}
	if len(jobs) == 0 {
		log.Infoln("No changes detected")
		return nil
	}
	syncing = true
	defer func() { syncing = false }()
	for _, job := range jobs {
		syncStatus = job.String()
		log.Infoln(syncStatus)
		if cfg.DryRun {
			continue
		}
		err = job.Run(ctx)
		if err != nil {
			syncStatus = "Error: " + err.Error()
			return err
		}
	}
	return nil
}

func waitNextSync(ctx context.Context) bool {
	if cfg.Repeat == 0 {
		return false
	}
	var tc <-chan time.Time
	startTimer := func() {
		if tc == nil {
			tc = time.After(5 * time.Second)
		}
	}
	var d time.Duration
	if notifier.Connected() && runtime.GOOS != "linux" {
		d = 2 * time.Hour
	} else {
		d = 15 * time.Minute
	}
	for {
		select {
		case <-time.After(d):
			return true
		case name := <-notifier.HasUpdates:
			log.Infof("Change detected at remote filesystem: %q", name)
			startTimer()
		case name := <-watcherUpdates:
			log.Infof("Change detected at local filesystem: %q", name)
			startTimer()
		case <-tc:
			return true
		case <-ctx.Done():
			return false
		}
	}
}

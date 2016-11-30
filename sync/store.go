package sync

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"os/user"
	"path/filepath"
	"time"

	"github.com/boltdb/bolt"
	"github.com/igungor/go-putio/putio"
)

// buckets
var (
	downloadItemsBucket   = []byte("download-items")
	watchedTorrentsBucket = []byte("watched-torrents")
	userAccountBucket     = []byte("user-account")
	configBucket          = []byte("configuration")
)

type Error string

func (e Error) Error() string { return string(e) }

const (
	ErrStateNotFound   = Error("state not found")
	ErrConfigNotFound  = Error("configuration not found")
	ErrSaveStateFailed = Error("state could not be saved")
)

type Store struct {
	path string
	db   *bolt.DB
}

func NewStore(path string) *Store {
	return &Store{path: path}
}

func (s *Store) Open() error {
	db, err := bolt.Open(s.path, 0666, &bolt.Options{Timeout: 10 * time.Second})
	if err != nil {
		return err
	}
	s.db = db

	err = s.db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists(downloadItemsBucket)
		_, err = tx.CreateBucketIfNotExists(watchedTorrentsBucket)
		_, err = tx.CreateBucketIfNotExists(userAccountBucket)
		_, err = tx.CreateBucketIfNotExists(configBucket)
		return err
	})
	if err != nil {
		return s.db.Close()
	}
	return nil
}

func (s *Store) Close() error { return s.db.Close() }

func (s *Store) Path() string { return s.path }

// SaveState inserts or updates the given state.
func (s *Store) SaveState(state *State) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(downloadItemsBucket)
		key := itob(state.FileID)

		var value bytes.Buffer
		err := gob.NewEncoder(&value).Encode(state)
		if err != nil {
			return err
		}

		return bucket.Put(key, value.Bytes())
	})
}

// State returns a state by the given file ID.
func (s *Store) State(id int64) (*State, error) {
	var state State
	err := s.db.View(func(tx *bolt.Tx) error {
		fileid := itob(id)

		value := tx.Bucket(downloadItemsBucket).Get(fileid)
		if value == nil {
			return ErrStateNotFound
		}

		return gob.NewDecoder(bytes.NewReader(value)).Decode(&state)
	})
	return &state, err
}

// States returns all the states in the store.
func (s *Store) States() ([]*State, error) {
	states := make([]*State, 0)
	err := s.db.View(func(tx *bolt.Tx) error {
		cursor := tx.Bucket(downloadItemsBucket).Cursor()

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var state State
			err := gob.NewDecoder(bytes.NewReader(v)).Decode(&state)
			if err != nil {
				return err
			}
			// dont include hidden downloads
			if state.IsHidden {
				continue
			}
			states = append(states, &state)
		}
		return nil
	})

	return states, err
}

// StateN returns the number of states in the store.
func (s *Store) StateN() (int, error) {
	var n int
	err := s.db.View(func(tx *bolt.Tx) error {
		cursor := tx.Bucket(downloadItemsBucket).Cursor()
		for k, _ := cursor.First(); k != nil; k, _ = cursor.Next() {
			n++
		}
		return nil
	})
	return n, err
}

func (s *Store) SaveUserAccount(info putio.AccountInfo) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(userAccountBucket)
		key := []byte(info.Username)
		var value bytes.Buffer

		err := gob.NewEncoder(&value).Encode(info)
		if err != nil {
			return err
		}

		return bucket.Put(key, value.Bytes())
	})
}

func (s *Store) Config() (*Config, error) {
	var cfg Config
	err := s.db.View(func(tx *bolt.Tx) error {
		key := []byte("config")
		value := tx.Bucket(configBucket).Get(key)
		if value == nil {
			u, err := user.Current()
			if err != nil {
				return err
			}

			cfg = Config{
				PollInterval:        Duration(defaultPollInterval),
				DownloadTo:          filepath.Join(u.HomeDir, "putio-sync"),
				DownloadFrom:        defaultDownloadFrom,
				SegmentsPerFile:     defaultSegmentsPerFile,
				MaxParallelFiles:    defaultMaxParallelFiles,
				WatchTorrentsFolder: false,
				TorrentsFolder:      "",
				IsPaused:            true,
			}
			return nil
		}

		return gob.NewDecoder(bytes.NewReader(value)).Decode(&cfg)
	})

	return &cfg, err
}

func (s *Store) SaveConfig(cfg *Config) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(configBucket)
		key := []byte("config")
		var value bytes.Buffer

		err := gob.NewEncoder(&value).Encode(cfg)
		if err != nil {
			return err
		}

		return bucket.Put(key, value.Bytes())
	})
}

func itob(v int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

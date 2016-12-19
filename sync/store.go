package sync

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"os/user"
	"path/filepath"
	"time"

	"github.com/boltdb/bolt"
)

// buckets
var (
	downloadItemsBucket   = []byte("download-items")
	watchedTorrentsBucket = []byte("watched-torrents")
	defaultsBucket        = []byte("defaults")
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
		_, err = tx.CreateBucketIfNotExists(defaultsBucket)
		return err
	})
	if err != nil {
		return s.db.Close()
	}
	return nil
}

func (s *Store) Close() error { return s.db.Close() }

func (s *Store) Path() string { return s.path }

func (s *Store) CreateBuckets(forUser string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		userBkt, err := tx.CreateBucketIfNotExists([]byte(forUser))
		if err != nil {
			return err
		}

		buckets := [][]byte{
			downloadItemsBucket,
			watchedTorrentsBucket,
		}

		for _, bucket := range buckets {
			_, err = userBkt.CreateBucketIfNotExists(bucket)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// SaveState inserts or updates the given state.
func (s *Store) SaveState(state *State, forUser string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		userBkt := tx.Bucket([]byte(forUser))
		downloadsBkt := userBkt.Bucket(downloadItemsBucket)

		key := itob(state.FileID)
		var value bytes.Buffer

		err := gob.NewEncoder(&value).Encode(state)
		if err != nil {
			return err
		}

		return downloadsBkt.Put(key, value.Bytes())
	})
}

// State returns a state by the given file ID.
func (s *Store) State(id int64, forUser string) (*State, error) {
	var state State
	err := s.db.View(func(tx *bolt.Tx) error {
		userBkt := tx.Bucket([]byte(forUser))
		downloadsBkt := userBkt.Bucket(downloadItemsBucket)
		fileID := itob(id)

		value := downloadsBkt.Get(fileID)
		if value == nil {
			return ErrStateNotFound
		}

		return gob.NewDecoder(bytes.NewReader(value)).Decode(&state)
	})
	return &state, err
}

// States returns all the states in the store.
func (s *Store) States(forUser string) ([]*State, error) {
	states := make([]*State, 0)

	if forUser == "" {
		return states, nil
	}

	err := s.db.View(func(tx *bolt.Tx) error {
		userBkt := tx.Bucket([]byte(forUser))
		downloadsBkt := userBkt.Bucket(downloadItemsBucket)

		cursor := downloadsBkt.Cursor()
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

func (s *Store) Config(forUser string) (*Config, error) {
	if forUser == "" {
		return s.DefaultConfig()
	}

	var cfg *Config
	err := s.db.View(func(tx *bolt.Tx) error {
		userBkt := tx.Bucket([]byte(forUser))

		key := []byte("config")
		value := userBkt.Get(key)

		if value == nil {
			return ErrConfigNotFound
		}

		return gob.NewDecoder(bytes.NewReader(value)).Decode(cfg)
	})

	if err == ErrConfigNotFound {
		cfg, err = s.DefaultConfig()
	}

	return cfg, err
}

func (s *Store) SaveConfig(cfg *Config, forUser string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		userBkt := tx.Bucket([]byte(forUser))

		key := []byte("config")
		var value bytes.Buffer

		err := gob.NewEncoder(&value).Encode(cfg)
		if err != nil {
			return err
		}

		return userBkt.Put(key, value.Bytes())
	})
}

func (s *Store) DefaultConfig() (*Config, error) {
	u, err := user.Current()
	if err != nil {
		return nil, err
	}
	return &Config{
		PollInterval:        Duration(defaultPollInterval),
		DownloadTo:          filepath.Join(u.HomeDir, "putio-sync"),
		DownloadFrom:        defaultDownloadFrom,
		SegmentsPerFile:     defaultSegmentsPerFile,
		MaxParallelFiles:    defaultMaxParallelFiles,
		IsPaused:            true,
		WatchTorrentsFolder: false,
		TorrentsFolder:      "",
	}, nil
}

func (s *Store) CurrentUser() (string, error) {
	var username []byte
	err := s.db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(defaultsBucket)
		username = bkt.Get([]byte("current-user"))
		return nil
	})
	if err != nil {
		return "", err
	}

	if username == nil {
		return "", nil
	}
	return string(username), nil
}

func (s *Store) SaveCurrentUser(username string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(defaultsBucket)
		key := []byte("current-user")
		return bkt.Put(key, []byte(username))
	})
}

func itob(v int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

package main

import (
	"encoding/json"
	"time"

	"go.etcd.io/bbolt"
)

type Status string

const (
	StatusSynced      = "synced"
	StatusDownloading = "downloading"
	StatusUploading   = "uploading"
)

type State struct {
	Status           Status
	LocalInode       uint64
	RemoteID         int64
	DownloadTempPath string
	UploadURL        string
	Offset           int64
	// Snapshot represents the final state of the file.
	// It is set after the sync decision and direction is made.
	Snapshot *Snapshot
	relpath  string
}

type Snapshot struct {
	ModTime time.Time
	Size    int64
	CRC32   string
}

func ReadAllStates() ([]State, error) {
	var l []State
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketFiles)
		return b.ForEach(func(key, val []byte) error {
			var s State
			err := json.Unmarshal(val, &s)
			if err != nil {
				return err
			}
			s.relpath = string(key)
			l = append(l, s)
			return nil
		})
	})
	return l, err
}

func (s State) Write() error {
	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketFiles)
		val, err := json.Marshal(s)
		if err != nil {
			return err
		}
		return b.Put([]byte(s.relpath), val)
	})
}

func (s State) Delete() error {
	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketFiles)
		return b.Delete([]byte(s.relpath))
	})
}

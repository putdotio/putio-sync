package main

import (
	"encoding/json"

	"go.etcd.io/bbolt"
)

type Status string

const (
	StatusSynced      = "synced"
	StatusDownloading = "downloading"
	StatusUploading   = "uploading"
)

// State stores information about syncing files and folders.
type State struct {
	Status           Status
	IsDir            bool
	LocalInode       uint64
	RemoteID         int64
	DownloadTempPath string
	UploadURL        string
	Offset           int64
	Size             int64
	CRC32            string
	relpath          string
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

// Move writes the state to database while changing the relpath key.
// Move also deletes the record at old relpath.
func (s *State) Move(target string) error {
	err := db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketFiles)
		err := b.Delete([]byte(s.relpath))
		if err != nil {
			return err
		}
		val, err := json.Marshal(s)
		if err != nil {
			return err
		}
		return b.Put([]byte(target), val)
	})
	if err != nil {
		return err
	}
	s.relpath = target
	return nil
}

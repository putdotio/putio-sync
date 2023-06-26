package putiosync

import (
	"encoding/json"

	"go.etcd.io/bbolt"
)

type status string

const (
	statusSynced      = "synced"
	statusDownloading = "downloading"
	statusUploading   = "uploading"
)

var bucketFiles = []byte("files")

// State stores information about syncing files and folders.
type stateType struct {
	Status           status
	IsDir            bool
	LocalInode       uint64
	RemoteID         int64
	DownloadTempName string
	UploadURL        string
	Offset           int64
	Size             int64
	CRC32            string
	LocalRoot        string
	relpath          string
}

func readAllStates() ([]stateType, error) {
	var l []stateType
	err := db.Update(func(tx *bbolt.Tx) error {
		remove := make([][]byte, 0)
		b := tx.Bucket(bucketFiles)
		err := b.ForEach(func(key, val []byte) error {
			var s stateType
			err := json.Unmarshal(val, &s)
			if err != nil {
				return err
			}
			if s.LocalRoot != cfg.LocalDir {
				remove = append(remove, key)
			} else {
				s.relpath = string(key)
				l = append(l, s)
			}
			return nil
		})
		if err != nil {
			return err
		}
		for _, key := range remove {
			err = b.Delete(key)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return l, err
}

func (s stateType) Write() error {
	s.LocalRoot = cfg.LocalDir
	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketFiles)
		val, err := json.Marshal(s)
		if err != nil {
			return err
		}
		return b.Put([]byte(s.relpath), val)
	})
}

func (s stateType) Delete() error {
	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketFiles)
		return b.Delete([]byte(s.relpath))
	})
}

// Move writes the state to database while changing the relpath key.
// Move also deletes the record at old relpath.
func (s *stateType) Move(target string) error {
	s.LocalRoot = cfg.LocalDir
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

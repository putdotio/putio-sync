package sync

import (
	"bytes"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/cenk/bitfield"
	"github.com/igungor/go-putio/putio"
)

// XXX: Changing any of these constants will break the control file alignment.
const (
	// Control-file version. Each version changes the previous layout and may
	// be backwards incompatible.
	version = 1

	// Each bit in a bitfield represent this amount of bytes
	bitfieldPieceLength = 16 * 1024
)

// DownloadStatus represents the current status of a download.
type DownloadStatus int

const (
	DownloadIdle DownloadStatus = iota
	DownloadFailed
	DownloadInQueue
	DownloadPaused
	DownloadInProgress
	DownloadCompleted
)

// String implements fmt.Stringer interface for DownloadStatus.
func (ds DownloadStatus) String() string {
	var s string
	switch ds {
	case DownloadIdle:
		s = "idle"
	case DownloadFailed:
		s = "failed"
	case DownloadInQueue:
		s = "inqueue"
	case DownloadPaused:
		s = "paused"
	case DownloadInProgress:
		s = "inprogress"
	case DownloadCompleted:
		s = "completed"
	}
	return s
}

// MarshalJSON implements json.Marshaler interface for DownloadStatus.
func (ds DownloadStatus) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%v\"", ds)), nil
}

// State stores all the metadata and state of a download. It is encoded as Gob
// and stored to a persistent storage.
type State struct {
	// State version number
	Version uint `json:"version"`

	// File metadata
	FileID              int64  `json:"file_id"`
	FileName            string `json:"file_name"`
	FileLength          int64  `json:"file_length"`
	FileIcon            string `json:"file_icon"`
	FileType            string `json:"file_type"`
	CRC32               string `json:"crc32"`
	BitfieldPieceLength int    `json:"-"`

	// Absolute path of the stored file
	LocalPath string `json:"local_path"`

	// Directory of the file relative to the Put.io root folder
	RemoteDir string `json:"-"`

	// Download states
	DownloadStatus     DownloadStatus `json:"download_status"`
	DownloadStartedAt  time.Time      `json:"download_started_at"`
	DownloadFinishedAt time.Time      `json:"download_finished_at"`
	DownloadSpeed      float64        `json:"download_speed"`

	IsHidden bool `json:"-"`

	Error string `json:"fail-reason"`

	// mu guards below
	mu                              sync.Mutex
	BytesTransferredSinceLastUpdate int64     `json:"-"`
	Bitfield                        *Bitfield `json:"bitfield"`
}

func NewState(f putio.File, savedTo string) *State {
	bflength := uint32(f.Size / bitfieldPieceLength)
	excess := f.Size % bitfieldPieceLength
	if excess > 0 {
		bflength++
	}

	return &State{
		Version:             version,
		FileID:              int64(f.ID),
		FileLength:          f.Size,
		FileName:            f.Name,
		FileIcon:            f.Screenshot,
		FileType:            f.ContentType,
		CRC32:               f.CRC32,
		BitfieldPieceLength: bitfieldPieceLength,
		LocalPath:           filepath.Join(savedTo, f.Name),
		DownloadStatus:      DownloadIdle,
		Bitfield: &Bitfield{
			length:   bflength,
			Bitfield: bitfield.New(bflength),
		},
	}
}

// String implements fmt.Stringer interface for State.
func (s *State) String() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("ID: %v\n", s.FileID))
	buf.WriteString(fmt.Sprintf("Name: %v\n", s.FileName))
	buf.WriteString(fmt.Sprintf("Length: %v\n", s.FileLength))
	buf.WriteString(fmt.Sprintf("CRC32: %v\n", s.CRC32))
	buf.WriteString(fmt.Sprintf("Download Status: %v\n", s.DownloadStatus))
	switch s.DownloadStatus {
	case DownloadCompleted:
		buf.WriteString(fmt.Sprintf("Downloaded At: %v\n", s.DownloadFinishedAt))
	case DownloadFailed:
		buf.WriteString(fmt.Sprintf("*** Fail reason: %v\n", s.Error))
	}
	buf.WriteString(fmt.Sprintf("Bitfield Piece Length: %v\n", s.BitfieldPieceLength))
	buf.WriteString(fmt.Sprintf("Bitfield Summary: %v/%v\n", s.Bitfield.Count(), s.Bitfield.Len()))

	return buf.String()
}

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

type DownloadStatus int

const (
	DownloadIdle DownloadStatus = iota
	DownloadFailed
	DownloadInQueue
	DownloadPaused
	DownloadInProgress
	DownloadCompleted
)

func (ds DownloadStatus) String() string {
	var s string
	switch ds {
	case DownloadIdle:
		s = "idle"
	case DownloadFailed:
		s = "failed"
	case DownloadInQueue:
		s = "inqueue"
	case DownloadInProgress:
		s = "inprogress"
	case DownloadCompleted:
		s = "completed"
	}
	return s
}

func (ds DownloadStatus) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%v\"", ds)), nil
}

// State is the high level representation of a Put.io control file.
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

	// Download states
	DownloadStatus     DownloadStatus `json:"download_status"`
	DownloadStartedAt  time.Time      `json:"download_started_at"`
	DownloadFinishedAt time.Time      `json:"download_finished_at"`
	DownloadSpeed      float64        `json:"download_speed"`

	IsHidden bool `json:"-"`

	Error string `json:"fail-reason"`

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

// calculateChunks splits a filesize into count parts, in which every chunk is
// either smaller than blocksize, or can be fully divisible. Chunks has to be
// divisible by blocksize since we maintain a bitfield for each file, and each
// bit represents the said constant.
func calculateChunks(state *State, segmentnum uint) []*chunk {
	count := int64(segmentnum)
	pieceLength := int64(state.BitfieldPieceLength)

	// calculate the chunks of a resumable file.
	if state.Bitfield.Count() != 0 {
		var chunks []*chunk
		var idx uint32
		for {
			start, ok := state.Bitfield.FirstClear(idx)
			if !ok {
				break
			}

			end, ok := state.Bitfield.FirstSet(start)
			if !ok {
				chunks = append(chunks, &chunk{
					offset: int64(start) * pieceLength,
					length: state.FileLength - int64(start)*pieceLength,
				})
				break
			}

			chunks = append(chunks, &chunk{
				offset: int64(start) * pieceLength,
				length: int64(end-start) * pieceLength,
			})

			idx = end
		}
		return chunks
	}

	// calculate the chunks of a fresh new file.

	filesize := state.FileLength
	// don't even consider smaller files
	if filesize <= pieceLength || count <= 1 {
		return []*chunk{{offset: 0, length: filesize}}
	}

	// how many blocks fit perfectly on a filesize
	blockCount := filesize / pieceLength
	// how many bytes are left out
	excessBytes := filesize % pieceLength

	// If there are no blocks available for the given blocksize, we're gonna
	// reduce the count to the max available block count.
	if blockCount < count {
		count = blockCount
	}

	blocksPerUnit := blockCount / count
	excessBlocks := blockCount % count

	var chunks []*chunk
	for i := int64(0); i < count; i++ {
		chunks = append(chunks, &chunk{
			offset: i * blocksPerUnit * pieceLength,
			length: blocksPerUnit * pieceLength,
		})
	}

	if excessBlocks > 0 {
		offset := count * blocksPerUnit * pieceLength
		length := excessBlocks * pieceLength
		chunks = append(chunks, &chunk{
			offset: offset,
			length: length,
		})
	}

	// append excess bytes to the last chunk
	if excessBytes > 0 {
		c := chunks[len(chunks)-1]
		c.length += excessBytes
	}

	return chunks
}

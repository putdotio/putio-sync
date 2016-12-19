package sync

import (
	"encoding"
	"time"
)

// Default configuration values
const (
	// default configuration values
	defaultSegmentsPerFile  = 3
	defaultMaxParallelFiles = 2
	defaultDownloadFrom     = -1
	defaultPollInterval     = 2 * time.Minute

	limitSegmentsPerFile = 8
	limitParallelFiles   = 8
)

// Config represents the configuration of putio-sync application.
type Config struct {
	// Walk DownloadFrom directory for every n interval
	PollInterval Duration `json:"poll-interval"`

	// Download Put.io files to this directory
	DownloadTo string `json:"download-to"`

	// Download files only in this directory (Put.io file ID)
	DownloadFrom int64 `json:"download-from"`

	// Max number of connections to server for each download
	SegmentsPerFile uint `json:"segments-per-file"`

	// Max number of parallel file downloads
	MaxParallelFiles uint `json:"max-parallel-files"`

	// User's OAuth2 token for this application
	OAuth2Token string `json:"oauth2-token"`

	// Reports whether the folder should be watched
	WatchTorrentsFolder bool `json:"watch-torrents-folder"`

	// User's prefered folder to watch for new .torrent files
	TorrentsFolder string `json:"torrents-folder"`

	// Last pause/resume state
	IsPaused bool `json:"is-paused"`
}

// Duration is a JSON wrapper type for time.Duration.
type Duration time.Duration

// ensure duration implements both these interfaces for our json encoding/decoding.
var _ encoding.TextMarshaler = new(Duration)
var _ encoding.TextUnmarshaler = new(Duration)

func (d Duration) String() string {
	return time.Duration(d).String()
}

// MarshalText converts a duration to a string for decoding json.
func (d Duration) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

// UnmarshalText parses a JSON value into a Duration value.
func (d *Duration) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		return nil
	}

	duration, err := time.ParseDuration(string(text))
	if err != nil {
		return err
	}

	*d = Duration(duration)
	return nil
}

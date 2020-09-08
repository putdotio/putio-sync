package putiosync

import (
	"fmt"
	"os"

	"github.com/putdotio/go-putio"
	"github.com/putdotio/putio-sync/v2/internal/walker"
	"golang.org/x/text/unicode/norm"
)

type iLocalFile interface {
	Info() os.FileInfo
	RelPath() string
	FullPath() string
}

type iRemoteFile interface {
	Info() os.FileInfo
	RelPath() string
	PutioFile() *putio.File
}

type syncFile struct {
	local   iLocalFile
	remote  iRemoteFile
	state   *stateType
	relpath string
	skip    bool
}

func (f *syncFile) String() string {
	flags := []byte("...")
	if f.state != nil {
		switch f.state.Status {
		case statusSynced:
			flags[0] = 'S'
		case statusDownloading:
			flags[0] = 'D'
		case statusUploading:
			flags[0] = 'U'
		default:
			flags[0] = '?'
		}
	}
	if f.local != nil {
		flags[1] = 'L'
	}
	if f.remote != nil {
		flags[2] = 'R'
	}
	return fmt.Sprintf("%s %s", string(flags), f.relpath)
}

func groupFiles(states []stateType, localFiles []*walker.LocalFile, remoteFiles []*walker.RemoteFile) map[string]*syncFile {
	m := make(map[string]*syncFile)
	initSyncFile := func(relpath string) *syncFile {
		relpath = norm.NFC.String(relpath)
		sf, ok := m[relpath]
		if ok {
			return sf
		}
		sf = &syncFile{relpath: relpath}
		m[relpath] = sf
		return sf
	}
	for _, lf := range localFiles {
		sf := initSyncFile(lf.RelPath())
		sf.local = lf
	}
	for _, rf := range remoteFiles {
		sf := initSyncFile(rf.RelPath())
		sf.remote = rf
	}
	for _, state := range states {
		sf := initSyncFile(state.relpath)
		s := state
		sf.state = &s
	}
	return m
}

package main

import (
	"fmt"
)

type SyncFile struct {
	local   *LocalFile
	remote  *RemoteFile
	state   *State
	relpath string
	skip    bool
}

func (f *SyncFile) String() string {
	flags := []byte("...")
	if f.state != nil {
		switch f.state.Status {
		case StatusSynced:
			flags[0] = 'S'
		case StatusDownloading:
			flags[0] = 'D'
		case StatusUploading:
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

func GroupFiles(states []State, localFiles []*LocalFile, remoteFiles []*RemoteFile) map[string]*SyncFile {
	m := make(map[string]*SyncFile)
	initSyncFile := func(relpath string) *SyncFile {
		sf, ok := m[relpath]
		if ok {
			return sf
		}
		sf = &SyncFile{relpath: relpath}
		m[relpath] = sf
		return sf
	}
	for _, lf := range localFiles {
		sf := initSyncFile(lf.relpath)
		sf.local = lf
	}
	for _, rf := range remoteFiles {
		sf := initSyncFile(rf.relpath)
		sf.remote = rf
	}
	for _, state := range states {
		sf := initSyncFile(state.relpath)
		s := state
		sf.state = &s
	}
	return m
}

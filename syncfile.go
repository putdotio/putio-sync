package main

import (
	"fmt"
)

type SyncFile struct {
	local   *LocalFile
	remote  *RemoteFile
	state   *State
	relpath string
}

func (f *SyncFile) String() string {
	flags := []byte("....")
	if f.state != nil {
		if f.state.Snapshot != nil {
			flags[0] = 'S'
		}
		switch f.state.Status {
		case StatusDownloading:
			flags[3] = 'D'
		case StatusUploading:
			flags[3] = 'U'
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

func GroupFiles(states []State, localFiles, remoteFiles []File) map[string]*SyncFile {
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
	for _, f := range localFiles {
		lf := f.(*LocalFile)
		sf := initSyncFile(lf.relpath)
		sf.local = lf
	}
	for _, f := range remoteFiles {
		rf := f.(*RemoteFile)
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

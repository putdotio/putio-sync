package main

import (
	"sort"

	"github.com/cenkalti/log"
)

// Reconciliation funciton does not perform any operation.
// It only takes the filesystem state as input and returns operations to be performed on those fileystems.
// It must be testable without side effects.
func Reconciliation(syncFiles map[string]*SyncFile) ([]Job, error) {
	var jobs []Job
	files := SortSyncFiles(syncFiles)
	for _, sf := range files {
		log.Debugln(sf.String())
	}
	for _, sf := range files {
		// TODO create jobs for non-existing folders
		// TODO detect deletes
		// TODO detect moves
		switch {
		case sf.local != nil && sf.remote == nil:
			if sf.local.info.IsDir() {
				jobs = append(jobs, &CreateRemoteFolder{
					relpath: sf.relpath,
					state:   sf.state,
				})
			} else {
				jobs = append(jobs, &Upload{
					localFile: sf.local,
					state:     sf.state,
				})
			}
		case sf.local == nil && sf.remote != nil:
			if sf.remote.putioFile.IsDir() {
				jobs = append(jobs, &CreateLocalFolder{
					relpath: sf.relpath,
					state:   sf.state,
				})
			} else {
				jobs = append(jobs, &Download{
					remoteFile: sf.remote,
					state:      sf.state,
				})
			}
		case sf.local != nil && sf.remote != nil:
			// Exist on both sides
			// TODO handle conflicts for files exist on both sides
		case sf.state != nil:
			jobs = append(jobs, &DeleteState{
				state: *sf.state,
			})
		}
	}
	return jobs, nil
}

func SortSyncFiles(m map[string]*SyncFile) []*SyncFile {
	var a []*SyncFile
	for _, sf := range m {
		a = append(a, sf)
	}
	sort.Slice(a, func(i, j int) bool {
		return a[i].relpath < a[j].relpath
	})
	return a
}

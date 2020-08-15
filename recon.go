package main

import (
	"sort"
	"strings"

	"github.com/cenkalti/log"
)

// Reconciliation function does not perform any operation.
// It only takes the filesystem state as input and returns operations to be performed on those fileystems.
// It must be testable without side effects.
func Reconciliation(syncFiles map[string]*SyncFile) []Job {
	var jobs []Job

	// Sort files by path for deterministic output
	files := sortSyncFiles(syncFiles)
	for _, sf := range files {
		log.Debugln(sf.String())
	}

	// Group files by certain keys for faster lookup
	filesByRemoteID := mapRemoteFilesByID(syncFiles)
	filesByInode := mapLocalFilesByInode(syncFiles)

	// First, sync files with known state
	// This is required for detecting simple move operations correctly.
	for _, sf := range files {
		if sf.state != nil {
			for _, job := range syncWithState(sf, filesByRemoteID, filesByInode) {
				if job != nil {
					jobs = append(jobs, job)
				}
			}
		}
	}

	// Then, sync first seen files
	for _, sf := range files {
		if sf.state == nil && !sf.skip {
			job := syncFresh(sf)
			if job != nil {
				jobs = append(jobs, job)
			}
		}
	}

	return jobs
}

func syncFresh(sf *SyncFile) Job {
	switch {
	case sf.local != nil && sf.remote == nil:
		// File present only on local side. Copy to the remote side.
		if sf.local.info.IsDir() {
			return &CreateRemoteFolder{
				relpath: sf.relpath,
			}
		}
		return &Upload{
			localFile: sf.local,
		}
	case sf.local == nil && sf.remote != nil:
		// File present only on remote side. Copy to the local side.
		if sf.remote.putioFile.IsDir() {
			return &CreateLocalFolder{
				relpath:  sf.relpath,
				remoteID: sf.remote.putioFile.ID,
			}
		}
		return &Download{
			remoteFile: sf.remote,
		}
	case sf.local != nil && sf.remote != nil:
		// File exists on both sides
		switch {
		case sf.local.info.IsDir() && sf.remote.Info().IsDir():
			// Dir exists on both sides, save this state to db
			return &WriteDirState{
				remoteID: sf.remote.putioFile.ID,
				relpath:  sf.remote.relpath,
			}
		case sf.local.info.IsDir() || sf.remote.Info().IsDir():
			// One of the sides is a dir, the other is a file
			log.Warningf("Conflicting file, skipping sync: %q", sf.relpath)
			return nil
		// Both sides are file, not folder
		case sf.local.info.Size() != sf.remote.putioFile.Size:
			log.Warningf("File sizes differ, skipping sync: %q", sf.relpath)
			return nil
		default:
			// Assume files are same if they are in same size
			// TODO check crc32 for local file
			return &WriteFileState{
				localFile:  *sf.local,
				remoteFile: *sf.remote,
			}
		}
	default:
		log.Errorf("Unhandled case for %q", sf.relpath)
		return nil
	}
}

func syncWithState(sf *SyncFile, filesByRemoteID map[int64]*SyncFile, filesByInode map[uint64]*SyncFile) []Job {
	// We have a state from previous sync. Compare local and remote sides with existing state.
	switch sf.state.Status {
	case StatusSynced:
		switch {
		case sf.local != nil && sf.remote != nil:
			// Exist on both sides
			if sf.local.info.IsDir() && sf.remote.Info().IsDir() {
				// Both sides are directory
				return nil
			}
			if sf.local.info.IsDir() || sf.remote.Info().IsDir() {
				// One of the sides is a file
				log.Warningf("Conflicting file, skipping sync: %q", sf.relpath)
				return nil
			}
			if sf.state.Size != sf.local.info.Size() || sf.state.Size != sf.remote.putioFile.Size {
				log.Warningf("File sizes differ, skipping sync: %q", sf.relpath)
				return nil
			}
			// Assume files didn't change if their size didn't change
			// This is the hottest case that is executed most because once all files are in sync no operations will be done later.
			return nil
		case sf.local != nil && sf.remote == nil:
			// File missing in remote side, could be deleted or moved elsewhere
			target, ok := filesByRemoteID[sf.state.RemoteID]
			if ok { // nolint: nestif
				// File with the same id is found on another path
				if target.state == nil {
					// There is no existing state in move target
					if target.remote.putioFile.CRC32 == sf.state.CRC32 {
						// Remote file is not changed
						inode, _ := GetInode(sf.local.info)
						if inode == sf.state.LocalInode {
							// Local file is not changed too
							// Then, file must be moved. We can move the local file to same path.
							target.skip = true
							return []Job{&MoveLocalFile{
								localFile: *sf.local,
								toRelpath: target.relpath,
								state:     *sf.state,
							}}
						}
					}
				}
			}
			// File is deleted on remote side
			return []Job{&DeleteLocalFile{
				localFile: *sf.local,
				state:     *sf.state,
			}}
		case sf.local == nil && sf.remote != nil:
			// File missing in local side, could be deleted or moved elsewhere
			target, ok := filesByInode[sf.state.LocalInode]
			if ok { // nolint: nestif
				// File with same inode is found on another path
				if target.state == nil {
					if sf.remote.putioFile.CRC32 == sf.state.CRC32 {
						// Remote file is not changed
						inode, _ := GetInode(target.local.info)
						if inode == sf.state.LocalInode {
							// Local file is not changed too
							// Then, file must be moved. We can move the remote file to same path.
							target.skip = true
							return []Job{&MoveRemoteFile{
								remoteFile: *sf.remote,
								toRelpath:  target.relpath,
								state:      *sf.state,
							}}
						}
					}
				}
			}
			// File is delete on local side
			return []Job{&DeleteRemoteFile{
				remoteFile: *sf.remote,
				state:      *sf.state,
			}}
		case sf.local == nil && sf.remote == nil:
			// File deleted on both sides, handled by other cases above, let's delete the state.
			return []Job{&DeleteState{
				state: *sf.state,
			}}
		default:
			log.Errorf("Unhandled case for %q", sf.relpath)
			return nil
		}
	case StatusDownloading:
		if sf.local == nil && sf.remote != nil {
			if sf.remote.putioFile.CRC32 == sf.state.CRC32 {
				// Remote file is still the same, resume download
				return []Job{&Download{
					remoteFile: sf.remote,
					state:      sf.state,
				}}
			}
		}
		// Cancel current download and make a new sync decision
		return []Job{
			&DeleteState{
				state: *sf.state,
			},
			syncFresh(sf),
		}
	case StatusUploading:
		if sf.local != nil && sf.remote == nil {
			inode, _ := GetInode(sf.local.info)
			if inode == sf.state.LocalInode {
				// Local file is still the same, resume upload
				return []Job{&Upload{
					localFile: sf.local,
					state:     sf.state,
				}}
			}
		}
		// Cancel current upload and make a new sync decision
		return []Job{
			&DeleteState{
				state: *sf.state,
			},
			syncFresh(sf),
		}
	default:
		// Invalid status, should not happen in normal usage.
		return []Job{&DeleteState{
			state: *sf.state,
		}}
	}
}

// sortSyncFiles so that folders come after regular files.
// This is required for correct moving of folders with files inside.
// First, files are synced, then folder can be moved or deleted.
func sortSyncFiles(m map[string]*SyncFile) []*SyncFile {
	a := make([]*SyncFile, 0, len(m))
	for _, sf := range m {
		a = append(a, sf)
	}
	sort.Slice(a, func(i, j int) bool {
		x, y := a[i], a[j]
		if isChildOf(x.relpath, y.relpath) {
			return true
		}
		if isChildOf(y.relpath, x.relpath) {
			return false
		}
		return x.relpath < y.relpath
	})
	return a
}

func isChildOf(child, parent string) bool {
	return strings.HasPrefix(child, parent+"/")
}

func mapRemoteFilesByID(syncFiles map[string]*SyncFile) map[int64]*SyncFile {
	m := make(map[int64]*SyncFile, len(syncFiles))
	for _, sf := range syncFiles {
		if sf.remote != nil {
			m[sf.remote.putioFile.ID] = sf
		}
	}
	return m
}

func mapLocalFilesByInode(syncFiles map[string]*SyncFile) map[uint64]*SyncFile {
	m := make(map[uint64]*SyncFile, len(syncFiles))
	for _, sf := range syncFiles {
		if sf.local != nil {
			inode, err := GetInode(sf.local.info)
			if err != nil {
				log.Error(err)
				continue
			}
			m[inode] = sf
		}
	}
	return m
}

package putiosync

import (
	"sort"
	"strings"

	"github.com/cenkalti/log"
	"github.com/putdotio/putio-sync/internal/inode"
)

// Reconciliation function does not perform any operation.
// It only takes the filesystem state as input and returns operations to be performed on those fileystems.
// It must be testable without side effects.
func reconciliation(syncFiles map[string]*syncFile) []iJob {
	var jobs []iJob

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

func syncFresh(sf *syncFile) iJob {
	switch {
	case sf.local != nil && sf.remote == nil:
		// File present only on local side. Copy to the remote side.
		if sf.local.Info().IsDir() {
			return &createRemoteFolderJob{
				relpath: sf.relpath,
			}
		}
		return &uploadJob{
			localFile: sf.local,
		}
	case sf.local == nil && sf.remote != nil:
		// File present only on remote side. Copy to the local side.
		if sf.remote.PutioFile().IsDir() {
			return &createLocalFolderJob{
				relpath:  sf.relpath,
				remoteID: sf.remote.PutioFile().ID,
			}
		}
		return &downloadJob{
			remoteFile: sf.remote,
		}
	case sf.local != nil && sf.remote != nil:
		// File exists on both sides
		switch {
		case sf.local.Info().IsDir() && sf.remote.Info().IsDir():
			// Dir exists on both sides, save this state to db
			return &writeDirStateJob{
				remoteID: sf.remote.PutioFile().ID,
				relpath:  sf.remote.RelPath(),
			}
		case sf.local.Info().IsDir() || sf.remote.Info().IsDir():
			// One of the sides is a dir, the other is a file
			log.Warningf("Conflicting file, skipping sync: %q", sf.relpath)
			return nil
		// Both sides are file, not folder
		case sf.local.Info().Size() != sf.remote.PutioFile().Size:
			log.Warningf("File sizes differ, skipping sync: %q", sf.relpath)
			return nil
		default:
			// Assume files are same if they are in same size
			return &writeFileStateJob{
				localFile:  sf.local,
				remoteFile: sf.remote,
			}
		}
	default:
		log.Errorf("Unhandled case for %q", sf.relpath)
		return nil
	}
}

func syncWithState(sf *syncFile, filesByRemoteID map[int64]*syncFile, filesByInode map[uint64]*syncFile) []iJob {
	// We have a state from previous sync. Compare local and remote sides with existing state.
	switch sf.state.Status {
	case statusSynced:
		switch {
		case sf.local != nil && sf.remote != nil:
			// Exist on both sides
			if sf.local.Info().IsDir() && sf.remote.Info().IsDir() {
				// Both sides are directory
				return nil
			}
			if sf.local.Info().IsDir() || sf.remote.Info().IsDir() {
				// One of the sides is a file
				log.Warningf("Conflicting file, skipping sync: %q", sf.relpath)
				return nil
			}
			if sf.state.Size != sf.local.Info().Size() || sf.state.Size != sf.remote.PutioFile().Size {
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
					if target.remote.PutioFile().CRC32 == sf.state.CRC32 {
						// Remote file is not changed
						in, _ := inode.Get(sf.local.Info())
						if in == sf.state.LocalInode {
							// Local file is not changed too
							// Then, file must be moved. We can move the local file to same path.
							target.skip = true
							return []iJob{&moveLocalFileJob{
								localFile: sf.local,
								toRelpath: target.relpath,
								state:     *sf.state,
							}}
						}
					}
				}
			}
			// File is deleted on remote side
			return []iJob{&deleteLocalFileJob{
				localFile: sf.local,
				state:     *sf.state,
			}}
		case sf.local == nil && sf.remote != nil:
			// File missing in local side, could be deleted or moved elsewhere
			target, ok := filesByInode[sf.state.LocalInode]
			if ok { // nolint: nestif
				// File with same inode is found on another path
				if target.state == nil {
					if sf.remote.PutioFile().CRC32 == sf.state.CRC32 {
						// Remote file is not changed
						in, _ := inode.Get(target.local.Info())
						if in == sf.state.LocalInode {
							// Local file is not changed too
							// Then, file must be moved. We can move the remote file to same path.
							target.skip = true
							return []iJob{&moveRemoteFileJob{
								remoteFile: sf.remote,
								toRelpath:  target.relpath,
								state:      *sf.state,
							}}
						}
					}
				}
			}
			// File is delete on local side
			return []iJob{&deleteRemoteFileJob{
				remoteFile: sf.remote,
				state:      *sf.state,
			}}
		case sf.local == nil && sf.remote == nil:
			// File deleted on both sides, handled by other cases above, let's delete the state.
			return []iJob{&deleteStateJob{
				state: *sf.state,
			}}
		default:
			log.Errorf("Unhandled case for %q", sf.relpath)
			return nil
		}
	case statusDownloading:
		if sf.local == nil && sf.remote != nil {
			if sf.remote.PutioFile().CRC32 == sf.state.CRC32 {
				// Remote file is still the same, resume download
				return []iJob{&downloadJob{
					remoteFile: sf.remote,
					state:      sf.state,
				}}
			}
		}
		// Cancel current download and make a new sync decision
		return []iJob{
			&deleteStateJob{
				state: *sf.state,
			},
			syncFresh(sf),
		}
	case statusUploading:
		if sf.local != nil && sf.remote == nil {
			in, _ := inode.Get(sf.local.Info())
			if in == sf.state.LocalInode {
				// Local file is still the same, resume upload
				return []iJob{&uploadJob{
					localFile: sf.local,
					state:     sf.state,
				}}
			}
		}
		// Cancel current upload and make a new sync decision
		return []iJob{
			&deleteStateJob{
				state: *sf.state,
			},
			syncFresh(sf),
		}
	default:
		// Invalid status, should not happen in normal usage.
		return []iJob{&deleteStateJob{
			state: *sf.state,
		}}
	}
}

// sortSyncFiles so that folders come after regular files.
// This is required for correct moving of folders with files inside.
// First, files are synced, then folder can be moved or deleted.
func sortSyncFiles(m map[string]*syncFile) []*syncFile {
	a := make([]*syncFile, 0, len(m))
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

func mapRemoteFilesByID(syncFiles map[string]*syncFile) map[int64]*syncFile {
	m := make(map[int64]*syncFile, len(syncFiles))
	for _, sf := range syncFiles {
		if sf.remote != nil {
			m[sf.remote.PutioFile().ID] = sf
		}
	}
	return m
}

func mapLocalFilesByInode(syncFiles map[string]*syncFile) map[uint64]*syncFile {
	m := make(map[uint64]*syncFile, len(syncFiles))
	for _, sf := range syncFiles {
		if sf.local != nil {
			in, err := inode.Get(sf.local.Info())
			if err != nil {
				log.Error(err)
				continue
			}
			m[in] = sf
		}
	}
	return m
}

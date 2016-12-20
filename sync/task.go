package sync

import (
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"io"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

// Task represent a download task, which is closely associated with a Put.io
// file.
type Task struct {
	state  *State
	cwd    string
	chunks []*chunk
}

// NewTask creates a new Task, with a fresh internal state.
func NewTask(state *State, cwd string, segmentnum uint) *Task {
	chunks := calculateChunks(state, segmentnum)
	return &Task{
		state:  state,
		cwd:    cwd,
		chunks: chunks,
	}
}

// String implements fmt.Stringer interface for the Task.
func (t Task) String() string {
	return fmt.Sprintf("task<name: %q, size: %v, chunks: %v, bitfield: %v>",
		trimPath(path.Join(t.cwd, t.state.FileName)),
		t.state.FileLength,
		t.chunks,
		t.state.Bitfield.Len(),
	)
}

// Verify checks bitfield integrity and computes CRC32 of the task.
func (t *Task) Verify(r io.Reader) error {
	if !t.state.Bitfield.All() {
		return fmt.Errorf("Not all bits are downloaded for task: %q\n", t)
	}

	h := crc32.NewIEEE()
	_, err := io.Copy(h, r)
	if err != nil {
		return err
	}

	sum := h.Sum(nil)
	sumHex := hex.EncodeToString(sum)
	if sumHex != t.state.CRC32 {
		return fmt.Errorf("CRC32 check failed. got: %x want: %v", sumHex, t.state.CRC32)
	}

	return nil
}

// Tasks stores active tasks.
type Tasks struct {
	sync.Mutex
	s map[int64]*Task
}

// NewTasks allocated a new Tasks, which stores active tasks.
func NewTasks() *Tasks {
	return &Tasks{s: make(map[int64]*Task)}
}

// Add adds given task to the store.
func (m *Tasks) Add(t *Task) {
	m.Lock()
	defer m.Unlock()

	m.s[t.state.FileID] = t
}

// Remove removes task from the store.
func (m *Tasks) Remove(t *Task) {
	m.Lock()
	defer m.Unlock()

	delete(m.s, t.state.FileID)
}

// Exists reports whether given task is in the store.
func (m *Tasks) Exists(t *Task) bool {
	m.Lock()
	defer m.Unlock()

	_, ok := m.s[t.state.FileID]
	return ok
}

// Empty reports whether there are active tasks.
func (m *Tasks) Empty() bool {
	m.Lock()
	defer m.Unlock()

	return len(m.s) == 0
}

// trimPath trims the given path.
// E.g. /usr/local/bin/foo becomes /u/l/b/foo.
func trimPath(p string) string {
	if len(p) < 60 {
		return p
	}
	p = filepath.Clean(p)
	parts := strings.Split(p, string(filepath.Separator))
	for i, part := range parts {
		if part == "" {
			parts[i] = string(filepath.Separator)
			continue
		}
		// skip last element
		if i == len(parts)-1 {
			continue
		}
		parts[i] = string(part[0])
	}

	return filepath.Join(parts...)
}

package putiosync

import (
	"os"
	"testing"
	"time"

	"github.com/putdotio/go-putio"
)

func TestReconciliation(t *testing.T) {
	m := map[string]*syncFile{
		"foo": {
			relpath: "foo",
			local:   fakeLocalFile(t, "foo"),
		},
		"bar": {
			relpath: "bar",
			remote:  fakeRemoteFile("bar"),
		},
	}
	jobs := reconciliation(m)
	if len(jobs) != 2 {
		t.FailNow()
	}
	var j iJob
	var ok bool
	j = jobs[0]
	_, ok = j.(*downloadJob)
	if !ok {
		t.Fatal("job is not download")
	}
	j = jobs[1]
	_, ok = j.(*uploadJob)
	if !ok {
		t.Fatal("job is not upload")
	}
}

type FakeLocalFile struct {
	info     os.FileInfo
	relpath  string
	fullpath string
}

func fakeLocalFile(t *testing.T, relpath string) *FakeLocalFile {
	f, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		t.Fatal(err)
	}
	err = f.Close()
	if err != nil {
		t.Fatal(err)
	}
	return &FakeLocalFile{
		info:     fi,
		relpath:  relpath,
		fullpath: f.Name(),
	}
}

func (f *FakeLocalFile) Info() os.FileInfo { return f.info }
func (f *FakeLocalFile) RelPath() string   { return f.relpath }
func (f *FakeLocalFile) FullPath() string  { return f.fullpath }

type FakeRemoteFile struct {
	putioFile putio.File
	relpath   string
}

func fakeRemoteFile(relpath string) *FakeRemoteFile {
	return &FakeRemoteFile{
		relpath:   relpath,
		putioFile: putio.File{},
	}
}

func (f *FakeRemoteFile) Info() os.FileInfo      { return &PutioFileInfo{File: f.putioFile} }
func (f *FakeRemoteFile) RelPath() string        { return f.relpath }
func (f *FakeRemoteFile) PutioFile() *putio.File { return &f.putioFile }

type PutioFileInfo struct {
	putio.File
}

var _ os.FileInfo = (*PutioFileInfo)(nil)

func (fi *PutioFileInfo) Name() string       { return fi.File.Name }
func (fi *PutioFileInfo) Size() int64        { return fi.File.Size }
func (fi *PutioFileInfo) ModTime() time.Time { return fi.File.CreatedAt.Time }
func (fi *PutioFileInfo) Sys() interface{}   { return nil }
func (fi *PutioFileInfo) IsDir() bool        { return fi.File.IsDir() }
func (fi *PutioFileInfo) Mode() os.FileMode {
	if fi.IsDir() {
		return 755 | os.ModeDir
	}
	return 644
}

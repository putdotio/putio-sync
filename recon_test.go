package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/putdotio/go-putio"
)

func TestReconciliation(t *testing.T) {
	m := map[string]*SyncFile{
		"foo": {
			relpath: "foo",
			local:   fakeLocalFile(t, "foo"),
		},
		"bar": {
			relpath: "bar",
			remote:  fakeRemoteFile("bar"),
		},
	}
	jobs := Reconciliation(m)
	if len(jobs) != 2 {
		t.FailNow()
	}
	var j Job
	var ok bool
	j = jobs[0]
	_, ok = j.(*Download)
	if !ok {
		t.Fatal("job is not download")
	}
	j = jobs[1]
	_, ok = j.(*Upload)
	if !ok {
		t.Fatal("job is not upload")
	}
}

func fakeLocalFile(t *testing.T, relpath string) (lf *LocalFile) {
	f, err := ioutil.TempFile("", "")
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
	return &LocalFile{
		info:    fi,
		relpath: relpath,
	}
}

func fakeRemoteFile(relpath string) (rf *RemoteFile) {
	return &RemoteFile{
		relpath:   relpath,
		putioFile: putio.File{},
	}
}

package putio

import (
	"fmt"
	"net/http"
	"testing"
)

func TestZips_Get(t *testing.T) {
	setup()
	defer teardown()

	fixture := `
{
	"missing_files": [],
	"size": 27039611973,
	"status": "OK",
	"url": "https://some-valid-storage-url.com/12345"
}
`
	mux.HandleFunc("/v2/zips/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintln(w, fixture)
	})

	zip, err := client.Zips.Get(nil, 1)
	if err != nil {
		t.Error(err)
	}

	if zip.URL != "https://some-valid-storage-url.com/12345" {
		t.Errorf("got: %v, want: https://some-valid-storage-url.com/12345", zip.URL)
	}

	// negative id
	_, err = client.Zips.Get(nil, -1)
	if err == nil {
		t.Errorf("negative id accepted")
	}
}

func TestZips_List(t *testing.T) {
	setup()
	defer teardown()

	fixture := `
{
	"status": "OK",
	"zips": [
		{
			"created_at": "2016-07-15T10:42:12",
			"id": 4177262
		}
	]
}
`
	mux.HandleFunc("/v2/zips/list", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintln(w, fixture)
	})

	zips, err := client.Zips.List(nil)
	if err != nil {
		t.Error(err)
	}

	if len(zips) != 1 {
		t.Errorf("got: %v, want: 1", len(zips))
	}

	if zips[0].ID != 4177262 {
		t.Errorf("got: %v, want: 4177262", zips[0].ID)
	}
}

func TestZips_Create(t *testing.T) {
	setup()
	defer teardown()

	fixture := `
{
	"status": "OK",
	"zip_id": 4177264
}
`
	mux.HandleFunc("/v2/zips/create", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Content-Type", "application/x-www-form-urlencoded")
		fmt.Fprintln(w, fixture)
	})

	id, err := client.Zips.Create(nil, 666)
	if err != nil {
		t.Error(err)
	}

	if id != 4177264 {
		t.Errorf("got: %v, want 4177264", id)
	}

	// negative id
	_, err = client.Zips.Create(nil, 1, 2, -1)
	if err == nil {
		t.Errorf("negative id accepted")
	}

	_, err = client.Zips.Create(nil)
	if err == nil {
		t.Errorf("empty params accepted")
	}
}

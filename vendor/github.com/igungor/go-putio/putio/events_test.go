package putio

import (
	"fmt"
	"net/http"
	"testing"
)

func TestEvents_List(t *testing.T) {
	setup()
	defer teardown()

	fixture := `
{
  "events": [
    {
      "created_at": "2016-06-16 11:12:57",
      "file_id": 409621890,
      "id": 26494492,
      "transfer_name": "Ubuntu 16.04.LTS.iso",
      "transfer_size": 334014003,
      "type": "transfer_completed"
    },
    {
      "created_at": "2016-06-16 05:44:24",
      "file_id": 409583253,
      "id": 26490777,
      "transfer_name": "Big.Buck.Bunny.mkv",
      "transfer_size": 613389392,
      "type": "transfer_completed"
    }
  ],
  "status": "OK"
}
`
	mux.HandleFunc("/v2/events/list", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintln(w, fixture)
	})

	events, err := client.Events.List(nil)
	if err != nil {
		t.Error(err)
	}

	if len(events) != 2 {
		t.Errorf("got: %v, want 2", len(events))
	}

	if events[0].ID != 26494492 {
		t.Errorf("got: %v, want: 26494492", events[0].ID)
	}
}

func TestEvents_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/events/delete", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Content-Type", "application/x-www-form-urlencoded")
		fmt.Fprintln(w, `{"status":"OK"}`)
	})

	err := client.Events.Delete(nil)
	if err != nil {
		t.Error(err)
	}
}

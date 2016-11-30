package putio

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestTransfers_Get(t *testing.T) {
	setup()
	defer teardown()

	fixture := `
	{
	"status": "OK",
	"transfer": {
		"availability": null,
		"callback_url": null,
		"client_ip": null,
		"created_at": "2016-07-15T09:45:15",
		"created_torrent": false,
		"current_ratio": 0.00,
		"down_speed": 0,
		"download_id": 20448117,
		"downloaded": 0,
		"error_message": null,
		"estimated_time": null,
		"extract": false,
		"file_id": 415107363,
		"finished_at": "2016-07-15T09:45:16",
		"id": 1,
		"is_private": false,
		"magneturi": "magnet:?xt=urn:btih:4344503b7e797ebf31582327a5baae35b11bda01&dn=ubuntu-16.04-desktop-amd64.iso",
		"name": "ubuntu-16.04-desktop-amd64.iso",
		"peers_connected": 0,
		"peers_getting_from_us": 0,
		"peers_sending_to_us": 0,
		"percent_done": 100,
		"save_parent_id": 400004368,
		"seconds_seeding": 0,
		"simulated": true,
		"size": 1485881344,
		"source": "http://releases.ubuntu.com/16.04/ubuntu-16.04-desktop-amd64.iso.torrent",
		"status": "COMPLETED",
		"status_message": "Completed 4 mins ago.",
		"subscription_id": null,
		"torrent_link": "/v2/transfers/36003178/torrent",
		"tracker_message": null,
		"trackers": null,
		"type": "TORRENT",
		"up_speed": 0,
		"uploaded": 0
	}
	}
	`
	mux.HandleFunc("/v2/transfers/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintln(w, fixture)
	})

	transfer, err := client.Transfers.Get(nil, 1)
	if err != nil {
		t.Error(err)
	}

	if transfer.ID != 1 {
		t.Errorf("got: %v, want: 1", transfer.ID)
	}

	// negative id
	_, err = client.Transfers.Get(nil, -1)
	if err == nil {
		t.Errorf("negative id accepted")
	}
}

func TestTransfers_List(t *testing.T) {
	setup()
	defer teardown()

	fixture := `
{
"status": "OK",
"transfers": [
	{
	"availability": null,
	"callback_url": null,
	"client_ip": null,
	"created_at": "2016-07-15T09:45:15",
	"created_torrent": false,
	"current_ratio": 0.00,
	"down_speed": 0,
	"download_id": 20448117,
	"downloaded": 0,
	"error_message": null,
	"estimated_time": null,
	"extract": false,
	"file_id": 415107363,
	"finished_at": "2016-07-15T09:45:16",
	"id": 36003178,
	"is_private": false,
	"links": [
		{
			"label": "ud_logs",
			"url": "https://papertrailapp.com/events?q=user_download_36003178"
		},
		{
			"label": "d_logs",
			"url": "https://papertrailapp.com/events?q=download_id_20448117"
		},
		{
			"label": "torrent",
			"url": "/v2/transfers/36003178/torrent"
		},
		{
			"label": "magnet",
			"url": "magnet:?xt=urn:btih:4344503b7e797ebf31582327a5baae35b11bda01&dn=ubuntu-16.04-desktop-amd64.iso"
		}
	],
	"magneturi": "magnet:?xt=urn:btih:4344503b7e797ebf31582327a5baae35b11bda01&dn=ubuntu-16.04-desktop-amd64.iso",
	"name": "ubuntu-16.04-desktop-amd64.iso",
	"peers_connected": 0,
	"peers_getting_from_us": 0,
	"peers_sending_to_us": 0,
	"percent_done": 100,
	"save_parent_id": 400004368,
	"seconds_seeding": 0,
	"simulated": true,
	"size": 1485881344,
	"source": "http://releases.ubuntu.com/16.04/ubuntu-16.04-desktop-amd64.iso.torrent",
	"status": "COMPLETED",
	"status_message": "Completed 29 mins ago.",
	"subscription_id": null,
	"torrent_link": "/v2/transfers/36003178/torrent",
	"tracker_message": null,
	"trackers": null,
	"type": "TORRENT",
	"up_speed": 0,
	"uploaded": 0
}]
}
	`

	mux.HandleFunc("/v2/transfers/list", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintln(w, fixture)
	})

	transfers, err := client.Transfers.List(nil)
	if err != nil {
		t.Error(err)
	}

	if len(transfers) != 1 {
		t.Errorf("got: %v, want: 1", len(transfers))
	}

	if transfers[0].ID != 36003178 {
		t.Errorf("got: %v, want: 36003178", transfers[0].ID)
	}
}

func TestTransfers_Add(t *testing.T) {
	setup()
	defer teardown()

	fixture := `
{
	"status": "OK",
	"transfer": {
		"availability": null,
		"callback_url": null,
		"client_ip": null,
		"created_at": "2016-07-15T09:45:15",
		"created_torrent": false,
		"current_ratio": 0.00,
		"down_speed": 0,
		"download_id": null,
		"downloaded": 0,
		"error_message": null,
		"estimated_time": null,
		"extract": false,
		"file_id": null,
		"finished_at": null,
		"id": 3600317,
		"is_private": false,
		"magneturi": "magnet:?xt=urn:btih:4344503b7e797ebf31582327a5baae35b11bda01&dn=ubuntu-16.04-desktop-amd64.iso",
		"name": "ubuntu-16.04-desktop-amd64.iso",
		"peers_connected": 0,
		"peers_getting_from_us": 0,
		"peers_sending_to_us": 0,
		"percent_done": 0,
		"save_parent_id": 400004368,
		"seconds_seeding": 0,
		"simulated": false,
		"size": 1485881344,
		"source": "http://releases.ubuntu.com/16.04/ubuntu-16.04-desktop-amd64.iso.torrent",
		"status": "IN_QUEUE",
		"status_message": "In queue...",
		"subscription_id": null,
		"torrent_link": "/v2/transfers/36003178/torrent",
		"tracker_message": null,
		"trackers": null,
		"type": "TORRENT",
		"up_speed": 0,
		"uploaded": 0
	}
}
`
	mux.HandleFunc("/v2/transfers/add", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Content-Type", "application/x-www-form-urlencoded")

		// form values
		u := r.FormValue("url")
		if u == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		fmt.Fprintln(w, fixture)
	})

	transfer, err := client.Transfers.Add(nil, "http://releases.ubuntu.com/16.04/ubuntu-16.04-desktop-amd64.iso.torrent", 0, "")
	if err != nil {
		t.Error(err)
	}

	if transfer.ID != 3600317 {
		t.Errorf("got: %v, want: 3100317", transfer.ID)
	}

	// empty URL
	_, err = client.Transfers.Add(nil, "", 0, "")
	if err == nil {
		t.Errorf("empty URL accepted")
	}

	// negative parent folder means use the user's prefered download folder.
	_, err = client.Transfers.Add(nil, "filepath", -1, "")
	if err != nil {
		t.Error(err)
	}

	// callback-url
	_, err = client.Transfers.Add(nil, "filepath", -1, "https://some-valid-endpoint-for-post-hook.com")
	if err != nil {
		t.Error(err)
	}
}

func TestTransfers_Retry(t *testing.T) {
	setup()
	defer teardown()

	fixture := `
{
	"status": "OK",
	"transfer": {
		"availability": null,
		"callback_url": null,
		"client_ip": null,
		"created_at": "2016-07-15T09:45:15",
		"created_torrent": false,
		"current_ratio": 0.00,
		"down_speed": 0,
		"download_id": null,
		"downloaded": 0,
		"error_message": null,
		"estimated_time": null,
		"extract": false,
		"file_id": null,
		"finished_at": null,
		"id": 1,
		"is_private": false,
		"magneturi": "magnet:?xt=urn:btih:4344503b7e797ebf31582327a5baae35b11bda01&dn=ubuntu-16.04-desktop-amd64.iso",
		"name": "ubuntu-16.04-desktop-amd64.iso",
		"peers_connected": 0,
		"peers_getting_from_us": 0,
		"peers_sending_to_us": 0,
		"percent_done": 0,
		"save_parent_id": 400004368,
		"seconds_seeding": 0,
		"simulated": false,
		"size": 1485881344,
		"source": "http://releases.ubuntu.com/16.04/ubuntu-16.04-desktop-amd64.iso.torrent",
		"status": "IN_QUEUE",
		"status_message": "In queue...",
		"subscription_id": null,
		"torrent_link": "/v2/transfers/36003178/torrent",
		"tracker_message": null,
		"trackers": null,
		"type": "TORRENT",
		"up_speed": 0,
		"uploaded": 0
	}
}
`
	mux.HandleFunc("/v2/transfers/retry", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Content-Type", "application/x-www-form-urlencoded")

		id := r.FormValue("id")
		if id != "1" {
			http.NotFound(w, r)
			return
		}

		fmt.Fprintln(w, fixture)
	})

	transfer, err := client.Transfers.Retry(nil, 1)
	if err != nil {
		t.Error(err)
	}

	if transfer.ID != 1 {
		t.Errorf("got: %v, want: %v", transfer.ID, 1)
	}

	// negative transfer ID
	_, err = client.Transfers.Retry(nil, -1)
	if err == nil {
		t.Errorf("negative transfer ID accepted")
	}

	// non-existent tranfer iD
	_, err = client.Transfers.Retry(nil, 2)
	if err != ErrResourceNotFound {
		t.Errorf("got: %v, want: %v", err, ErrResourceNotFound)
	}
}

func TestTransfers_Cancel(t *testing.T) {
	setup()
	defer teardown()

	contains := func(s []string, i string) bool {
		for _, v := range s {
			if i == v {
				return true
			}
		}
		return false
	}

	mux.HandleFunc("/v2/transfers/cancel", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Content-Type", "application/x-www-form-urlencoded")

		transferIdsStr := r.FormValue("transfer_ids")
		if transferIdsStr == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		validIds := []string{"1", "2", "3"}
		ids := strings.Split(transferIdsStr, ",")

		for _, id := range ids {
			if !contains(validIds, id) {
				http.NotFound(w, r)
				return
			}
		}

		fmt.Fprintln(w, `{"status":"OK"}`)
	})

	err := client.Transfers.Cancel(nil, 1)
	if err != nil {
		t.Error(err)
	}

	// empty params
	err = client.Transfers.Cancel(nil)
	if err == nil {
		t.Errorf("no parameters given and accepted")
	}

	// negative id
	err = client.Transfers.Cancel(nil, 1, 2, -1)
	if err == nil {
		t.Errorf("negative id accepted")
	}

	// non-existent transfer
	err = client.Transfers.Cancel(nil, 1, 2, 3, 4)
	if err != ErrResourceNotFound {
		t.Errorf("got: %v, want: %v", err, ErrResourceNotFound)
	}
}

func TestTransfers_Clean(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/transfers/clean", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Content-Type", "application/x-www-form-urlencoded")
		fmt.Fprintln(w, `{"status":"OK"}`)
	})

	err := client.Transfers.Clean(nil)
	if err != nil {
		t.Error(err)
	}
}

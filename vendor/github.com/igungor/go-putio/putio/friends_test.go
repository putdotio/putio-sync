package putio

import (
	"fmt"
	"net/http"
	"testing"
)

func TestFriends_List(t *testing.T) {
	setup()
	defer teardown()

	fixture := `
{
"friends": [
{
	"avatar_url": "",
	"id": 1,
	"name": "jet"
},
{
	"avatar_url": "",
	"id": 2,
	"name": "spike"
},
{
	"avatar_url": "",
	"id": 3,
	"name": "faye"
},
{
	"avatar_url": "",
	"id": 4,
	"name": "ein"
},
{
	"avatar_url": "",
	"id": 5,
	"name": "ed"
}
],
"status": "OK",
"total": 5
}
`

	mux.HandleFunc("/v2/friends/list", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintln(w, fixture)
	})

	friends, err := client.Friends.List(nil)
	if err != nil {
		t.Error(err)
	}

	if len(friends) != 5 {
		t.Errorf("got: %v, want: 5", len(friends))
	}

	if friends[0].ID != 1 {
		t.Errorf("got: %v, want: 1", 1)
	}

	if friends[1].Name != "spike" {
		t.Errorf("got: %v, want: spike", friends[1].Name)
	}
}

func TestFriends_WaitingRequests(t *testing.T) {
	setup()
	defer teardown()

	fixture := `
{
"friends": [
{
	"avatar_url": "",
	"id": 6,
	"name": "julia"
},
{
	"avatar_url": "",
	"id": 7,
	"name": "vicious"
}
],
"status": "OK"
}
`

	mux.HandleFunc("/v2/friends/waiting-requests", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintln(w, fixture)
	})

	friends, err := client.Friends.WaitingRequests(nil)
	if err != nil {
		t.Error(err)
	}

	if len(friends) != 2 {
		t.Errorf("got: %v, want: 2", len(friends))
	}

	if friends[0].ID != 6 {
		t.Errorf("got: %v, want: 6", friends[0].ID)
	}

	if friends[1].Name != "vicious" {
		t.Errorf("got: %v, want: vicious", friends[1].Name)
	}
}

func TestFriends_Request(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/friends/annie/request", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Content-Type", "application/x-www-form-urlencoded")
		fmt.Fprintln(w, `{"status":"OK"}`)
	})

	err := client.Friends.Request(nil, "annie")
	if err != nil {
		t.Error(err)
	}

	// empty username
	err = client.Friends.Request(nil, "")
	if err == nil {
		t.Error("empty username accepted")
	}
}

func TestFriends_Approve(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/friends/bob/approve", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Content-Type", "application/x-www-form-urlencoded")
		fmt.Fprintln(w, `{"status":"OK"}`)
	})

	err := client.Friends.Approve(nil, "bob")
	if err != nil {
		t.Error(err)
	}

	// empty username
	err = client.Friends.Approve(nil, "")
	if err == nil {
		t.Error("empty username accepted")
	}
}

func TestFriends_Deny(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/friends/andy/deny", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Content-Type", "application/x-www-form-urlencoded")
		fmt.Fprintln(w, `{"status":"OK"}`)
	})

	err := client.Friends.Deny(nil, "andy")
	if err != nil {
		t.Error(err)
	}

	// empty username
	err = client.Friends.Deny(nil, "")
	if err == nil {
		t.Error("empty username accepted")
	}
}

func TestFriends_Unfriend(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/friends/lin/unfriend", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Content-Type", "application/x-www-form-urlencoded")
		fmt.Fprintln(w, `{"status":"OK"}`)
	})

	err := client.Friends.Unfriend(nil, "lin")
	if err != nil {
		t.Error(err)
	}

	// empty username
	err = client.Friends.Unfriend(nil, "")
	if err == nil {
		t.Error("empty username accepted")
	}
}

package putio

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestFiles_Get(t *testing.T) {
	setup()
	defer teardown()

	fixture := `
{
	"file": {
		"content_type": "text/plain",
		"crc32": "66a1512f",
		"created_at": "2013-09-07T21:32:03",
		"first_accessed_at": null,
		"icon": "https://put.io/images/file_types/text.png",
		"id": 6546533,
		"is_mp4_available": false,
		"is_shared": false,
		"name": "MyFile.txt",
		"opensubtitles_hash": null,
		"parent_id": 123,
		"screenshot": null,
		"size": 92
	},
    "status": "OK"
}
`

	mux.HandleFunc("/v2/files/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintln(w, fixture)
	})
	mux.HandleFunc("/v2/files/2", http.NotFound)

	file, err := client.Files.Get(nil, 1)
	if err != nil {
		t.Error(err)
	}

	if file.Size != 92 {
		t.Errorf("got: %v, want: 92", file.Size)
	}

	// negative id
	_, err = client.Files.Get(nil, -1)
	if err == nil {
		t.Errorf("negative id accepted")
	}

	// non-existent file
	_, err = client.Files.Get(nil, 2)
	if err != ErrResourceNotFound {
		t.Errorf("got: %v, want: %v", err, ErrResourceNotFound)
	}
}

func TestFiles_List(t *testing.T) {
	setup()
	defer teardown()

	fixture := `
{
"files": [
	{
		"content_type": "text/plain",
		"crc32": "66a1512f",
		"created_at": "2013-09-07T21:32:03",
		"first_accessed_at": null,
		"icon": "https://put.io/images/file_types/text.png",
		"id": 6546533,
		"is_mp4_available": false,
		"is_shared": false,
		"name": "MyFile.txt",
		"opensubtitles_hash": null,
		"parent_id": 123,
		"screenshot": null,
		"size": 92
	},
	{
		"content_type": "video/x-matroska",
		"crc32": "cb97ba70",
		"created_at": "2013-09-07T21:32:03",
		"first_accessed_at": "2013-09-07T21:32:13",
		"icon": "https://put.io/thumbnails/aF5rkZVtYV9pV1iWimSOZWJjWWFaXGZdaZBmY2OJY4uJlV5pj5FiXg%3D%3D.jpg",
		"id": 7645645,
		"is_mp4_available": false,
		"is_shared": false,
		"name": "MyVideo.mkv",
		"opensubtitles_hash": "acc2785ffa573c69",
		"parent_id": 123,
		"screenshot": "https://put.io/screenshots/aF5rkZVtYV9pV1iWimSOZWJjWWFaXGZdaZBmY2OJY4uJlV5pj5FiXg%3D%3D.jpg",
		"size": 1155197659
	}
],
"parent": {
	"content_type": "application/x-directory",
	"crc32": null,
	"created_at": "2013-09-07T21:32:03",
	"first_accessed_at": null,
	"icon": "https://put.io/images/file_types/folder.png",
	"id": 123,
	"is_mp4_available": false,
	"is_shared": false,
	"name": "MyFolder",
	"opensubtitles_hash": null,
	"parent_id": 0,
	"screenshot": null,
	"size": 1155197751
},
"status": "OK"
}
`
	mux.HandleFunc("/v2/files/list", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		// not found handler
		parentID := r.URL.Query().Get("parent_id")
		if parentID == "2" {
			http.NotFound(w, r)
			return
		}

		fmt.Fprintln(w, fixture)
	})

	files, parent, err := client.Files.List(nil, 0)
	if err != nil {
		t.Error(err)
	}

	if len(files) != 2 {
		t.Errorf("got: %v, want: 2", len(files))
	}
	if parent.ID != 123 {
		t.Errorf("got: %v, want: 123", parent.ID)
	}

	// negative id
	_, _, err = client.Files.List(nil, -1)
	if err == nil {
		t.Errorf("negative id accepted")
	}

	// non-existent parent folder
	_, _, err = client.Files.List(nil, 2)
	if err != ErrResourceNotFound {
		t.Errorf("got: %v, want: %v", err, ErrResourceNotFound)
	}
}

func TestFiles_CreateFolder(t *testing.T) {
	setup()
	defer teardown()

	fixture := `

{
	"file": {
		"content_type": "application/x-directory",
		"crc32": null,
		"created_at": "2016-07-15T09:21:03",
		"extension": null,
		"file_type": "FOLDER",
		"first_accessed_at": null,
		"folder_type": "REGULAR",
		"icon": "https://api.put.io/images/file_types/folder.png",
		"id": 415105276,
		"is_hidden": false,
		"is_mp4_available": false,
		"is_shared": false,
		"name": "foobar",
		"opensubtitles_hash": null,
		"parent_id": 0,
		"screenshot": null,
		"size": 0
	},
	"status": "OK"
}
`
	mux.HandleFunc("/v2/files/create-folder", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Content-Type", "application/x-www-form-urlencoded")
		fmt.Fprintln(w, fixture)
	})

	file, err := client.Files.CreateFolder(nil, "foobar", 0)
	if err != nil {
		t.Error(err)
	}

	if file.Name != "foobar" {
		t.Errorf("got: %v, want: foobar", file.Name)
	}

	// empty folder name
	_, err = client.Files.CreateFolder(nil, "", 0)
	if err == nil {
		t.Errorf("empty folder name accepted")
	}

	// negative id
	_, err = client.Files.CreateFolder(nil, "foobar", -1)
	if err == nil {
		t.Errorf("negative id accepted")
	}
}

func TestFiles_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/files/delete", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Content-Type", "application/x-www-form-urlencoded")
		fmt.Fprintln(w, `{"status": "OK"}`)
	})

	err := client.Files.Delete(nil, 1, 2, 3)
	if err != nil {
		t.Error(err)
	}

	// empty params
	err = client.Files.Delete(nil)
	if err == nil {
		t.Errorf("empty parameters accepted")
	}

	err = client.Files.Delete(nil, 1, 2, -1)
	if err == nil {
		t.Errorf("negative id accepted")
	}
}

func TestFiles_Rename(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/files/rename", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Content-Type", "application/x-www-form-urlencoded")
		fmt.Fprintln(w, `{"status":"OK"}`)
	})

	err := client.Files.Rename(nil, 1, "bar")
	if err != nil {
		t.Error(err)
	}

	// negative id
	err = client.Files.Rename(nil, -1, "bar")
	if err == nil {
		t.Errorf("negative file ID accepted")
	}

	// empty name
	err = client.Files.Rename(nil, 1, "")
	if err == nil {
		t.Errorf("empty filename accepted")
	}
}

func TestFiles_Move(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/files/move", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Content-Type", "application/x-www-form-urlencoded")
		fmt.Fprintln(w, `{"status":"OK"}`)
	})

	// move 1, 2, 3, 4 and 5 to root directory (0).
	err := client.Files.Move(nil, 0, 1, 2, 3, 4, 5)
	if err != nil {
		t.Error(err)
	}

	// negative parent id
	err = client.Files.Move(nil, -1, 1, 2, 3, 4, 5)
	if err == nil {
		t.Errorf("negative parent ID accepted")
	}

	// negative file id
	err = client.Files.Move(nil, 0, 1, 2, -3)
	if err == nil {
		t.Errorf("negative file ID accepted")
	}

	// no files
	err = client.Files.Move(nil, 0)
	if err == nil {
		t.Errorf("no files given and it is accepted")
	}
}

func TestFiles_Download(t *testing.T) {
	setup()
	defer teardown()

	fileContent := "0123456789"
	mux.HandleFunc("/v2/files/1/download", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		notunnel := r.URL.Query().Get("notunnel")
		if notunnel == "1" {
			http.Redirect(w, r, "/valid-tunnel-server-download-url", http.StatusFound)
		} else {
			http.Redirect(w, r, "/valid-storage-server-download-url", http.StatusFound)
		}
	})
	mux.HandleFunc("/valid-storage-server-download-url", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		buf := strings.NewReader(fileContent)
		http.ServeContent(w, r, "testfile", time.Now().UTC(), buf)
	})
	mux.HandleFunc("/valid-tunnel-server-download-url", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		buf := strings.NewReader(fileContent)
		http.ServeContent(w, r, "testfile", time.Now().UTC(), buf)
	})

	paymentRequiredHandler := func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		errorBody := `
{
	"status": "ERROR",
	"status_code": 402,
	"error_type": "BadRequest",
	"error_message": "Payment required",
	"error_uri": "http://api.put.io/v2/docs"
}`
		w.WriteHeader(http.StatusPaymentRequired)
		fmt.Fprintln(w, errorBody)
		return
	}
	mux.HandleFunc("/v2/files/2/download", paymentRequiredHandler)

	rc, err := client.Files.Download(nil, 1, false, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer rc.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, rc)
	if err != nil {
		t.Error(err)
	}

	if buf.String() != fileContent {
		t.Errorf("got: %q, want: %q", buf.String(), fileContent)
	}

	// negative id
	rc, err = client.Files.Download(nil, -1, false, nil)
	if err == nil {
		defer rc.Close()
		t.Errorf("negative id accepted")
	}

	// tunneled download
	rc, err = client.Files.Download(nil, 1, true, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer rc.Close()

	buf.Reset()
	_, err = io.Copy(&buf, rc)
	if err != nil {
		t.Error(err)
	}
	response := buf.String()
	if response != fileContent {
		t.Errorf("got: %v, want: %v", response, fileContent)
	}

	// range request
	rangeHeader := http.Header{}
	rangeHeader.Set("Range", fmt.Sprintf("bytes=%v-%v", 0, 3))
	rc, err = client.Files.Download(nil, 1, false, rangeHeader)
	if err != nil {
		t.Fatal(err)
	}
	defer rc.Close()

	buf.Reset()
	_, err = io.Copy(&buf, rc)
	if err != nil {
		t.Error(err)
	}

	response = buf.String()
	if response != "0123" {
		t.Errorf("got: %v, want: 0123", response)
	}

	// payment required
	_, err = client.Files.Download(nil, 2, false, nil)
	if err != ErrPaymentRequired {
		t.Errorf("payment-required error should be returned")
	}
}

func TestFiles_Search(t *testing.T) {
	setup()
	defer teardown()

	fixture := `
{
"files": [
	{
		"content_type": "video/x-msvideo",
		"crc32": "812ed74d",
		"created_at": "2013-04-30T21:40:04",
		"extension": "avi",
		"file_type": "VIDEO",
		"first_accessed_at": "2013-12-24T09:18:58",
		"folder_type": "REGULAR",
		"icon": "https://some-valid-screenhost-url.com",
		"id": 79905833,
		"is_hidden": false,
		"is_mp4_available": true,
		"is_shared": false,
		"name": "some-file.mkv",
		"opensubtitles_hash": "fb5414fd9b9e1e38",
		"parent_id": 79905827,
		"screenshot": "https://some-valid-screenhost-url.com",
		"sender_name": "hafifuyku",
		"size": 738705408,
		"start_from": 0
	},
	{
		"content_type": "application/x-directory",
		"crc32": null,
		"created_at": "2013-04-30T21:40:03",
		"extension": null,
		"file_type": "FOLDER",
		"first_accessed_at": null,
		"folder_type": "REGULAR",
		"icon": "https://some-valid-screenhost-url.com",
		"id": 79905827,
		"is_hidden": false,
		"is_mp4_available": false,
		"is_shared": false,
		"name": "Movie 43",
		"opensubtitles_hash": null,
		"parent_id": 2197,
		"screenshot": null,
		"sender_name": "hafifuyku",
		"size": 738831202
	},
	{
		"content_type": "application/x-directory",
		"crc32": null,
		"created_at": "2010-05-19T22:24:21",
		"extension": null,
		"file_type": "FOLDER",
		"first_accessed_at": null,
		"folder_type": "REGULAR",
		"icon": "https://some-valid-screenhost-url.com",
		"id": 5659875,
		"is_hidden": false,
		"is_mp4_available": false,
		"is_shared": false,
		"name": "MOVIE",
		"opensubtitles_hash": null,
		"parent_id": 0,
		"screenshot": null,
		"sender_name": "emsel",
		"size": 0
	}
],
"next": null,
"status": "OK",
"total": 3
}
`
	mux.HandleFunc("/v2/files/search/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintln(w, fixture)
	})

	s, err := client.Files.Search(nil, "naber", 1)
	if err != nil {
		t.Error(err)
	}

	if len(s.Files) != 3 {
		t.Errorf("got: %v, want: 3", len(s.Files))
	}

	if s.Files[0].Name != "some-file.mkv" {
		t.Errorf("got: %v, want: some-file.mkv", s.Files[0].Name)
	}

	// invalid page number
	_, err = client.Files.Search(nil, "naber", 0)
	if err == nil {
		t.Errorf("invalid page number accepted")
	}

	// empty query
	_, err = client.Files.Search(nil, "", 1)
	if err == nil {
		t.Errorf("empty query accepted")
	}
}

func TestFiles_SetVideoPosition(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/files/1/start-from", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Content-Type", "application/x-www-form-urlencoded")
		fmt.Fprintln(w, `{"statutus":"OK"}`)
	})

	err := client.Files.SetVideoPosition(nil, 1, 10)
	if err != nil {
		t.Error(err)
	}

	// negative id
	err = client.Files.SetVideoPosition(nil, -1, 10)
	if err == nil {
		t.Errorf("negative file id accepted")
	}

	// negative time
	err = client.Files.SetVideoPosition(nil, 1, -1)
	if err == nil {
		t.Errorf("negative time accepted")
	}
}

func TestFiles_DeleteVideoPosition(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/files/1/start-from/delete", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Content-Type", "application/x-www-form-urlencoded")
		fmt.Fprintln(w, `{"statutus":"OK"}`)
	})

	err := client.Files.DeleteVideoPosition(nil, 1)
	if err != nil {
		t.Error(err)
	}

	// negative id
	err = client.Files.DeleteVideoPosition(nil, -1)
	if err == nil {
		t.Errorf("negative file id accepted")
	}
}

func TestFiles_HLSPlaylist(t *testing.T) {
	setup()
	defer teardown()

	sampleHLS := `
#EXTM3U
#EXT-X-STREAM-INF:PROGRAM-ID=1, BANDWIDTH=688301
http://qthttp.apple.com.edgesuite.net/1010qwoeiuryfg/0640_vod.m3u8
#EXT-X-STREAM-INF:PROGRAM-ID=1, BANDWIDTH=165135
http://qthttp.apple.com.edgesuite.net/1010qwoeiuryfg/0150_vod.m3u8
#EXT-X-STREAM-INF:PROGRAM-ID=1, BANDWIDTH=262346
http://qthttp.apple.com.edgesuite.net/1010qwoeiuryfg/0240_vod.m3u8
#EXT-X-STREAM-INF:PROGRAM-ID=1, BANDWIDTH=481677
http://qthttp.apple.com.edgesuite.net/1010qwoeiuryfg/0440_vod.m3u8
#EXT-X-STREAM-INF:PROGRAM-ID=1, BANDWIDTH=1308077
http://qthttp.apple.com.edgesuite.net/1010qwoeiuryfg/1240_vod.m3u8
#EXT-X-STREAM-INF:PROGRAM-ID=1, BANDWIDTH=1927853
http://qthttp.apple.com.edgesuite.net/1010qwoeiuryfg/1840_vod.m3u8
#EXT-X-STREAM-INF:PROGRAM-ID=1, BANDWIDTH=2650941
http://qthttp.apple.com.edgesuite.net/1010qwoeiuryfg/2540_vod.m3u8
#EXT-X-STREAM-INF:PROGRAM-ID=1, BANDWIDTH=3477293
http://qthttp.apple.com.edgesuite.net/1010qwoeiuryfg/3340_vod.m3u8
`
	mux.HandleFunc("/v2/files/1/hls/media.m3u8", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		http.ServeContent(w, r, "media.m3u8", time.Now().UTC(), strings.NewReader(sampleHLS))
	})

	body, err := client.Files.HLSPlaylist(nil, 1, "all")
	if err != nil {
		t.Error(err)
	}
	defer body.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, body)
	if err != nil {
		t.Error(err)
	}
	if buf.String() != sampleHLS {
		t.Errorf("got: %v, want: %v", buf.String(), sampleHLS)
	}

	// negative id
	_, err = client.Files.HLSPlaylist(nil, -1, "all")
	if err == nil {
		t.Errorf("negative file ID accepted")
	}

	// empty key
	_, err = client.Files.HLSPlaylist(nil, 1, "")
	if err == nil {
		t.Errorf("empty key is accepted")
	}
}

func TestFiles_Share(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/files/share", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Content-Type", "application/x-www-form-urlencoded")
		fmt.Fprintln(w, `{"status":"OK"}`)
	})

	err := client.Files.share(nil, []int64{1, 2, 3}, "friend0", "friend1", "friend2")
	if err != nil {
		t.Error(err)
	}

	// negative file id
	err = client.Files.share(nil, []int64{-1, 1, 2}, "friend0")
	if err == nil {
		t.Errorf("negative file id accepted")
	}

	// no file id given
	err = client.Files.share(nil, []int64{}, "friend0")
	if err == nil {
		t.Errorf("no files given and accepted")
	}

	// case: everyone (given no friends share the files to every friend)
	err = client.Files.share(nil, []int64{1})
	if err != nil {
		t.Error(err)
	}
}

func TestFiles_Shared(t *testing.T) {
	setup()
	defer teardown()

	fixture := `
{
	"shared": [
    {
		"file_id": 388029022,
		"file_name": "cowboy",
		"shared_with": 1
    },
    {
		"file_id": 388029023,
		"file_name": "bebop",
		"shared_with": 1
    }
  ],
  "status": "OK"
}
`
	mux.HandleFunc("/v2/files/shared", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintln(w, fixture)
	})

	files, err := client.Files.shared(nil)
	if err != nil {
		t.Error(err)
	}
	if len(files) != 2 {
		t.Errorf("got: %v, want: %v", len(files), 2)
	}

	if files[0].FileID != 388029022 {
		t.Errorf("got: %v, want: %v", files[0].FileID, 388029022)
	}
}

func TestFiles_SharedWith(t *testing.T) {
	setup()
	defer teardown()

	fixture := `
{
	"shared-with": [
    {
		"share_id": 1,
		"user_avatar_url": "https://some-valid-avatar-url.com/avatar.jpg",
		"user_name": "spike"
    },
    {
		"share_id": 2,
		"user_avatar_url": "https://some-valid-avatar-url.com/avatar2.jpg",
		"user_name": "edward"
    }
  ],
  "status": "OK"
}
`

	mux.HandleFunc("/v2/files/1/shared-with", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintln(w, fixture)
	})

	files, err := client.Files.sharedWith(nil, 1)
	if err != nil {
		t.Error(err)
	}
	if len(files) != 2 {
		t.Errorf("got: %v, want: %v", len(files), 2)
	}
}

func TestFiles_Subtitles(t *testing.T) {
	setup()
	defer teardown()

	fixture := `
{
	"default": "key0",
	"status": "OK",
	"subtitles": [
		{
			"key": "key0",
			"language": "Turkish",
			"language_code": "tur",
			"name": "Big Buck Bunny",
			"source": "opensubtitles"
		},
		{
			"key": "key1",
			"language": "English",
			"language_code": "eng",
			"name": "Big Buck Bunny",
			"source": "opensubtitles"
		}
	]
}
`
	mux.HandleFunc("/v2/files/1/subtitles", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintln(w, fixture)
	})

	subtitles, err := client.Files.Subtitles(nil, 1)
	if err != nil {
		t.Error(err)
	}

	if len(subtitles) != 2 {
		t.Errorf("got: %v, want: %v", len(subtitles), 2)
	}

	if subtitles[0].Key != "key0" {
		t.Errorf("got: %v, want: %v", subtitles[0].Key, "key0")
	}

	// negative id
	_, err = client.Files.Subtitles(nil, -1)
	if err == nil {
		t.Errorf("negative file ID accepted")
	}
}

func TestFiles_DownloadSubtitle(t *testing.T) {
	setup()
	defer teardown()

	fileContent := `
1
00:03:07,834 --> 00:03:09,904
Let's go down.

2
00:03:20,474 --> 00:03:24,388
- You got some out. How many left?
- Three out, eight left.
`
	mux.HandleFunc("/v2/files/1/subtitles/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		// trim leading and trailing slashes and split the url path
		f := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

		var key string
		// empty key means `default`
		if len(f) == 4 {
			key = "default"
		}
		// grab the last item of the path
		if len(f) == 5 {
			key = f[len(f)-1]
		}

		switch key {
		case "default", "key0":
			http.ServeContent(w, r, "big.buck.bunny.srt", time.Now().UTC(), strings.NewReader(fileContent))
		default:
			http.NotFound(w, r)
		}
	})

	// valid file ID and valid key
	rc, err := client.Files.DownloadSubtitle(nil, 1, "key0", "")
	if err != nil {
		t.Error(err)
	}
	defer rc.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, rc)
	if err != nil {
		t.Error(err)
	}
	if buf.String() != fileContent {
		t.Errorf("got: %v, want: %v", buf.String(), fileContent)
	}

	// negative file ID
	rc, err = client.Files.DownloadSubtitle(nil, -1, "key0", "")
	if err == nil {
		defer rc.Close()
		t.Errorf("negative file ID accepted")
	}

	// invalid key
	rc, err = client.Files.DownloadSubtitle(nil, 1, "key3", "")
	if err == nil {
		defer rc.Close()
		t.Errorf("invalid key accepted")
	}

	// empty key
	rc, err = client.Files.DownloadSubtitle(nil, 1, "", "")
	if err != nil {
		t.Error(err)
	}
	defer rc.Close()

	buf.Reset()
	_, err = io.Copy(&buf, rc)
	if err != nil {
		t.Error(err)
	}
	if buf.String() != fileContent {
		t.Errorf("got: %v, want: %v", buf.String(), fileContent)
	}
}

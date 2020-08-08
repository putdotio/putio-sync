package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/cenkalti/log"
)

const UploadURL = "https://upload.put.io/files/"

func CreateUpload(ctx context.Context, token string, filename string, parentID, length int64) (location string, err error) {
	log.Debugf("Creating upload %q at parent=%d", filename, parentID)
	req, err := http.NewRequest(http.MethodPost, UploadURL, nil)
	if err != nil {
		return
	}
	req = req.WithContext(ctx)
	metadata := map[string]string{
		"name":       filename,
		"parent_id":  strconv.FormatInt(parentID, 10),
		"no-torrent": "true",
	}
	req.Header.Set("Content-Length", "0")
	req.Header.Set("Upload-Length", strconv.FormatInt(length, 10))
	req.Header.Set("Upload-Metadata", encodeMetadata(metadata))
	req.Header.Set("Authorization", "token "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	log.Debugln("Status code:", resp.StatusCode)
	if resp.StatusCode != http.StatusCreated {
		err = fmt.Errorf("unexpected status: %d", resp.StatusCode)
		return
	}
	location = resp.Header.Get("Location")
	return
}

func SendFile(token string, r io.Reader, location string, offset int64) (fileID int64, err error) {
	log.Debugf("Sending file %q offset=%d", location, offset)
	req, err := http.NewRequestWithContext(context.TODO(), http.MethodPatch, location, r)
	if err != nil {
		return
	}

	req.Header.Set("content-type", "application/offset+octet-stream")
	req.Header.Set("upload-offset", strconv.FormatInt(offset, 10))
	req.Header.Set("Authorization", "token "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	log.Debugln("Status code:", resp.StatusCode)
	if resp.StatusCode != http.StatusNoContent {
		err = fmt.Errorf("unexpected status: %d", resp.StatusCode)
		return
	}
	return strconv.ParseInt(resp.Header.Get("putio-file-id"), 10, 64)
}

func GetUploadOffset(token string, location string) (n int64, err error) {
	log.Debugf("Getting upload offset %q", location)
	req, err := http.NewRequestWithContext(context.TODO(), http.MethodHead, location, nil)
	if err != nil {
		return
	}

	req.Header.Set("Authorization", "token "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	log.Debugln("Status code:", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("unexpected status: %d", resp.StatusCode)
		return
	}
	n, err = strconv.ParseInt(resp.Header.Get("upload-offset"), 10, 64)
	log.Debugln("Upload offset:", n)
	return n, err
}

func encodeMetadata(metadata map[string]string) string {
	encoded := make([]string, 0, len(metadata))
	for k, v := range metadata {
		encoded = append(encoded, fmt.Sprintf("%s %s", k, base64.StdEncoding.EncodeToString([]byte(v))))
	}
	return strings.Join(encoded, ",")
}

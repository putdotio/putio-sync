package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cenkalti/log"
)

const UploadURL = "https://upload.put.io/files/"

func CreateUpload(baseCtx context.Context, token string, filename string, parentID, length int64) (location string, err error) {
	log.Debugf("Creating upload %q at parent=%d", filename, parentID)
	ctx, cancel := context.WithTimeout(baseCtx, defaultTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, UploadURL, nil)
	if err != nil {
		return
	}
	metadata := map[string]string{
		"name":       filename,
		"parent_id":  strconv.FormatInt(parentID, 10),
		"no-torrent": "true",
	}
	req.Header.Set("Content-Length", "0")
	req.Header.Set("Upload-Length", strconv.FormatInt(length, 10))
	req.Header.Set("Upload-Metadata", encodeMetadata(metadata))
	req.Header.Set("Authorization", "token "+token)

	resp, err := httpClient.Do(req)
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

func SendFile(ctx context.Context, token string, r io.Reader, location string, offset int64) (fileID int64, crc32 string, err error) {
	log.Debugf("Sending file %q offset=%d", location, offset)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Stop upload if speed is too slow.
	// Wrap reader so each read call resets the timer that cancels the request on certain duration.
	r = &TimerResetReader{r: r, timer: time.AfterFunc(defaultTimeout, cancel)}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, location, r)
	if err != nil {
		return
	}

	req.Header.Set("content-type", "application/offset+octet-stream")
	req.Header.Set("upload-offset", strconv.FormatInt(offset, 10))
	req.Header.Set("Authorization", "token "+token)
	resp, err := httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	log.Debugln("Status code:", resp.StatusCode)
	if resp.StatusCode != http.StatusNoContent {
		err = fmt.Errorf("unexpected status: %d", resp.StatusCode)
		return
	}
	fileID, err = strconv.ParseInt(resp.Header.Get("putio-file-id"), 10, 64)
	if err != nil {
		err = fmt.Errorf("cannot parse putio-file-id header: %w", err)
		return
	}
	crc32 = resp.Header.Get("putio-file-crc32")
	return
}

func GetUploadOffset(ctx context.Context, token string, location string) (n int64, err error) {
	log.Debugf("Getting upload offset %q", location)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, location, nil)
	if err != nil {
		return
	}

	req.Header.Set("Authorization", "token "+token)
	resp, err := httpClient.Do(req)
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

func TerminateUpload(ctx context.Context, token string, location string) (err error) {
	log.Debugf("Terminating upload %q", location)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, location, nil)
	if err != nil {
		return
	}

	req.Header.Set("Authorization", "token "+token)
	resp, err := httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	log.Debugln("Status code:", resp.StatusCode)
	if resp.StatusCode != http.StatusNoContent {
		err = fmt.Errorf("unexpected status: %d", resp.StatusCode)
		return
	}
	return nil
}

func encodeMetadata(metadata map[string]string) string {
	encoded := make([]string, 0, len(metadata))
	for k, v := range metadata {
		encoded = append(encoded, fmt.Sprintf("%s %s", k, base64.StdEncoding.EncodeToString([]byte(v))))
	}
	return strings.Join(encoded, ",")
}

type TimerResetReader struct {
	r     io.Reader
	timer *time.Timer
}

func (r *TimerResetReader) Read(p []byte) (int, error) {
	r.timer.Reset(defaultTimeout)
	return r.r.Read(p)
}

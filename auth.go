package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/cenkalti/log"
	"github.com/putdotio/go-putio"
	"go.etcd.io/bbolt"
	"golang.org/x/oauth2"
)

const (
	baseURL        = "https://api.put.io"
	tokenURL       = "https://api.put.io/v2/oauth2/authorizations/clients/"
	clientID       = "4785"
	clientSecret   = "YGRIVM3BKAPGTYCR7PEC"
	defaultTimeout = 10 * time.Second
)

var (
	bucketConfig = []byte("config")
	keyToken     = []byte("token")
)

func ensureValidClient() error {
	token, err := getExistingToken()
	if err != nil {
		return err
	}
	if token == "" {
		// no existing token
		return authenticate()
	}
	c := newPutioClient(token)
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	_, err = c.ValidateToken(ctx)
	if err != nil {
		if er, ok := err.(*putio.ErrorResponse); ok && er.Response.StatusCode == 401 {
			// existing token is not valid anymore
			return authenticate()
		}
		return err
	}
	// existing token is valid
	client = c
	return nil
}

func getExistingToken() (string, error) {
	var token string
	err := db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(bucketConfig)
		if err != nil {
			return err
		}
		val := b.Get(keyToken)
		if val != nil {
			token = string(val)
		}
		return nil
	})
	return token, err
}

func authenticate() error {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	fingerprint := url.PathEscape(hostname)
	clientName := url.QueryEscape(hostname)
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "PUT", tokenURL+clientID+"/"+fingerprint+"?client_secret="+clientSecret+"&client_name="+clientName, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(config.Username, config.Password)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 403 {
		return errors.New("invalid username or password")
	}

	if resp.StatusCode >= 400 && resp.StatusCode <= 499 {
		return fmt.Errorf("client error, status: %v", resp.Status)
	}

	if resp.StatusCode >= 500 {
		return fmt.Errorf("server error while requesting bearer token, status: %v", resp.Status)
	}

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
		UserID      int64  `json:"user_id"`
	}

	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	if err != nil {
		return fmt.Errorf("json decode error: %v", err)
	}

	client = newPutioClient(tokenResponse.AccessToken)
	return nil
}

func newPutioClient(token string) *putio.Client {
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	oauthClient := oauth2.NewClient(context.Background(), tokenSource)
	client := putio.NewClient(oauthClient)
	client.BaseURL, _ = url.Parse(baseURL)
	return client
}

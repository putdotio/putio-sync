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
	"golang.org/x/oauth2"
)

const (
	tokenURL       = "https://api.put.io/v2/oauth2/authorizations/clients/"
	clientID       = "4785"
	clientSecret   = "YGRIVM3BKAPGTYCR7PEC"
	defaultTimeout = 10 * time.Second
)

func authenticate() error {
	log.Infof("Authenticating as user: %q", config.Username)
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
	}

	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	if err != nil {
		return fmt.Errorf("json decode error: %v", err)
	}

	token = tokenResponse.AccessToken
	oauthToken := &oauth2.Token{AccessToken: token}
	tokenSource := oauth2.StaticTokenSource(oauthToken)
	oauthClient := oauth2.NewClient(context.Background(), tokenSource)
	client = putio.NewClient(oauthClient)
	return nil
}

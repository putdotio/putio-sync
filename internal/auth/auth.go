package auth

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
	tokenURL     = "https://api.put.io/v2/oauth2/authorizations/clients/" // nolint: gosec
	clientID     = "4785"
	clientSecret = "YGRIVM3BKAPGTYCR7PEC" // nolint: gosec
)

func Authenticate(ctx context.Context, httpClient *http.Client, timeout time.Duration, username, password string) (token string, client *putio.Client, err error) {
	log.Infof("Authenticating as user: %q", username)
	hostname, err := os.Hostname()
	if err != nil {
		return
	}
	fingerprint := url.PathEscape(hostname)
	clientName := url.QueryEscape(hostname)
	authCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(authCtx, "PUT", tokenURL+clientID+"/"+fingerprint+"?client_secret="+clientSecret+"&client_name="+clientName, nil)
	if err != nil {
		return
	}
	req.SetBasicAuth(username, password)

	resp, err := httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden {
		err = errors.New("invalid username or password")
		return
	}

	if resp.StatusCode >= http.StatusBadRequest && resp.StatusCode < http.StatusInternalServerError {
		err = fmt.Errorf("client error, status: %v", resp.Status)
		return
	}

	if resp.StatusCode >= http.StatusInternalServerError {
		err = fmt.Errorf("server error while requesting bearer token, status: %v", resp.Status)
		return
	}

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
	}

	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	if err != nil {
		err = fmt.Errorf("json decode error: %v", err)
		return
	}

	token = tokenResponse.AccessToken
	oauthToken := &oauth2.Token{AccessToken: token}
	tokenSource := oauth2.StaticTokenSource(oauthToken)
	clientCtx := context.WithValue(ctx, oauth2.HTTPClient, httpClient)
	oauthClient := oauth2.NewClient(clientCtx, tokenSource)
	client = putio.NewClient(oauthClient)
	return
}

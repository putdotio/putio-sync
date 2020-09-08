package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/cenkalti/log"
	"github.com/putdotio/go-putio"
	"golang.org/x/oauth2"
)

const (
	tokenURL      = "https://api.put.io/v2/oauth2/authorizations/clients/" // nolint: gosec
	clientID      = "4785"
	clientSecret  = "YGRIVM3BKAPGTYCR7PEC" // nolint: gosec
	clientTimeout = 10 * time.Second
)

var ErrInvalidCredentials = errors.New("invalid credentials")

func Authenticate(ctx context.Context, httpClient *http.Client, timeout time.Duration, username, password string) (token string, client *putio.Client, err error) {
	if strings.HasPrefix(password, "token/") {
		// User may use a token instead of password.
		log.Infof("Validating authentication token")
		token = password[6:]
		client = newClient(ctx, httpClient, token)
		authCtx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		userID, verr := client.ValidateToken(authCtx)
		if verr != nil {
			err = verr
			return
		}
		if userID == nil {
			err = ErrInvalidCredentials
		}
		return
	}
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

	if resp.StatusCode == http.StatusUnauthorized {
		err = ErrInvalidCredentials
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
	client = newClient(ctx, httpClient, token)
	return
}

func newClient(ctx context.Context, httpClient *http.Client, token string) *putio.Client {
	oauthToken := &oauth2.Token{AccessToken: token}
	tokenSource := oauth2.StaticTokenSource(oauthToken)
	clientCtx := context.WithValue(ctx, oauth2.HTTPClient, httpClient)
	oauthClient := oauth2.NewClient(clientCtx, tokenSource)
	client := putio.NewClient(oauthClient)
	client.Timeout = clientTimeout
	return client
}

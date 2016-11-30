package putio

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	defaultUserAgent = "go-putio"
	defaultMediaType = "application/json"
	defaultBaseURL   = "https://api.put.io"
	defaultUploadURL = "https://upload.put.io"
)

type Error string

func (e Error) Error() string { return string(e) }

const (
	ErrResourceNotFound = Error("resource does not exist")
	ErrPaymentRequired  = Error("payment required")

	errRedirect   = Error("redirect attempt on a no-redirect client")
	errNegativeID = Error("file id cannot be negative")
)

// Client manages communication with Put.io v2 API.
type Client struct {
	// HTTP client used to communicate with Put.io API
	client *http.Client

	// Base URL for API requests
	BaseURL *url.URL

	// base url for upload requests
	uploadURL *url.URL

	// User agent for client
	UserAgent string

	// Services used for communicating with the API
	Account   *AccountService
	Files     *FilesService
	Transfers *TransfersService
	Zips      *ZipsService
	Friends   *FriendsService
	Events    *EventsService
}

// NewClient returns a new Put.io API client, using the htttpClient, which must
// be a new Oauth2 enabled http.Client. If httpClient is not defined, default
// HTTP client is used.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, _ := url.Parse(defaultBaseURL)
	uploadURL, _ := url.Parse(defaultUploadURL)
	c := &Client{
		client:    httpClient,
		BaseURL:   baseURL,
		uploadURL: uploadURL,
		UserAgent: defaultUserAgent,
	}

	// redirect once client. it's necessary to create a new client just for
	// download operations.
	roc := *c
	roc.client.CheckRedirect = redirectOnceFunc

	c.Account = &AccountService{client: c}
	c.Files = &FilesService{client: c, redirectOnceClient: &roc}
	c.Transfers = &TransfersService{client: c}
	c.Zips = &ZipsService{client: c}
	c.Friends = &FriendsService{client: c}
	c.Events = &EventsService{client: c}

	return c
}

// NewRequest creates an API request. A relative URL can be provided via
// relURL, which will be resolved to the BaseURL of the Client.
func (c *Client) NewRequest(ctx context.Context, method, relURL string, body io.Reader) (*http.Request, error) {
	rel, err := url.Parse(relURL)
	if err != nil {
		return nil, err
	}

	var u *url.URL
	// XXX: workaroud for upload endpoint. upload method has a different base url,
	// so we've a special case for testing purposes.
	if relURL == "/v2/files/upload" {
		u = c.uploadURL.ResolveReference(rel)
	} else {
		u = c.BaseURL.ResolveReference(rel)
	}

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	if ctx == nil {
		ctx = context.Background()
	}
	req = req.WithContext(ctx)

	req.Header.Set("Accept", defaultMediaType)
	req.Header.Set("User-Agent", c.UserAgent)

	return req, nil
}

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred. Response body is closed at all cases except
// v is nil. If v is nil, response body is not closed and the body can be used
// for streaming.
func (c *Client) Do(r *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(r)
	if err != nil {
		return nil, err
	}

	err = checkResponse(resp)
	if err != nil {
		// close the body at all times if there is an http error
		resp.Body.Close()
		return resp, err
	}

	if v == nil {
		return resp, nil
	}

	// close the body for all cases from here
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(v)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// redirectOnceFunc follows the redirect only once, and copies the original
// request headers to the new one.
func redirectOnceFunc(req *http.Request, via []*http.Request) error {
	if len(via) == 0 {
		return nil
	}

	if len(via) > 1 {
		return errRedirect
	}

	// merge headers with request headers
	for header, values := range via[0].Header {
		for _, value := range values {
			req.Header.Add(header, value)
		}
	}
	return nil
}

// ErrorResponse reports the error caused by an API request.
type ErrorResponse struct {
	Response *http.Response `json:"-"`

	Message string `json:"error_message"`
	Type    string `json:"error_type"`
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf(
		"Type: %v Message: %q. Original error: %v %v: %v",
		e.Type,
		e.Message,
		e.Response.Request.Method,
		e.Response.Request.URL,
		e.Response.Status,
	)
}

// checkResponse is the entrypoint to reading the API response. If the response
// status code is not in success range, it will try to return a structured
// error.
func checkResponse(r *http.Response) error {
	statusCode := r.StatusCode
	if statusCode >= 200 && statusCode <= 299 {
		return nil
	}

	if statusCode == http.StatusNotFound {
		return ErrResourceNotFound
	}

	if statusCode == http.StatusPaymentRequired {
		return ErrPaymentRequired
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, errorResponse)
		if err != nil {
			return err
		}
	}
	return errorResponse
}

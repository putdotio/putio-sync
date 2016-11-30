package sync

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cenk/backoff"
)

const backoffMaxElapsedTime = 1 * time.Minute

// retryTransport wraps an http.RoundTripper and retries the same request over
// and over again until a certain treshold is met.
type retryTransport struct {
	transport http.RoundTripper
}

var _ http.RoundTripper = &retryTransport{}

func (r *retryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	expbackoff := backoff.NewExponentialBackOff()
	expbackoff.MaxElapsedTime = backoffMaxElapsedTime

	var response *http.Response
	var err error
	op := func() error {
		response, err = r.transport.RoundTrip(req)
		if err != nil {
			return err
		}
		if response.StatusCode >= http.StatusInternalServerError {
			return fmt.Errorf("Unexpected HTTP Status: %v", response.Status)
		}
		return nil
	}

	// TODO: backoff.Retry is not cancellable by context.Context.
	err = backoff.Retry(op, expbackoff)
	return response, err
}

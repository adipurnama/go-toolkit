package httpclient

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	phttp "github.com/pinpoint-apm/pinpoint-go-agent/plugin/http"
	"github.com/pkg/errors"
)

// ErrNonNilContext when required context param is nil.
var ErrNonNilContext = errors.New("web/httpclient: context cannot be nil")

// ContextHTTPClient wrapper of http.Client with required context param.
type ContextHTTPClient struct {
	client *http.Client
}

// NewContextHTTPClient returns instance of ContextHTTPClient.
func NewContextHTTPClient(c *http.Client) *ContextHTTPClient {
	client := phttp.WrapClient(c)

	return &ContextHTTPClient{client: client}
}

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred. If v implements the io.Writer interface,
// the raw response body will be written to v, without attempting to first
// decode it. If v is nil, and no error hapens, the response is returned as is.
// If rate limit is exceeded and reset time is in the future, Do returns
// *RateLimitError immediately without making a network API call.
//
// The provided ctx must be non-nil, if it is nil an error is returned. If it
// is canceled or times out, ctx.Err() will be returned.
func (c *ContextHTTPClient) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	if ctx == nil {
		return nil, ErrNonNilContext
	}

	req, err := http.NewRequestWithContext(ctx, req.Method, req.URL.String(), req.Body)
	if err != nil {
		return nil, errors.Wrap(err, "web/httpclient: failed to build new request with context")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		return nil, err
	}

	defer resp.Body.Close()

	switch v := v.(type) {
	case nil:
	case io.Writer:
		_, err = io.Copy(v, resp.Body)
	default:
		decErr := json.NewDecoder(resp.Body).Decode(v)
		if errors.Is(decErr, io.EOF) {
			decErr = nil // ignore EOF errors caused by empty response body
		}

		if decErr != nil {
			err = decErr
		}
	}

	return resp, errors.Wrap(err, "web/httpclient: Do request")
}

package loki

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"

	"github.com/grafana/loki/pkg/logproto"
)

const (
	queryPath = "/api/prom/query?query=%s&limit=%d&start=%d&end=%d&direction=%s&regexp=%s"
)

// Client is an API client for Loki
type Client interface {
	// Query will request logs from Loki, returning a parsed log stream response
	Query(context.Context, logproto.QueryRequest) (*logproto.QueryResponse, error)
}

// http.Client wrapper for adding new methods, particularly sendJsonReq
type httpClient struct {
	parent      *http.Client
	lokiBaseURL string
}

// New returns a new Loki httpClient
func New(base string) Client {
	return &httpClient{
		parent:      http.DefaultClient,
		lokiBaseURL: base,
	}
}

func (c *httpClient) Query(ctx context.Context, req logproto.QueryRequest) (*logproto.QueryResponse, error) {
	path := fmt.Sprintf(
		queryPath,
		url.QueryEscape(req.Query),
		req.Limit,
		req.Start.UnixNano(),
		req.End.UnixNano(),
		req.Direction.String(),
		req.Regex,
	)

	var resp logproto.QueryResponse
	err := c.doRequest(ctx, path, &resp)
	if err != nil {
		return nil, errors.Wrap(err, "failed query")
	}

	return &resp, nil
}

func (c *httpClient) doRequest(ctx context.Context, path string, out interface{}) error {
	log := logrus.WithField("method", "doRequest")

	req, err := http.NewRequest("GET", c.lokiBaseURL+path, nil)
	if err != nil {
		return errors.Wrap(err, "failed to build Loki request")
	}
	req = req.WithContext(ctx)

	log.Debug(req.URL.String())

	resp, err := c.parent.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed request")
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Errorf("error closing body: %s", err)
		}
	}()

	if !isOK(resp.StatusCode) {
		return errors.Wrap(err, "error response from server")
	}

	return json.NewDecoder(resp.Body).Decode(out)
}

func isOK(s int) bool {
	return s/100 == 2
}

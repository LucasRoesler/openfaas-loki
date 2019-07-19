package loki

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/grafana/loki/pkg/logproto"
)

const forward = "FORWARD"
const backward = "BACKWARD"

func Test_RequestConstruction(t *testing.T) {
	ctx := context.Background()

	now := time.Now()
	nowUnix := fmt.Sprintf("%d", now.UnixNano())
	cases := []struct {
		name     string
		req      logproto.QueryRequest
		expected expectedQueryParmas
	}{
		{
			"query value only",
			logproto.QueryRequest{Query: "test"},
			expectedQueryParmas{query: "test", direction: forward},
		},
		{
			"query with limit",
			logproto.QueryRequest{Query: "test", Limit: 3},
			expectedQueryParmas{query: "test", direction: forward, limit: "3"},
		},
		{
			"query with start",
			logproto.QueryRequest{Query: "test", Start: now},
			expectedQueryParmas{query: "test", direction: forward, start: nowUnix},
		},
		{
			"query with end equal now",
			logproto.QueryRequest{Query: "test", End: now},
			expectedQueryParmas{query: "test", direction: forward, end: nowUnix},
		},
		{
			"query with regex",
			logproto.QueryRequest{Query: "test", Regex: "^abc$"},
			expectedQueryParmas{query: "test", direction: forward, regexp: "^abc$"},
		},
		{
			"query with backward direction",
			logproto.QueryRequest{Query: "test", Direction: logproto.BACKWARD},
			expectedQueryParmas{query: "test", direction: backward},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, http.MethodGet, r.Method)
				require.Equal(t, "/api/prom/query", r.URL.Path)

				query := r.URL.Query()
				require.Equal(t, tc.expected.query, query.Get("query"))
				require.Equal(t, tc.expected.limit, query.Get("limit"))
				require.Equal(t, tc.expected.start, query.Get("start"))
				require.Equal(t, tc.expected.end, query.Get("end"))
				require.Equal(t, tc.expected.direction, query.Get("direction"))
				require.Equal(t, tc.expected.regexp, query.Get("regexp"))
				fmt.Fprintln(w, "{}")
			}))
			c := New(ts.URL)

			resp, err := c.Query(ctx, tc.req)
			require.NoError(t, err)
			require.Equal(t, 0, len(resp.Streams))
		})
	}
}

type expectedQueryParmas struct {
	query     string
	limit     string
	start     string
	end       string
	direction string
	regexp    string
}

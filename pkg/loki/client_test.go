package loki

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/grafana/loki/pkg/logproto"
)

func Test_RequestConstruction(t *testing.T) {
	ctx := context.Background()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/api/prom/query", r.URL.Path)

		query := r.URL.Query()
		require.Equal(t, "test", query.Get("query"))
		require.Equal(t, "0", query.Get("limit"))
		require.Equal(t, "", query.Get("start"))
		require.Equal(t, "", query.Get("end"))
		require.Equal(t, "", query.Get("direction"))
		require.Equal(t, "", query.Get("regexp"))
		fmt.Fprintln(w, "{}")
	}))
	c := New(ts.URL)

	req := logproto.QueryRequest{Query: "test"}
	resp, err := c.Query(ctx, req)
	require.NoError(t, err)
	require.Equal(t, 0, len(resp.Streams))
}

package faas

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/pkg/errors"

	"github.com/stretchr/testify/require"

	"github.com/grafana/loki/pkg/logproto"

	"github.com/openfaas/faas-provider/logs"
)

func Test_QueryErrorIsPropagated(t *testing.T) {
	ctx := context.Background()
	lokiClient := mockLokiClient{err: errors.New("query failed")}
	requester := New(&lokiClient)

	_, err := requester.Query(ctx, logs.Request{Name: "testFnc"})
	require.EqualError(t, errors.Cause(err), "query failed")
}

func Test_QueryWithInvalidLabels(t *testing.T) {
	ctx := context.Background()

	now := time.Now()
	lokiClient := mockLokiClient{
		labels: map[string]string{"somethingelse": "and more"},
		entries: []logproto.Entry{
			{Timestamp: now, Line: "test message0"},
		},
	}
	requester := New(&lokiClient)

	stream, err := requester.Query(ctx, logs.Request{Name: "testFnc"})
	require.NoError(t, err)

	var received int
	for msg := range stream {
		require.Equal(t, now.Add(time.Second*time.Duration(received)), msg.Timestamp)
		require.Equal(t, fmt.Sprintf("test message%v", received), msg.Text)
		require.Empty(t, msg.Name)
		received++
	}
	require.Equal(t, len(lokiClient.entries), received)
}

func Test_QueryHappyPath(t *testing.T) {
	ctx := context.Background()

	now := time.Now()
	lokiClient := mockLokiClient{
		labels: map[string]string{"faas_function": "testFnc"},
		entries: []logproto.Entry{
			{Timestamp: now, Line: "test message0"},
			{Timestamp: now.Add(time.Second), Line: "test message1"},
			{Timestamp: now.Add(2 * time.Second), Line: "test message2"},
		},
	}
	requester := New(&lokiClient)

	stream, err := requester.Query(ctx, logs.Request{Name: "testFnc"})
	require.NoError(t, err)

	var received int
	for msg := range stream {
		require.Equal(t, now.Add(time.Second*time.Duration(received)), msg.Timestamp)
		require.Equal(t, fmt.Sprintf("test message%v", received), msg.Text)
		require.Equal(t, "testFnc", msg.Name)
		received++
	}
	require.Equal(t, len(lokiClient.entries), received)
}

type mockLokiClient struct {
	err     error
	entries []logproto.Entry
	labels  map[string]string
}

// Query will request logs from Loki, returning a parsed log stream response
func (c *mockLokiClient) Query(ctx context.Context, req logproto.QueryRequest) (*logproto.QueryResponse, error) {
	if c.err != nil {
		return nil, c.err
	}

	var b strings.Builder
	for label, value := range c.labels {
		b.WriteString(label)
		b.WriteString("=\"")
		b.WriteString(value)
		b.WriteString("\",")
	}
	// built based on the example from https://github.com/grafana/loki/blob/master/docs/api.md
	labels := fmt.Sprintf("{%s}", strings.TrimRight(b.String(), ","))

	resp := &logproto.QueryResponse{
		Streams: []logproto.Stream{
			{Labels: labels, Entries: c.entries},
		},
	}
	return resp, nil
}

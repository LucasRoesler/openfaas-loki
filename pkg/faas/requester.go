package faas

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"sync"

	"github.com/pkg/errors"

	"github.com/grafana/loki/pkg/logproto"

	"github.com/LucasRoesler/openfaas-loki/pkg/loki"

	"github.com/openfaas/faas-provider/logs"
)

type lokiRequester struct {
	client loki.Client
}

// New returns a new Loki log Requester
func New(client loki.Client) logs.Requester {
	return &lokiRequester{
		client: client,
	}
}

// Query submits a log request to the actual logging system.
func (l *lokiRequester) Query(ctx context.Context, req logs.Request) (<-chan logs.Message, error) {
	logger := slog.Default().With("method", "Query").With("name", req.Name)

	queryReq := l.buildRequest(req)
	resp, err := l.client.Query(ctx, queryReq)
	if err != nil {
		return nil, errors.Wrap(err, "failed requester request")
	}

	var wg sync.WaitGroup
	logStream := make(chan logs.Message, 100)
	for _, stream := range resp.Streams {
		logger.LogAttrs(ctx, slog.LevelDebug, "starting stream")
		wg.Add(1)
		go func(s logproto.Stream) {
			defer wg.Done()
			// TODO: is this safe, since stream is []*logproto.Stream?
			l.sendEntries(ctx, logStream, s)
		}(stream)
	}

	go func() {
		wg.Wait()
		logger.LogAttrs(ctx, slog.LevelDebug, "all streams closed")
		close(logStream)
	}()

	return logStream, err
}

// buildRequest currently ignores Follow
func (l *lokiRequester) buildRequest(logReq logs.Request) (req logproto.QueryRequest) {
	if logReq.Tail > 0 {
		req.Limit = uint32(logReq.Tail)
	}
	if logReq.Since != nil {
		req.Start = *logReq.Since
	}

	if logReq.Instance != "" {
		req.Selector = fmt.Sprintf("{faas_function=\"%s\",instance=\"%s\"}", logReq.Name, logReq.Instance)
	} else {
		req.Selector = fmt.Sprintf("{faas_function=\"%s\"}", logReq.Name)
	}
	slog.LogAttrs(context.TODO(), slog.LevelDebug, "buildRequest", slog.Any("request", logReq), slog.Any("lokiRequest", req))
	return req
}

// sendEntries will parse the stream entries and push them into the log stream channel
func (l *lokiRequester) sendEntries(ctx context.Context, logStream chan logs.Message, stream logproto.Stream) {
	logger := slog.Default().With("method", "sendEntries")
	labels := parseLabels(stream.Labels)
	for _, entry := range stream.Entries {
		if ctx.Err() != nil {
			logger.LogAttrs(ctx, slog.LevelDebug, "context canceled, stopping stream")
			return
		}
		logStream <- parseEntry(entry, labels)
	}
}

// parseLabels parse the output of Labels.String from
// from prometheus https://github.com/prometheus/prometheus/blob/8624913a3489f28a173f5c6e49a8a762785ae2c0/pkg/labels/labels.go#L49
// because this is currently the format sen back by Loki
// parsing errors are quiently skipped
func parseLabels(value string) map[string]string {
	logger := slog.With("method", "parseLabels")
	logger.LogAttrs(context.TODO(), slog.LevelDebug, value)

	parsed := map[string]string{}

	labelCSV := strings.Trim(value, "{}")
	labels := strings.Split(labelCSV, ",")
	for _, label := range labels {
		parts := strings.SplitN(strings.TrimSpace(label), "=", 2)
		if len(parts) != 2 {
			logger.LogAttrs(context.TODO(), slog.LevelError, "unexpected number of label parts", slog.String("label", label))
			continue
		}

		value, err := strconv.Unquote(parts[1])
		if err != nil {
			logger.LogAttrs(context.TODO(), slog.LevelError, "failed to unquote label value", slog.String("label", label), slog.String("err", err.Error()))
			continue
		}
		parsed[parts[0]] = value
	}

	return parsed
}

// parseEntry
func parseEntry(entry logproto.Entry, labels map[string]string) logs.Message {
	return logs.Message{
		Name:      labels["faas_function"],
		Instance:  labels["instance"],
		Timestamp: entry.Timestamp,
		Text:      entry.Line,
	}
}

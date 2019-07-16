package faas

import (
	"context"
	"fmt"
	"sync"

	"github.com/prometheus/prometheus/promql"

	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"

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
	log := logrus.WithField("method", "Query").WithField("name", req.Name)

	queryReq := l.buildRequest(req)
	resp, err := l.client.Query(ctx, queryReq)
	if err != nil {
		return nil, errors.Wrap(err, "failed requester request")
	}

	var wg sync.WaitGroup
	logStream := make(chan logs.Message, 100)
	for _, stream := range resp.Streams {
		log.Debug("starting stream")
		wg.Add(1)
		go func(s *logproto.Stream) {
			defer wg.Done()
			// TODO: is this safe, since stream is []*logproto.Stream?
			l.sendEntries(ctx, logStream, s)
		}(stream)
	}

	go func() {
		wg.Wait()
		log.Debug("all streams closed")
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
		req.Query = fmt.Sprintf("{faas_function=\"%s\",instance=\"%s\"}", logReq.Name, logReq.Instance)
	} else {
		req.Query = fmt.Sprintf("{faas_function=\"%s\"}", logReq.Name)
	}
	logrus.WithField("method", "buildRequest").Debugf("%v => %v", logReq, req)
	return req
}

// sendEntries will parse the stream entries and push them into the log stream channel
func (l *lokiRequester) sendEntries(ctx context.Context, logStream chan logs.Message, stream *logproto.Stream) {
	log := logrus.WithField("method", "sendEntries")
	if stream == nil {
		log.Debug("received nil stream")
		return
	}

	labels := parseLabels(stream.Labels)
	for _, entry := range stream.Entries {
		if ctx.Err() != nil {
			log.Debug("context cancelled, stopping stream")
			return
		}
		logStream <- parseEntry(entry, labels)
	}
}

// parseLabels uses the prometheus labels parser to convert the labels into a map
func parseLabels(value string) map[string]string {
	ls, err := promql.ParseMetric(value)
	if err != nil {
		logrus.
			WithField("method", "parseLabels").
			WithField("value", value).
			Error(err)
		return nil
	}

	return ls.Map()
}

// parseEntry
func parseEntry(entry logproto.Entry, labels map[string]string) logs.Message {
	return logs.Message{
		Name:      labels["faas_function"],
		Timestamp: entry.Timestamp,
		Text:      entry.Line,
	}
}

package middlewares

import (
	"bufio"
	"errors"
	"net"
	"net/http"
)

// ResponseWriter wraps the standard http.ResponseWriter to provide additional information about
// the response.
type ResponseWriter interface {
	http.ResponseWriter
	http.Flusher
	// Status returns the status code of the response or 0 if the response has
	// not been written
	Status() int
	// Written returns whether or not the ResponseWriter has been written.
	Written() bool
	// Size returns the size of the response body.
	Size() int
}

// NewResponseWriter creates a ResponseWriter that wraps an http.ResponseWriter
func NewResponseWriter(rw http.ResponseWriter) ResponseWriter {
	nrw := &responseWriter{
		ResponseWriter: rw,
	}

	//nolint:staticcheck // ignore deprecation error
	_, ok := rw.(http.CloseNotifier)
	if ok {
		return &responseWriterCloseNotifer{nrw}
	}

	return nrw
}

type responseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (rw *responseWriter) WriteHeader(s int) {
	if rw.Written() {
		return
	}
	rw.status = s
	rw.ResponseWriter.WriteHeader(s)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.Written() {
		// The status will be StatusOK if WriteHeader has not been called yet
		rw.WriteHeader(http.StatusOK)
	}
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) Size() int {
	return rw.size
}

func (rw *responseWriter) Written() bool {
	return rw.status != 0
}

func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := rw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("the ResponseWriter doesn't support the Hijacker interface")
	}
	return hijacker.Hijack()
}

func (rw *responseWriter) Flush() {
	flusher, ok := rw.ResponseWriter.(http.Flusher)
	if ok {
		if !rw.Written() {
			// The status will be StatusOK if WriteHeader has not been called yet
			rw.WriteHeader(http.StatusOK)
		}
		flusher.Flush()
	}
}

// Deprecated: the CloseNotifier interface predates Go's context package.
// New code should use Request.Context instead.
//
// We need to update the log provider before we can remove this.
type responseWriterCloseNotifer struct {
	*responseWriter
}

func (rw *responseWriterCloseNotifer) CloseNotify() <-chan bool {
	//nolint:staticcheck // we need to update the log provider before we can remove this.
	return rw.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

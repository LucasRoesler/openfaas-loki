package middlewares

import (
	"log/slog"
	"net/http"
	"time"
)

// RequestLogger returns a logger handler using a custom LogFormatter.
func RequestLogger(skipList []string) func(next http.Handler) http.Handler {
	skip := map[string]bool{}
	for _, s := range skipList {
		skip[s] = true
	}

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if skip[r.URL.Path] {
				next.ServeHTTP(w, r)
				return
			}

			ww := NewResponseWriter(w)

			t1 := time.Now()
			defer func() {
				slog.Default().LogAttrs(
					r.Context(),
					slog.LevelInfo, "request finished",
					slog.String("method", r.Method),
					slog.String("url", r.URL.String()),
					slog.Int("status", ww.Status()),
					slog.Int("size", ww.Size()),
					slog.Duration("duration", time.Since(t1)),
				)
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}

package middlewares

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
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
				log.Info().
					Str("method", r.Method).
					Str("url", r.URL.String()).
					Int("status", ww.Status()).
					Dur("duration", time.Since(t1)).
					Int("size", ww.Size()).
					Msg("request finished")
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}

package httputil

import (
	"context"
	"net"
	"net/http"

	"github.com/felixge/httpsnoop"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gorilla/mux"

	"github.com/micromdm/micromdm/pkg/logutil"
)

func NewLogger(base log.Logger) mux.MiddlewareFunc {
	h := handler{base}
	return h.decorate
}

type handler struct{ logger log.Logger }

func (h handler) decorate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = logutil.NewContext(ctx, h.logger)
		var metrics httpsnoop.Metrics

		defer func() {
			logRequest(ctx, metrics.Code, r)
		}()

		metrics = httpsnoop.CaptureMetrics(next, w, r.WithContext(ctx))
	})
}

func logRequest(ctx context.Context, code int, r *http.Request) {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}

	url := *r.URL
	uri := r.RequestURI

	// Requests using the CONNECT method over HTTP/2.0 must use
	// the authority field (aka r.Host) to identify the target.
	// Refer: https://httpwg.github.io/specs/rfc7540.html#CONNECT
	if r.ProtoMajor == 2 && r.Method == "CONNECT" {
		uri = r.Host
	}

	if uri == "" {
		uri = url.RequestURI()
	}

	keyvals := []interface{}{
		"method", r.Method,
		"status", code,
		"proto", r.Proto,
		"host", host,
		"user_agent", r.UserAgent(),
		"path", uri,
	}

	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		keyvals = append(keyvals, "x_forwarded_for", fwd)
	}

	if referer := r.Referer(); referer != "" {
		keyvals = append(keyvals, "referer", referer)
	}
	if code >= 500 {
		level.Info(logutil.FromContext(ctx)).Log(keyvals...)
	} else {
		level.Debug(logutil.FromContext(ctx)).Log(keyvals...)
	}
}

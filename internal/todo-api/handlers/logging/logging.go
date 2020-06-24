package logging

import (
	"net/http"
	"time"

	"github.com/justinas/alice"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

func NewHandler(logger zerolog.Logger) http.Handler {
	c := alice.New()
	c = c.Append(hlog.NewHandler(logger))
	c = c.Append(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		hlog.FromRequest(r).Info().
			Str("verb", r.Method).
			Stringer("url", r.URL).
			Int("status", status).
			Int("size", size).
			Int64("µs", duration.Microseconds()).
			Send()
	}))
	c = c.Append(hlog.RemoteAddrHandler("ip"))
	c = c.Append(hlog.UserAgentHandler("agent"))
	c = c.Append(hlog.RefererHandler("referer"))
	c = c.Append(hlog.RequestIDHandler("req_id", "Request-Id"))

	return c.ThenFunc(func(_ http.ResponseWriter, _ *http.Request) {})
}

package logging

import (
	"net/http"
	"time"

	"github.com/justinas/alice"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

func NewHandlerFunc(logger zerolog.Logger) func(http.Handler) http.Handler {
	c := alice.New()
	c = c.Append(hlog.NewHandler(logger))
	c = c.Append(hlog.RemoteAddrHandler("ip"))
	c = c.Append(hlog.UserAgentHandler("agent"))
	c = c.Append(hlog.RefererHandler("referer"))
	c = c.Append(hlog.RequestIDHandler("req_id", "Request-Id"))
	c = c.Append(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		hlog.FromRequest(r).Info().
			Str("verb", r.Method).
			Stringer("url", r.URL).
			Int("size", size).
			Int("status", status).
			Int64("duration", duration.Milliseconds()).
			Msg("HTTP Request")
	}))

	return c.Then
}

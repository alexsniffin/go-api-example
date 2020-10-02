package utils

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

func GetSubLoggerCtx(logger zerolog.Logger, ctx context.Context) context.Context {
	subLogger := logger
	reqId, ok := hlog.IDFromCtx(ctx)
	if ok {
		subLogger = subLogger.With().Str("reqID", reqId.String()).Logger()
	}
	id, ok := ctx.Value("id").(int)
	if ok {
		subLogger = subLogger.With().Int("id", id).Logger()
	}
	return subLogger.WithContext(ctx)
}

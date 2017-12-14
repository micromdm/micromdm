package connect

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
)

type Middleware func(Service) Service

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func (mw loggingMiddleware) Acknowledge(ctx context.Context, req MDMConnectRequest) (payload []byte, err error) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "Acknowledge",
			"udid", req.MDMResponse.UDID,
			"command_uuid", req.MDMResponse.CommandUUID,
			"status", req.MDMResponse.Status,
			"request_type", req.MDMResponse.RequestType,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	payload, err = mw.next.Acknowledge(ctx, req)
	return
}

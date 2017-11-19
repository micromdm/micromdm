package connect

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
)

type loggingMiddleware struct {
	logger log.Logger
	next   ConnectService
}

func NewLoggingService(svc ConnectService, logger log.Logger) loggingMiddleware {
	return loggingMiddleware{
		next:   svc,
		logger: logger,
	}
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

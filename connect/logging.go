package connect

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/micromdm/mdm"
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

func (mw loggingMiddleware) Acknowledge(ctx context.Context, req mdm.Response) (payload []byte, err error) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "Acknowledge",
			"udid", req.UDID,
			"command_uuid", req.CommandUUID,
			"status", req.Status,
			"request_type", req.RequestType,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	payload, err = mw.next.Acknowledge(ctx, req)
	return
}

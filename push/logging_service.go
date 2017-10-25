package push

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
)

type loggingMiddleware struct {
	logger log.Logger
	next   Service
}

func NewLoggingService(svc Service, logger log.Logger) loggingMiddleware {
	return loggingMiddleware{
		next:   svc,
		logger: logger,
	}
}

func (mw loggingMiddleware) Push(ctx context.Context, udid string) (id string, err error) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "Push",
			"udid", udid,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	id, err = mw.next.Push(ctx, udid)
	return
}

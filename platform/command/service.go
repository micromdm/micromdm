// Package command provides utilities for creating MDM Payloads.
package command

import (
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/micromdm/mdm"
	"golang.org/x/net/context"
)

type Service interface {
	NewCommand(context.Context, *mdm.CommandRequest) (*mdm.Payload, error)
}

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(Service) Service

// ServiceLoggingMiddleware returns a service middleware that logs the
// parameters and result of each method invocation.
func ServiceLoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return serviceLoggingMiddleware{
			logger: logger,
			next:   next,
		}
	}
}

func (mw serviceLoggingMiddleware) NewCommand(ctx context.Context, req *mdm.CommandRequest) (p *mdm.Payload, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "NewCommand",
			"error", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.next.NewCommand(ctx, req)
}

type serviceLoggingMiddleware struct {
	logger log.Logger
	next   Service
}

// ServiceInstrumentingMiddleware returns a service middleware that tracks the
// number of payloads created by the service.
func ServiceInstrumentingMiddleware(p metrics.Counter) Middleware {
	return func(next Service) Service {
		return serviceInstrumentingMiddleware{
			payloads: p,
			next:     next,
		}
	}
}

type serviceInstrumentingMiddleware struct {
	payloads metrics.Counter
	next     Service
}

func (mw serviceInstrumentingMiddleware) NewCommand(ctx context.Context, req *mdm.CommandRequest) (*mdm.Payload, error) {
	p, err := mw.next.NewCommand(ctx, req)
	mw.payloads.Add(1)
	return p, err
}

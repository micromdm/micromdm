package command

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/micromdm/micromdm/mdm/mdm"
)

var errEmptyRequest = errors.New("request must contain UDID of the device")

type Endpoints struct {
	NewCommandEndpoint endpoint.Endpoint
}

// MakeNewCommandEndpoint creates an endpoint which creates new MDM Commands.
func MakeNewCommandEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(newCommandRequest)
		if req.UDID == "" || req.RequestType == "" {
			return newCommandResponse{Err: errEmptyRequest}, nil
		}
		payload, err := svc.NewCommand(ctx, &req.CommandRequest)
		if err != nil {
			return newCommandResponse{Err: err}, nil
		}
		return newCommandResponse{Payload: payload}, nil
	}
}

// EndpointInstrumentingMiddleware returns an endpoint middleware that records
// the duration of each invocation to the passed histogram. The middleware adds
// a single field: "success", which is "true" if no error is returned, and
// "false" otherwise.
func EndpointInstrumentingMiddleware(duration metrics.Histogram) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time) {
				duration.With("success", fmt.Sprint(err == nil)).Observe(time.Since(begin).Seconds())
			}(time.Now())
			return next(ctx, request)

		}
	}
}

// EndpointLoggingMiddleware returns an endpoint middleware that logs the
// duration of each invocation, and the resulting error, if any.
func EndpointLoggingMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {

			defer func(begin time.Time) {
				logger.Log("error", err, "took", time.Since(begin))
			}(time.Now())
			return next(ctx, request)

		}
	}
}

type newCommandRequest struct {
	mdm.CommandRequest
}

type newCommandResponse struct {
	Payload *mdm.CommandPayload `json:"payload,omitempty"`
	Err     error               `json:"error,omitempty"`
}

func (r newCommandResponse) error() error { return r.Err }
func (r newCommandResponse) status() int  { return http.StatusCreated }

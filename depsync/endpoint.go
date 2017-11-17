package depsync

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	SyncNowEndpoint endpoint.Endpoint
}

func MakeSyncNowEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		s.SyncNow(ctx)
		return syncNowResponse{}, nil
	}
}

type syncNowResponse struct{}

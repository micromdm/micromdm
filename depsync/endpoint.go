package depsync

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	RefreshEndpoint endpoint.Endpoint
}

func MakeRefreshEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		s.Refresh(ctx)
		return refreshResponse{}, nil
	}
}

type refreshResponse struct{}

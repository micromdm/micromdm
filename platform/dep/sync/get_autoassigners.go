package sync

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/pkg/errors"

	"github.com/micromdm/micromdm/pkg/httputil"
)

func (s DEPSyncService) GetAutoAssigners(ctx context.Context) ([]AutoAssigner, error) {
	aa, err := s.db.LoadAutoAssigners(ctx)
	return aa, errors.Wrap(err, "get AutoAssigners")
}

type getAutoAssignersResponse struct {
	AutoAssigners []AutoAssigner `json:"autoassigners"`
	Err           error          `json:"err,omitempty"`
}

func (r getAutoAssignersResponse) Failed() error { return r.Err }

func MakeGetAutoAssignersEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		assigners, err := s.GetAutoAssigners(ctx)
		return &getAutoAssignersResponse{
			AutoAssigners: assigners,
			Err:           err,
		}, nil
	}
}

func decodeGetAutoAssignersResponse(ctx context.Context, r *http.Response) (interface{}, error) {
	var req getAutoAssignersResponse
	err := httputil.DecodeJSONResponse(r, &req)
	return req, err
}

func (e Endpoints) GetAutoAssigners(ctx context.Context) ([]AutoAssigner, error) {
	resp, err := e.GetAutoAssignersEndpoint(ctx, nil)
	if err != nil {
		return nil, err
	}
	response := resp.(getAutoAssignersResponse)
	return response.AutoAssigners, response.Err
}

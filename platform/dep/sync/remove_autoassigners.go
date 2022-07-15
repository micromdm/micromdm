package sync

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/micromdm/micromdm/pkg/httputil"
	"github.com/pkg/errors"
)

func (s DEPSyncService) RemoveAutoAssigner(ctx context.Context, filter string) error {
	err := s.db.DeleteAutoAssigner(filter)
	return errors.Wrap(err, "remove AutoAssigner")
}

type removeAutoAssignerRequest struct {
	Filter string `json:"filter"`
}

type removeAutoAssignerResponse struct {
	Err error `json:"err,omitempty"`
}

func (r removeAutoAssignerResponse) Failed() error { return r.Err }

func MakeRemoveAutoAssignerEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(removeAutoAssignerRequest)
		err := s.RemoveAutoAssigner(ctx, req.Filter)
		return &removeAutoAssignerResponse{
			Err: err,
		}, nil
	}
}

func decodeRemoveAutoAssignerResponse(ctx context.Context, r *http.Response) (interface{}, error) {
	var req removeAutoAssignerResponse
	err := httputil.DecodeJSONResponse(r, &req)
	return req, err
}

func decodeRemoveAutoAssignerRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req removeAutoAssignerRequest
	err := httputil.DecodeJSONRequest(r, &req)
	return req, err
}

func (e Endpoints) RemoveAutoAssigner(ctx context.Context, filter string) error {
	request := removeAutoAssignerRequest{Filter: filter}
	resp, err := e.RemoveAutoAssignerEndpoint(ctx, request)
	if err != nil {
		return err
	}
	response := resp.(removeAutoAssignerResponse)
	return response.Err
}

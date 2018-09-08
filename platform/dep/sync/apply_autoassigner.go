package sync

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"

	"github.com/micromdm/micromdm/pkg/httputil"
)

func (s DEPSyncService) ApplyAutoAssigner(ctx context.Context, aa *AutoAssigner) error {
	conf := s.syncer.GetConfig()
	return conf.saveAutoAssigner(aa)
}

type applyAutoAssignerRequest struct {
	*AutoAssigner
}
type applyAutoAssignerResponse struct {
	Err error `json:"err,omitempty"`
}

func (r applyAutoAssignerResponse) Failed() error { return r.Err }

func MakeApplyAutoAssignerEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(applyAutoAssignerRequest)
		err := s.ApplyAutoAssigner(ctx, req.AutoAssigner)
		return &applyAutoAssignerResponse{
			Err: err,
		}, nil
	}
}

func decodeApplyAutoAssignerResponse(ctx context.Context, r *http.Response) (interface{}, error) {
	var req applyAutoAssignerResponse
	err := httputil.DecodeJSONResponse(r, &req)
	return req, err
}

func decodeApplyAutoAssignerRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req applyAutoAssignerRequest
	err := httputil.DecodeJSONRequest(r, &req)
	return req, err
}

func (e Endpoints) ApplyAutoAssigner(ctx context.Context, aa *AutoAssigner) error {
	request := applyAutoAssignerRequest{AutoAssigner: aa}
	resp, err := e.ApplyAutoAssignerEndpoint(ctx, request)
	if err != nil {
		return err
	}
	response := resp.(applyAutoAssignerResponse)
	return response.Err
}

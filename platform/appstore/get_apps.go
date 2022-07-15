package appstore

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/groob/plist"
	"github.com/micromdm/micromdm/pkg/httputil"
	"github.com/pkg/errors"
)

type ListAppsOption struct {
	FilterName []string `json:"filter_name"`
}

type AppDTO struct {
	Name    string `json:"name"`
	Payload []byte `json:"payload,omitempty"`
}

func (svc *AppService) ListApplications(ctx context.Context, opts ListAppsOption) ([]AppDTO, error) {
	var filter string
	if len(opts.FilterName) == 1 {
		filter = opts.FilterName[0]
	}
	apps, err := svc.store.Apps(filter)
	if err != nil {
		return nil, err
	}
	var appList []AppDTO
	for name, app := range apps {
		payload, err := plist.MarshalIndent(&app, "  ")
		if err != nil {
			return nil, errors.Wrap(err, "create dto payload")
		}
		appList = append(appList, AppDTO{
			Name:    name,
			Payload: payload,
		})
	}
	return appList, nil
}

type appListRequest struct {
	Opts ListAppsOption
}

type appListResponse struct {
	Apps []AppDTO `json:"apps,omitempty"`
	Err  error    `json:"err,omitempty"`
}

func (r appListResponse) Failed() error { return r.Err }

func decodeListAppsRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req appListRequest
	err := httputil.DecodeJSONRequest(r, &req)
	return req, err
}

func decodeListAppsResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var resp appListResponse
	err := httputil.DecodeJSONResponse(r, &resp)
	return resp, err
}

func MakeListAppsEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(appListRequest)
		apps, err := svc.ListApplications(ctx, req.Opts)
		return appListResponse{
			Apps: apps,
			Err:  err,
		}, nil
	}
}

func (e Endpoints) ListApplications(ctx context.Context, opts ListAppsOption) ([]AppDTO, error) {
	request := appListRequest{opts}
	response, err := e.ListAppsEndpoint(ctx, request.Opts)
	if err != nil {
		return nil, err
	}
	return response.(appListResponse).Apps, response.(appListResponse).Err
}

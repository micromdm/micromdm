package remove

import (
	"context"
	"net/http"
	"net/url"

	"github.com/go-kit/kit/endpoint"
	"github.com/gorilla/mux"
	"github.com/micromdm/micromdm/pkg/httputil"
	"github.com/pkg/errors"
)

func (svc *RemoveService) BlockDevice(ctx context.Context, udid string) error {
	return svc.store.Save(&Device{UDID: udid})
}

type blockDeviceRequest struct {
	UDID string
}

type blockDeviceResponse struct {
	Err error `json:"err,omitempty"`
}

func (r blockDeviceResponse) Failed() error { return r.Err }

func decodeBlockDeviceRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var errBadRoute = errors.New("bad route")
	var req blockDeviceRequest
	vars := mux.Vars(r)
	udid, ok := vars["udid"]
	if !ok {
		return 0, errBadRoute
	}
	req.UDID = udid
	return req, nil
}

func encodeBlockDeviceRequest(_ context.Context, r *http.Request, request interface{}) error {
	req := request.(blockDeviceRequest)
	udid := url.QueryEscape(req.UDID)
	r.Method, r.URL.Path = "POST", "/v1/devices/"+udid+"/block"
	return nil
}

func decodeBlockDeviceResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var resp blockDeviceResponse
	err := httputil.DecodeJSONResponse(r, &resp)
	return resp, err
}

func MakeBlockDeviceEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(blockDeviceRequest)
		err = svc.BlockDevice(ctx, req.UDID)
		return &blockDeviceResponse{
			Err: err,
		}, nil
	}
}

func (e Endpoints) BlockDevice(ctx context.Context, udid string) error {
	request := blockDeviceRequest{
		UDID: udid,
	}
	resp, err := e.BlockDeviceEndpoint(ctx, request)
	if err != nil {
		return err
	}
	return resp.(blockDeviceResponse).Err
}

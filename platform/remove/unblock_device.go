package remove

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"github.com/go-kit/kit/endpoint"
	"github.com/gorilla/mux"
	"github.com/micromdm/micromdm/pkg/httputil"
)

func (svc *RemoveService) UnblockDevice(ctx context.Context, udid string) error {
	return svc.store.Delete(udid)
}

type unblockDeviceRequest struct {
	UDID string
}

type unblockDeviceResponse struct {
	Err error `json:"err,omitempty"`
}

func (r unblockDeviceResponse) Failed() error { return r.Err }

func decodeUnblockDeviceRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var errBadRoute = errors.New("bad route")
	var req unblockDeviceRequest
	vars := mux.Vars(r)
	udid, ok := vars["udid"]
	if !ok {
		return 0, errBadRoute
	}
	req.UDID = udid
	return req, nil
}

func encodeUnblockDeviceRequest(_ context.Context, r *http.Request, request interface{}) error {
	req := request.(unblockDeviceRequest)
	udid := url.QueryEscape(req.UDID)
	r.Method, r.URL.Path = "POST", "/v1/devices/"+udid+"/unblock"
	return nil
}

func decodeUnblockDeviceResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var resp unblockDeviceResponse
	err := httputil.DecodeJSONResponse(r, &resp)
	return resp, err
}

func MakeUnblockDeviceEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(unblockDeviceRequest)
		err = svc.UnblockDevice(ctx, req.UDID)
		return unblockDeviceResponse{
			Err: err,
		}, nil
	}
}

func (e Endpoints) UnblockDevice(ctx context.Context, udid string) error {
	request := unblockDeviceRequest{UDID: udid}
	resp, err := e.UnblockDeviceEndpoint(ctx, request)
	if err != nil {
		return err
	}
	return resp.(unblockDeviceResponse).Err
}

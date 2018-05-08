package device

import (
	"context"
	"net/http"
	"time"

	"github.com/go-kit/kit/endpoint"

	"github.com/micromdm/micromdm/pkg/httputil"
)

type ListDevicesOption struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`

	FilterSerial []string `json:"filter_serial"`
	FilterUDID   []string `json:"filter_udid"`
}

type DeviceDTO struct {
	SerialNumber     string    `json:"serial_number"`
	UDID             string    `json:"udid"`
	EnrollmentStatus bool      `json:"enrollment_status"`
	LastSeen         time.Time `json:"last_seen"`
}

func (svc *DeviceService) ListDevices(ctx context.Context, opt ListDevicesOption) ([]DeviceDTO, error) {
	devices, err := svc.store.List(opt)
	var dto []DeviceDTO
	for _, d := range devices {
		dto = append(dto, DeviceDTO{
			SerialNumber:     d.SerialNumber,
			UDID:             d.UDID,
			EnrollmentStatus: d.Enrolled,
			LastSeen:         d.LastSeen,
		})
	}
	return dto, err
}

type getDevicesRequest struct{ Opts ListDevicesOption }
type getDevicesResponse struct {
	Devices []DeviceDTO `json:"devices"`
	Err     error       `json:"err,omitempty"`
}

func (r getDevicesResponse) Failed() error { return r.Err }

func decodeListDevicesRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var opts ListDevicesOption
	err := httputil.DecodeJSONRequest(r, &opts)
	req := getDevicesRequest{
		Opts: opts,
	}
	return req, err
}

func decodeListDevicesResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var resp getDevicesResponse
	err := httputil.DecodeJSONResponse(r, &resp)
	return resp, err
}

func MakeListDevicesEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getDevicesRequest)
		dto, err := svc.ListDevices(ctx, req.Opts)
		return getDevicesResponse{
			Devices: dto,
			Err:     err,
		}, nil
	}
}

func (e Endpoints) ListDevices(ctx context.Context, opts ListDevicesOption) ([]DeviceDTO, error) {
	request := getDevicesRequest{opts}
	response, err := e.ListDevicesEndpoint(ctx, request.Opts)
	if err != nil {
		return nil, err
	}
	return response.(getDevicesResponse).Devices, response.(getDevicesResponse).Err
}

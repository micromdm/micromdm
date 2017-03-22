package list

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	ListDevicesEndpoint endpoint.Endpoint
}

func (e Endpoints) ListDevices(ctx context.Context, opts ListDevicesOption) ([]DeviceDTO, error) {
	request := devicesRequest{opts}
	response, err := e.ListDevicesEndpoint(ctx, request.Opts)
	if err != nil {
		return nil, err
	}
	return response.(devicesResponse).Devices, response.(devicesResponse).Err
}

func MakeListDevicesEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(devicesRequest)
		dto, err := svc.ListDevices(ctx, req.Opts)
		return devicesResponse{
			Devices: dto,
			Err:     err,
		}, nil
	}
}

type DeviceDTO struct {
	SerialNumber     string    `json:"serial_number"`
	UDID             string    `json:"udid"`
	EnrollmentStatus bool      `json:"enrollment_status"`
	LastSeen         time.Time `json:"last_seen"`
}

type devicesRequest struct{ Opts ListDevicesOption }
type devicesResponse struct {
	Devices []DeviceDTO `json:"devices"`
	Err     error       `json:"err,omitempty"`
}

package list

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/micromdm/dep"
)

type Endpoints struct {
	GetDEPAccountInfoEndpoint endpoint.Endpoint
	GetDEPDeviceEndpoint      endpoint.Endpoint
	GetDEPProfileEndpoint     endpoint.Endpoint
}

func (e Endpoints) GetDEPAccountInfo(ctx context.Context) (*dep.Account, error) {
	request := depAccountInforequest{}
	response, err := e.GetDEPAccountInfoEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.(depAccountInfoResponse).Account, response.(depAccountInfoResponse).Err
}

func (e Endpoints) GetDEPDevice(ctx context.Context, serials []string) (*dep.DeviceDetailsResponse, error) {
	request := depDeviceDetailsRequest{Serials: serials}
	response, err := e.GetDEPDeviceEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.(depDeviceDetailsResponse).DeviceDetailsResponse, response.(depDeviceDetailsResponse).Err
}

func (e Endpoints) GetDEPProfile(ctx context.Context, uuid string) (*dep.Profile, error) {
	request := depProfileRequest{UUID: uuid}
	response, err := e.GetDEPProfileEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.(depProfileResponse).Profile, response.(depProfileResponse).Err
}

func MakeGetDEPAccountInfoEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		account, err := svc.GetDEPAccountInfo(ctx)
		return depAccountInfoResponse{Account: account, Err: err}, nil
	}
}

func MakeGetDEPDeviceDetailsEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(depDeviceDetailsRequest)
		details, err := svc.GetDEPDevice(ctx, req.Serials)
		return depDeviceDetailsResponse{DeviceDetailsResponse: details, Err: err}, nil
	}
}

func MakeGetDEPProfileEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(depProfileRequest)
		profile, err := svc.GetDEPProfile(ctx, req.UUID)
		return depProfileResponse{Profile: profile, Err: err}, nil
	}
}

type depAccountInforequest struct{}
type depAccountInfoResponse struct {
	*dep.Account
	Err error `json:"err,omitempty"`
}

func (r depAccountInfoResponse) error() error { return r.Err }

type depDeviceDetailsRequest struct {
	Serials []string `json:"serials"`
}

type depDeviceDetailsResponse struct {
	*dep.DeviceDetailsResponse
	Err error `json:"err,omitempty"`
}

func (r depDeviceDetailsResponse) error() error { return r.Err }

type depProfileRequest struct {
	UUID string `json:"uuid"`
}
type depProfileResponse struct {
	*dep.Profile
	Err error `json:"err,omitempty"`
}

func (r depProfileResponse) error() error { return r.Err }

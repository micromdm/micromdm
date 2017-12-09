package list

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/micromdm/dep"

	"github.com/micromdm/micromdm/platform/deptoken"
	"github.com/micromdm/micromdm/platform/user"
)

type Endpoints struct {
	ListDevicesEndpoint       endpoint.Endpoint
	GetDEPTokensEndpoint      endpoint.Endpoint
	GetDEPAccountInfoEndpoint endpoint.Endpoint
	GetDEPDeviceEndpoint      endpoint.Endpoint
	GetDEPProfileEndpoint     endpoint.Endpoint
	ListAppsEndpont           endpoint.Endpoint
	ListUserEndpoint          endpoint.Endpoint
}

func (e Endpoints) ListUsers(ctx context.Context, opts ListUsersOption) ([]user.User, error) {
	request := userRequest{opts}
	response, err := e.ListUserEndpoint(ctx, request.Opts)
	if err != nil {
		return nil, err
	}
	return response.(userResponse).Users, response.(userResponse).Err
}

func (e Endpoints) ListDevices(ctx context.Context, opts ListDevicesOption) ([]DeviceDTO, error) {
	request := devicesRequest{opts}
	response, err := e.ListDevicesEndpoint(ctx, request.Opts)
	if err != nil {
		return nil, err
	}
	return response.(devicesResponse).Devices, response.(devicesResponse).Err
}

func (e Endpoints) ListApplications(ctx context.Context, opts ListAppsOption) ([]AppDTO, error) {
	request := appListRequest{opts}
	response, err := e.ListAppsEndpont(ctx, request.Opts)
	if err != nil {
		return nil, err
	}
	return response.(appListResponse).Apps, response.(appListResponse).Err
}

func (e Endpoints) GetDEPTokens(ctx context.Context) ([]deptoken.DEPToken, []byte, error) {
	resp, err := e.GetDEPTokensEndpoint(ctx, nil)
	if err != nil {
		return nil, nil, err
	}
	return resp.(depTokenResponse).DEPTokens, resp.(depTokenResponse).DEPPubKey, nil
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

func MakeListUsersEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(userRequest)
		dto, err := svc.ListUsers(ctx, req.Opts)
		return userResponse{
			Users: dto,
			Err:   err,
		}, nil
	}
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

func (e Endpoints) GetDEPProfile(ctx context.Context, uuid string) (*dep.Profile, error) {
	request := depProfileRequest{UUID: uuid}
	response, err := e.GetDEPProfileEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.(depProfileResponse).Profile, response.(depProfileResponse).Err
}

func MakeGetDEPTokensEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		tokens, pubkey, err := svc.GetDEPTokens(ctx)
		return depTokenResponse{
			DEPTokens: tokens,
			DEPPubKey: pubkey,
			Err:       err,
		}, nil
	}
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

type DeviceDTO struct {
	SerialNumber     string    `json:"serial_number"`
	UDID             string    `json:"udid"`
	EnrollmentStatus bool      `json:"enrollment_status"`
	LastSeen         time.Time `json:"last_seen"`
}

type userRequest struct{ Opts ListUsersOption }
type userResponse struct {
	Users []user.User `json:"users"`
	Err   error       `json:"err,omitempty"`
}

func (r userResponse) error() error { return r.Err }

type devicesRequest struct{ Opts ListDevicesOption }
type devicesResponse struct {
	Devices []DeviceDTO `json:"devices"`
	Err     error       `json:"err,omitempty"`
}

type depTokenResponse struct {
	DEPTokens []deptoken.DEPToken `json:"dep_tokens"`
	DEPPubKey []byte              `json:"public_key"`
	Err       error               `json:"err,omitempty"`
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

type appListRequest struct {
	Opts ListAppsOption
}

type AppDTO struct {
	Name    string `json:"name"`
	Payload []byte `json:"payload,omitempty"`
}

type appListResponse struct {
	Apps []AppDTO `json:"apps,omitempty"`
	Err  error    `json:"err,omitempty"`
}

func (r appListResponse) error() error { return r.Err }

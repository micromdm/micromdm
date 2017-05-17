package list

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/micromdm/micromdm/blueprint"
	"github.com/micromdm/micromdm/profile"
)

type Endpoints struct {
	ListDevicesEndpoint   endpoint.Endpoint
	GetDEPTokensEndpoint  endpoint.Endpoint
	GetBlueprintsEndpoint endpoint.Endpoint
	GetProfilesEndpoint   endpoint.Endpoint
}

func (e Endpoints) ListDevices(ctx context.Context, opts ListDevicesOption) ([]DeviceDTO, error) {
	request := devicesRequest{opts}
	response, err := e.ListDevicesEndpoint(ctx, request.Opts)
	if err != nil {
		return nil, err
	}
	return response.(devicesResponse).Devices, response.(devicesResponse).Err
}

func (e Endpoints) GetDEPTokens(ctx context.Context) ([]DEPToken, []byte, error) {
	resp, err := e.GetDEPTokensEndpoint(ctx, nil)
	if err != nil {
		return nil, nil, err
	}
	return resp.(depTokenResponse).DEPTokens, resp.(depTokenResponse).DEPPubKey, nil
}

func (e Endpoints) GetBlueprints(ctx context.Context, opt GetBlueprintsOption) ([]blueprint.Blueprint, error) {
	request := blueprintsRequest{opt}
	response, err := e.GetBlueprintsEndpoint(ctx, request.Opts)
	if err != nil {
		return nil, err
	}
	return response.(blueprintsResponse).Blueprints, response.(blueprintsResponse).Err
}

func (e Endpoints) GetProfiles(ctx context.Context, opt GetProfilesOption) ([]profile.Profile, error) {
	request := profilesRequest{opt}
	response, err := e.GetProfilesEndpoint(ctx, request.Opts)
	if err != nil {
		return nil, err
	}
	return response.(profilesResponse).Profiles, response.(profilesResponse).Err
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

func MakeGetBlueprintsEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(blueprintsRequest)
		blueprints, err := svc.GetBlueprints(ctx, req.Opts)
		return blueprintsResponse{
			Blueprints: blueprints,
			Err:        err,
		}, nil
	}
}

func MakeGetProfilesEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(profilesRequest)
		profiles, err := svc.GetProfiles(ctx, req.Opts)
		return profilesResponse{
			Profiles: profiles,
			Err:      err,
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

type DEPToken struct {
	ConsumerKey       string    `json:"consumer_key"`
	ConsumerSecret    string    `json:"consumer_secret"`
	AccessToken       string    `json:"access_token"`
	AccessSecret      string    `json:"access_secret"`
	AccessTokenExpiry time.Time `json:"access_token_expiry"`
}

type depTokenResponse struct {
	DEPTokens []DEPToken `json:"dep_tokens"`
	DEPPubKey []byte     `json:"public_key"`
	Err       error      `json:"err,omitempty"`
}

type blueprintsRequest struct{ Opts GetBlueprintsOption }
type blueprintsResponse struct {
	Blueprints []blueprint.Blueprint `json:"blueprints"`
	Err        error                 `json:"err,omitempty"`
}

func (r blueprintsResponse) error() error { return r.Err }

type profilesRequest struct{ Opts GetProfilesOption }
type profilesResponse struct {
	Profiles []profile.Profile `json:"profiles"`
	Err      error             `json:"err,omitempty"`
}

func (r profilesResponse) error() error { return r.Err }

package list

import (
	"context"
	"errors"

	"github.com/micromdm/dep"
)

type DEPService interface {
	GetDEPAccountInfo(ctx context.Context) (*dep.Account, error)
	GetDEPDevice(ctx context.Context, serials []string) (*dep.DeviceDetailsResponse, error)
	GetDEPProfile(ctx context.Context, uuid string) (*dep.Profile, error)
}

func (svc *ListService) GetDEPAccountInfo(ctx context.Context) (*dep.Account, error) {
	if svc.DEPClient == nil {
		return nil, errors.New("DEP not configured yet. add a DEP token to enable DEP")
	}
	return svc.DEPClient.Account()
}

func (svc *ListService) GetDEPDevice(ctx context.Context, serials []string) (*dep.DeviceDetailsResponse, error) {
	if svc.DEPClient == nil {
		return nil, errors.New("DEP not configured yet. add a DEP token to enable DEP")
	}
	return svc.DEPClient.DeviceDetails(serials)
}

func (svc *ListService) GetDEPProfile(ctx context.Context, uuid string) (*dep.Profile, error) {
	if svc.DEPClient == nil {
		return nil, errors.New("DEP not configured yet. add a DEP token to enable DEP")
	}
	return svc.DEPClient.FetchProfile(uuid)
}

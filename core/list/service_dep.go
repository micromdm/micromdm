package list

import (
	"context"

	"github.com/micromdm/dep"
)

type DEPService interface {
	GetDEPAccountInfo(ctx context.Context) (*dep.Account, error)
	GetDEPDevice(ctx context.Context, serials []string) (*dep.DeviceDetailsResponse, error)
	GetDEPProfile(ctx context.Context, uuid string) (*dep.Profile, error)
}

func (svc *ListService) GetDEPAccountInfo(ctx context.Context) (*dep.Account, error) {
	return svc.DEPClient.Account()
}

func (svc *ListService) GetDEPDevice(ctx context.Context, serials []string) (*dep.DeviceDetailsResponse, error) {
	return svc.DEPClient.DeviceDetails(serials)
}

func (svc *ListService) GetDEPProfile(ctx context.Context, uuid string) (*dep.Profile, error) {
	return svc.DEPClient.FetchProfile(uuid)
}

package dep

import (
	"context"
	"sync"

	"github.com/micromdm/micromdm/dep"
	"github.com/micromdm/micromdm/platform/pubsub"
)

type Service interface {
	DefineProfile(ctx context.Context, p *dep.Profile) (*dep.ProfileResponse, error)
	AssignProfile(ctx context.Context, uuid string, serials ...string) (*dep.ProfileResponse, error)
	GetAccountInfo(ctx context.Context) (*dep.Account, error)
	GetDeviceDetails(ctx context.Context, serials []string) (*dep.DeviceDetailsResponse, error)
	FetchProfile(ctx context.Context, uuid string) (*dep.Profile, error)
}

type DEPClient interface {
	DefineProfile(*dep.Profile) (*dep.ProfileResponse, error)
	AssignProfile(string, ...string) (*dep.ProfileResponse, error)
	FetchProfile(string) (*dep.Profile, error)
	Account() (*dep.Account, error)
	DeviceDetails(...string) (*dep.DeviceDetailsResponse, error)
}

type DEPService struct {
	mtx        sync.RWMutex
	client     DEPClient
	subscriber pubsub.Subscriber
}

func (svc *DEPService) Run() error {
	return svc.watchTokenUpdates(svc.subscriber)
}

func New(client DEPClient, subscriber pubsub.Subscriber) *DEPService {
	return &DEPService{client: client, subscriber: subscriber}
}

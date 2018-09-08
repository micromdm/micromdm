package dep

import (
	"context"
	"sync"

	"github.com/micromdm/micromdm/dep"
	"github.com/micromdm/micromdm/platform/pubsub"
)

type Service interface {
	DefineProfile(ctx context.Context, p *dep.Profile) (*dep.ProfileResponse, error)
	GetAccountInfo(ctx context.Context) (*dep.Account, error)
	GetDeviceDetails(ctx context.Context, serials []string) (*dep.DeviceDetailsResponse, error)
	FetchProfile(ctx context.Context, uuid string) (*dep.Profile, error)
}

type DEPService struct {
	mtx        sync.RWMutex
	client     dep.Client
	subscriber pubsub.Subscriber
}

func (svc *DEPService) Run() error {
	return svc.watchTokenUpdates(svc.subscriber)
}

func New(client dep.Client, subscriber pubsub.Subscriber) *DEPService {
	return &DEPService{client: client, subscriber: subscriber}
}

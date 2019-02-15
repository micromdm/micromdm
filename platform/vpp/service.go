package vpp

import (
	"context"
	"sync"

	"github.com/micromdm/micromdm/platform/pubsub"
	"github.com/micromdm/micromdm/vpp"
)

type Service interface {
	GetVPPApps(ctx context.Context, ids []string) (*vpp.VPPAppResponse, error)
	GetContentMetadata(ctx context.Context, options vpp.ContentMetadataOptions) (*vpp.ContentMetadata, error)
	GetAssetsSrv(ctx context.Context, options vpp.AssetsSrvOptions) (*vpp.AssetsSrv, error)
	GetLicensesSrv(ctx context.Context, options vpp.LicensesSrvOptions) (*vpp.LicensesSrv, error)
	GetServiceConfigSrv(ctx context.Context, options vpp.ServiceConfigSrvOptions) (*vpp.ServiceConfigSrv, error)
	ManageVPPLicensesByAdamIdSrv(ctx context.Context, options vpp.ManageVPPLicensesByAdamIdSrvOptions) (*vpp.ManageVPPLicensesByAdamIdSrv, error)
}

type VPPClient interface {
	GetVPPApps(...string) (*vpp.VPPAppResponse, error)
	GetContentMetadata(vpp.ContentMetadataOptions) (*vpp.ContentMetadata, error)
	GetAssetsSrv(vpp.AssetsSrvOptions) (*vpp.AssetsSrv, error)
	GetLicensesSrv(vpp.LicensesSrvOptions) (*vpp.LicensesSrv, error)
	GetServiceConfigSrv(vpp.ServiceConfigSrvOptions) (*vpp.ServiceConfigSrv, error)
	ManageVPPLicensesByAdamIdSrv(vpp.ManageVPPLicensesByAdamIdSrvOptions) (*vpp.ManageVPPLicensesByAdamIdSrv, error)
}

type VPPService struct {
	mtx        sync.RWMutex
	client     VPPClient
	subscriber pubsub.Subscriber
}

func (svc *VPPService) Run(serverURL string) error {
	return svc.watchTokenUpdates(svc.subscriber, serverURL)
}

func New(client VPPClient, subscriber pubsub.Subscriber) *VPPService {
	return &VPPService{client: client, subscriber: subscriber}
}

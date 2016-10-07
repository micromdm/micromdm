package management

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/micromdm/micromdm/application"
	"golang.org/x/net/context"
)

type installedAppsRequest struct {
	UUID string
}

type installedAppsResponse struct {
	applications []application.Application
	Err          error `json:"error,omitempty"`
}

func makeInstalledAppsEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(installedAppsRequest)
		apps, err := svc.InstalledApps(req.UUID)
		if err != nil {
			return installedAppsResponse{Err: err}, nil
		}
		return installedAppsResponse{applications: apps}, nil
	}
}

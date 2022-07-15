package appstore

import (
	"context"
	"io"

	"github.com/micromdm/micromdm/mdm/appmanifest"
)

type Service interface {
	UploadApp(ctx context.Context, manifestName string, manifest io.Reader, pkgName string, pkg io.Reader) error
	ListApplications(ctx context.Context, opt ListAppsOption) ([]AppDTO, error)
}

type AppService struct {
	store Store
}

type Store interface {
	SaveFile(name string, f io.Reader) error
	Manifest(name string) (*appmanifest.Manifest, error)
	Apps(name string) (map[string]appmanifest.Manifest, error)
}

func New(store Store) *AppService {
	return &AppService{store: store}
}

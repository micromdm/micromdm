package profile

import (
	"context"
)

type Service interface {
	ApplyProfile(ctx context.Context, p *Profile) error
	GetProfiles(ctx context.Context, opt GetProfilesOption) ([]Profile, error)
	RemoveProfiles(ctx context.Context, ids []string) error
}

type GetProfilesOption struct {
	Identifier string `json:"id"`
}

type Store interface {
	ProfileById(ctx context.Context, id string) (*Profile, error)
	Save(ctx context.Context, p *Profile) error
	List(ctx context.Context) ([]Profile, error)
	Delete(ctx context.Context, id string) error
}

func New(store Store) *ProfileService {
	return &ProfileService{store: store}
}

type ProfileService struct {
	store Store
}

func IsNotFound(err error) bool {
	type notFoundError interface {
		error
		NotFound() bool
	}

	_, ok := err.(notFoundError)
	return ok
}

package remove

import (
	"context"
	"github.com/micromdm/micromdm/blueprint"
	"github.com/micromdm/micromdm/profile"
)

type Service interface {
	RemoveBlueprints(ctx context.Context, names []string) error
	RemoveProfiles(ctx context.Context, ids []string) error
}

type RemoveService struct {
	Blueprints *blueprint.DB
	Profiles   *profile.DB
}

func (svc *RemoveService) RemoveBlueprints(ctx context.Context, names []string) error {
	// TODO: Wrap deletion(s) in transactions so as to not have
	// incomplete removals?
	for _, name := range names {
		err := svc.Blueprints.Delete(name)
		if err != nil {
			return err
		}
	}
	return nil
}

func (svc *RemoveService) RemoveProfiles(ctx context.Context, ids []string) error {
	// TODO: Wrap deletion(s) in transactions so as to not have
	// incomplete removals?
	for _, id := range ids {
		err := svc.Profiles.Delete(id)
		if err != nil {
			return err
		}
	}
	return nil
}

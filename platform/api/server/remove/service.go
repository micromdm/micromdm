package remove

import (
	"context"

	"github.com/micromdm/micromdm/platform/blueprint"
	"github.com/micromdm/micromdm/platform/remove"
)

type Service interface {
	RemoveBlueprints(ctx context.Context, names []string) error
	UnblockDevice(ctx context.Context, udid string) error
}

type RemoveService struct {
	Blueprints *blueprint.DB
	*remove.RemoveService
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

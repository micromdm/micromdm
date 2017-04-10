package apply

import (
	"context"
	"errors"

	"github.com/micromdm/micromdm/blueprint"
)

type Service interface {
	ApplyBlueprint(ctx context.Context, bp *blueprint.Blueprint) error
	ApplyDEPToken(ctx context.Context, P7MContent []byte) error
}

type ApplyService struct {
	Blueprints *blueprint.DB
}

func (svc *ApplyService) ApplyBlueprint(ctx context.Context, bp *blueprint.Blueprint) error {
	return svc.Blueprints.Save(bp)
}

func (svc *ApplyService) ApplyDEPToken(ctx context.Context, P7MContent []byte) error {
	return errors.New("not implemented yet")
}

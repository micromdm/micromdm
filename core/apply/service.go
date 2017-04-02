package apply

import (
	"context"

	"github.com/micromdm/micromdm/blueprint"
)

type Service interface {
	ApplyBlueprint(ctx context.Context, bp *blueprint.Blueprint) error
}

type ApplyService struct {
	Blueprints *blueprint.DB
}

func (svc *ApplyService) ApplyBlueprint(ctx context.Context, bp *blueprint.Blueprint) error {
	return svc.Blueprints.Save(bp)
}

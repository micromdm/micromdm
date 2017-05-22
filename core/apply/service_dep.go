package apply

import (
	"context"

	"github.com/micromdm/dep"
)

type DEPService interface {
	DefineDEPProfile(ctx context.Context, p *dep.Profile) (*dep.ProfileResponse, error)
}

func (svc *ApplyService) DefineDEPProfile(ctx context.Context, p *dep.Profile) (*dep.ProfileResponse, error) {
	return svc.DEPClient.DefineProfile(p)
}

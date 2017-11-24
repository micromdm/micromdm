package apply

import (
	"context"
	"errors"

	"github.com/micromdm/dep"
)

type DEPService interface {
	DefineDEPProfile(ctx context.Context, p *dep.Profile) (*dep.ProfileResponse, error)
}

func (svc *ApplyService) DefineDEPProfile(ctx context.Context, p *dep.Profile) (*dep.ProfileResponse, error) {
	if svc.DEPClient == nil {
		return nil, errors.New("DEP not configured yet. add a DEP token to enable DEP")
	}
	return svc.DEPClient.DefineProfile(p)
}

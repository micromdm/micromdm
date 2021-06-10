package challenge

import (
	"context"
	"errors"

	"github.com/micromdm/scep/v2/challenge"
)

type Service interface {
	SCEPChallenge(ctx context.Context) (string, error)
}

type ChallengeService struct {
	challenge.Store
}

func (c *ChallengeService) SCEPChallenge(_ context.Context) (string, error) {
	if c.Store == nil {
		return "", errors.New("SCEP challenge store missing")
	}
	return c.Store.SCEPChallenge()
}

func NewService(challengeStore challenge.Store) *ChallengeService {
	return &ChallengeService{
		Store: challengeStore,
	}
}

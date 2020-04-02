package challenge

import (
	"context"
	"errors"

	challengestore "github.com/micromdm/scep/challenge/bolt"
)

type Service interface {
	SCEPChallenge(ctx context.Context) (string, error)
}

type ChallengeService struct {
	scepChallengeStore *challengestore.Depot
}

func (c *ChallengeService) SCEPChallenge(ctx context.Context) (string, error) {
	if c.scepChallengeStore == nil {
		return "", errors.New("SCEP challenge store missing")
	}
	return c.scepChallengeStore.SCEPChallenge()
}

func NewService(cs *challengestore.Depot) *ChallengeService {
	return &ChallengeService{
		scepChallengeStore: cs,
	}
}

package challenge

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type challengeResponse struct {
	Challenge string `json:"string"`
	Err       error  `json:"err,omitempty"`
}

func (r challengeResponse) Failed() error { return r.Err }

func MakeChallengeEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := challengeResponse{}
		r.Challenge, r.Err = svc.SCEPChallenge(ctx)
		return r, nil
	}
}

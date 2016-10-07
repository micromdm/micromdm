package enroll

import (
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
)

type Endpoints struct {
	GetEnrollEndpoint endpoint.Endpoint
}

type mdmEnrollRequest struct{}

type mdmEnrollResponse struct {
	Profile
	Err error `plist:"error,omitempty"`
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		GetEnrollEndpoint: MakeGetEnrollEndpoint(s),
	}
}

func MakeGetEnrollEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		profile, err := s.Enroll(ctx)
		return mdmEnrollResponse{profile, err}, nil
	}
}

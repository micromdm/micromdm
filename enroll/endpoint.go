package enroll

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-kit/kit/endpoint"
	"github.com/micromdm/mdm"
)

type Endpoints struct {
	GetEnrollEndpoint endpoint.Endpoint
	OTAEnrollEndpoint endpoint.Endpoint
}

type depEnrollmentRequest struct {
	mdm.DEPEnrollmentRequest
}

type mdmEnrollRequest struct{}

type mdmEnrollResponse struct {
	Profile
	Err error `plist:"error,omitempty"`
}

type mdmOTAEnrollResponse struct {
	Payload
	Err error `plist:"error,omitempty"`
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		GetEnrollEndpoint: MakeGetEnrollEndpoint(s),
		OTAEnrollEndpoint: MakeOTAEnrollEndpoint(s),
	}
}

func MakeGetEnrollEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		switch req := request.(type) {
		case mdmEnrollRequest:
			profile, err := s.Enroll(ctx)
			return mdmEnrollResponse{profile, err}, nil
		case depEnrollmentRequest:
			fmt.Printf("got DEP enrollment request from %s\n", req.Serial)
			profile, err := s.Enroll(ctx)
			return mdmEnrollResponse{profile, err}, nil
		default:
			return nil, errors.New("unknown enrollment type")
		}
	}
}

func MakeOTAEnrollEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		payload, err := s.OTAEnroll(ctx)
		return mdmOTAEnrollResponse{payload, err}, nil
	}
}

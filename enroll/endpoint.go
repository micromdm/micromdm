package enroll

import (
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
)

type mdmEnrollRequest struct{}

type mdmEnrollResponse struct {
	*Profile
}

func makeEnrollEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		//req := request.(mdmEnrollRequest)
		profile, err := svc.Enroll()
		if err != nil {
			return mdmEnrollResponse{}, err
		}
		return mdmEnrollResponse{profile}
	}
}

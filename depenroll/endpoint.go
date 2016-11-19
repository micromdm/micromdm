package depenroll

import (
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"

	"github.com/micromdm/mdm"
)

type depEnrollmentRequest struct {
	mdm.DEPEnrollmentRequest
}

type depEnrollmentResponse struct {
	Profile []byte // MDM Enrollment Profile
	Err     error  `plist:"error,omitempty"`
}

func (r depEnrollmentResponse) error() error { return r.Err }

func makeDEPEnrollmentEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(depEnrollmentRequest)
		profile, err := svc.EnrollDEP(req.DEPEnrollmentRequest.UDID, req.DEPEnrollmentRequest.Serial)
		if err != nil {
			return depEnrollmentResponse{Err: err}, nil
		}
		return depEnrollmentResponse{Profile: profile}, nil
	}
}

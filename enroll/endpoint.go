package enroll

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-kit/kit/endpoint"
	"github.com/micromdm/mdm"
)

type Endpoints struct {
	GetEnrollEndpoint       endpoint.Endpoint
	OTAEnrollEndpoint       endpoint.Endpoint
	OTAPhase2Phase3Endpoint endpoint.Endpoint
}

type depEnrollmentRequest struct {
	mdm.DEPEnrollmentRequest
}

// TODO: may overlap at some point with mdm.DEPEnrollmentRequest
type otaEnrollmentRequest struct {
	Challenge     string `plist:"CHALLENGE"`
	Product       string `plist:"PRODUCT"`
	Serial        string `plist:"SERIAL"`
	UDID          string `plist:"UDID"`
	Version       string `plist:"VERSION"` // build no.
	IMSI          string `plist:"IMSI"`
	IMEI          string `plist:"IMEI,omitempty"`
	MEID          string `plist:"MEID,omitempty"`
	ICCID         string `plist:"ICCID"`
	MacAddressEn0 string `plist:"MAC_ADDRESS_EN0"`
	DeviceName    string `plist:"DEVICE_NAME"`
	NotOnConsole  bool
	UserID        string // GUID of User
	UserLongName  string
	UserShortName string
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
		GetEnrollEndpoint:       MakeGetEnrollEndpoint(s),
		OTAEnrollEndpoint:       MakeOTAEnrollEndpoint(s),
		OTAPhase2Phase3Endpoint: MakeOTAPhase2Phase3Endpoint(s),
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

func MakeOTAPhase2Phase3Endpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		profile, err := s.OTAPhase2(ctx)
		// TODO: if Phase 3
		// s.OTAPhase3(ctx)
		return mdmEnrollResponse{profile, err}, nil
	}
}

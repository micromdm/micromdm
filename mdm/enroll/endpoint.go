package enroll

import (
	"context"
	"errors"
	"fmt"

	"github.com/fullsailor/pkcs7"
	"github.com/go-kit/kit/endpoint"
	boltdepot "github.com/micromdm/scep/depot/bolt"

	"github.com/micromdm/micromdm/pkg/crypto"
	"github.com/micromdm/micromdm/platform/profile"
)

type Endpoints struct {
	GetEnrollEndpoint       endpoint.Endpoint
	OTAEnrollEndpoint       endpoint.Endpoint
	OTAPhase2Phase3Endpoint endpoint.Endpoint
}

type depEnrollmentRequest struct {
	Language string `plist:"LANGUAGE"`
	Product  string `plist:"PRODUCT"`
	Serial   string `plist:"SERIAL"`
	UDID     string `plist:"UDID"`
	Version  string `plist:"VERSION"`
	IMEI     string `plist:"IMEI,omitempty"`
	MEID     string `plist:"MEID,omitempty"`
}

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
	MACAddressEN0 string `plist:"MAC_ADDRESS_EN0"`
	DeviceName    string `plist:"DEVICE_NAME"`
	NotOnConsole  bool
	UserID        string // GUID of User
	UserLongName  string
	UserShortName string
}

type mdmEnrollRequest struct{}

type mobileconfigResponse struct {
	profile.Mobileconfig
	Err error `plist:"error,omitempty"`
}

type mdmOTAPhase2Phase3Request struct {
	otaEnrollmentRequest otaEnrollmentRequest
	p7                   *pkcs7.PKCS7
}

func MakeServerEndpoints(s Service, scepDepot *boltdepot.Depot) Endpoints {
	return Endpoints{
		GetEnrollEndpoint:       MakeGetEnrollEndpoint(s),
		OTAEnrollEndpoint:       MakeOTAEnrollEndpoint(s),
		OTAPhase2Phase3Endpoint: MakeOTAPhase2Phase3Endpoint(s, scepDepot),
	}
}

func MakeGetEnrollEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		switch req := request.(type) {
		case mdmEnrollRequest:
			mc, err := s.Enroll(ctx)
			return mobileconfigResponse{mc, err}, nil
		case depEnrollmentRequest:
			fmt.Printf("got DEP enrollment request from %s\n", req.Serial)
			mc, err := s.Enroll(ctx)
			return mobileconfigResponse{mc, err}, nil
		default:
			return nil, errors.New("unknown enrollment type")
		}
	}
}

func MakeOTAEnrollEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		mc, err := s.OTAEnroll(ctx)
		return mobileconfigResponse{mc, err}, nil
	}
}

func MakeOTAPhase2Phase3Endpoint(s Service, scepDepot *boltdepot.Depot) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(mdmOTAPhase2Phase3Request)

		if req.p7 == nil || req.p7.GetOnlySigner() == nil {
			return nil, errors.New("invalid signer/signer not provided")
		}

		// TODO: currently only verifying the signing certificate but ought to
		// verify the whole provided chain. Note this will be difficult to do
		// given the inconsist certificate chain returned by macOS in OTA mode,
		// macOS in DEP mode, and iOS in either mode. See:
		// https://openradar.appspot.com/radar?id=4957320861712384
		if err := crypto.VerifyFromAppleDeviceCA(req.p7.GetOnlySigner()); err == nil {
			// signing certificate is signed by the Apple Device CA. this means
			// we don't yet have a SCEP identity and thus are in Phase 2 of the
			// OTA enrollment
			mc, err := s.OTAPhase2(ctx)
			return mobileconfigResponse{mc, err}, nil
		}

		caChain, _, err := scepDepot.CA(nil)
		if err != nil {
			return nil, err
		}

		if len(caChain) < 1 {
			return nil, errors.New("invalid SCEP CA chain")
		}

		if req.p7.GetOnlySigner().CheckSignatureFrom(caChain[0]) == nil {
			// signing certificate is signed by our SCEP CA. this means we
			// we are in Phase 3 of OTA enrollment (as we already have a
			// identified certificate)

			// TODO: possibly deliver a different enrollment profile based
			// on device certificates
			// TODO: we can encrypt the enrollment (or any profile) at this
			// point: we have a device identity that we can encrypt to that
			// device's private key that it can decrypt
			// TODO: the SCEP CA checking ought to be more robust
			// see: https://github.com/micromdm/scep/issues/32

			mc, err := s.Enroll(ctx)
			// profile, err := s.OTAPhase3(ctx)
			return mobileconfigResponse{mc, err}, nil
		}
		return mobileconfigResponse{profile.Mobileconfig{}, errors.New("unauthorized client")}, nil
	}
}

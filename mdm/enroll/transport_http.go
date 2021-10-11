package enroll

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/micromdm/micromdm/pkg/crypto"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/groob/plist"
	"go.mozilla.org/pkcs7"
)

type HTTPHandlers struct {
	EnrollHandler    http.Handler
	OTAEnrollHandler http.Handler

	// In Apple's Over-the-Air design Phases 2 and 3 happen over the same URL.
	// The differentiator is which certificate signed the CMS POST body.
	OTAPhase2Phase3Handler http.Handler
}

func MakeHTTPHandlers(ctx context.Context, endpoints Endpoints, opts ...httptransport.ServerOption) HTTPHandlers {
	h := HTTPHandlers{
		EnrollHandler: httptransport.NewServer(
			endpoints.GetEnrollEndpoint,
			decodeMDMEnrollRequest,
			encodeMobileconfigResponse,
			opts...,
		),
		OTAEnrollHandler: httptransport.NewServer(
			endpoints.OTAEnrollEndpoint,
			decodeEmptyRequest,
			encodeMobileconfigResponse,
			opts...,
		),
		OTAPhase2Phase3Handler: httptransport.NewServer(
			endpoints.OTAPhase2Phase3Endpoint,
			decodeOTAPhase2Phase3Request,
			encodeMobileconfigResponse,
			opts...,
		),
	}
	return h
}

func decodeEmptyRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeMDMEnrollRequest(_ context.Context, r *http.Request) (interface{}, error) {
	switch r.Method {
	case "GET":
		return mdmEnrollRequest{}, nil
	case "POST": // DEP request
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		p7, err := pkcs7.Parse(data)
		if err != nil {
			return nil, err
		}
		err = p7.Verify()
		if err != nil {
			return nil, err
		}
		// TODO: for thse errors provide better feedback as 4xx HTTP status
		signer := p7.GetOnlySigner()
		if signer == nil {
			return nil, errors.New("invalid CMS signer during enrollment")
		}
		err = crypto.VerifyFromAppleDeviceCA(signer)
		if err != nil {
			return nil, errors.New("unauthorized enrollment client: not signed by Apple Device CA")
		}
		var request depEnrollmentRequest
		if err := plist.Unmarshal(p7.Content, &request); err != nil {
			return nil, err
		}
		return request, nil
	default:
		return nil, errors.New("unknown enrollment method")
	}
}

func encodeMobileconfigResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/x-apple-aspen-config")
	mcResp := response.(mobileconfigResponse)
	_, err := w.Write(mcResp.Mobileconfig)
	return err
}

func decodeOTAPhase2Phase3Request(_ context.Context, r *http.Request) (interface{}, error) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	p7, err := pkcs7.Parse(data)
	if err != nil {
		return nil, err
	}
	err = p7.Verify()
	if err != nil {
		return nil, err
	}
	var request otaEnrollmentRequest
	err = plist.Unmarshal(p7.Content, &request)
	if err != nil {
		return nil, err
	}
	return mdmOTAPhase2Phase3Request{request, p7}, nil
}

package mdm

import (
	"context"
	"io/ioutil"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type Endpoints struct {
	CheckinEndpoint     endpoint.Endpoint
	AcknowledgeEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		CheckinEndpoint:     MakeCheckinEndpoint(s),
		AcknowledgeEndpoint: MakeAcknowledgeEndpoint(s),
	}
}

func RegisterHTTPHandlers(r *mux.Router, e Endpoints, verifier SignatureVerifier, logger log.Logger) {
	decoder := &requestDecoder{verifier: verifier}
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerBefore(httptransport.PopulateRequestContext),
	}

	r.Methods(http.MethodPut).Path("/mdm/checkin").Handler(httptransport.NewServer(
		e.CheckinEndpoint,
		decoder.decodeCheckinRequest,
		encodeResponse,
		options...,
	))

	r.Methods(http.MethodPut).Path("/mdm/connect").Handler(httptransport.NewServer(
		e.AcknowledgeEndpoint,
		decoder.decodeAcknowledgeRequest,
		encodeResponse,
		options...,
	))
}

type SignatureVerifier interface {
	VerifySignature(sig string, message []byte) error
}

type requestDecoder struct {
	verifier SignatureVerifier
}

func (d *requestDecoder) readBody(r *http.Request) ([]byte, error) {
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrap(err, "reading MDM Response HTTP Body")
	}

	if d.verifier != nil {
		b64sig := r.Header.Get("Mdm-Signature")
		if err := d.verifier.VerifySignature(b64sig, body); err != nil {
			return nil, errors.Wrap(err, "verify signature")
		}
	}

	return body, nil
}

// According to the MDM Check-in protocol, the server must respond with 200 OK
// to successful Check-in requests.
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	type failer interface {
		Failed() error
	}

	if e, ok := response.(failer); ok && e.Failed() != nil {
		encodeError(ctx, e.Failed(), w)
		return nil
	}

	w.WriteHeader(http.StatusOK)

	type payloader interface {
		Response() []byte
	}

	var err error
	if r, ok := response.(payloader); ok {
		_, err = w.Write(r.Response())
	}
	return errors.Wrap(err, "write acknowledge response")
}

func encodeError(ctx context.Context, err error, w http.ResponseWriter) {
	err = errors.Cause(err)
	type rejectUserAuthError interface {
		error
		UserAuthReject() bool
	}
	if e, ok := err.(rejectUserAuthError); ok && e.UserAuthReject() {
		w.WriteHeader(http.StatusGone)
		return
	}

	type checkoutErr interface {
		error
		Checkout() bool
	}
	if e, ok := err.(checkoutErr); ok && e.Checkout() {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
}

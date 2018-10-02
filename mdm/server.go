package mdm

import (
	"bytes"
	"context"
	"crypto/x509"
	"encoding/base64"
	"io/ioutil"
	"net/http"

	"github.com/fullsailor/pkcs7"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/groob/plist"
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

func RegisterHTTPHandlers(r *mux.Router, e Endpoints, logger log.Logger) {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerBefore(httptransport.PopulateRequestContext),
		httptransport.ServerBefore(populateDeviceCertificateFromSignRequestHeader),
	}

	r.Methods(http.MethodPut).Path("/mdm/checkin").Handler(httptransport.NewServer(
		e.CheckinEndpoint,
		decodeCheckinRequest,
		encodeResponse,
		options...,
	))

	r.Methods(http.MethodPut).Path("/mdm/connect").Handler(httptransport.NewServer(
		e.AcknowledgeEndpoint,
		decodeAcknowledgeRequest,
		encodeResponse,
		options...,
	))
}

type contextKey int

const (
	ContextKeyDeviceCertificate contextKey = iota
	ContextKeyDeviceCertificateVerifyError
)

func DeviceCertificateFromContext(ctx context.Context) (*x509.Certificate, error) {
	cert := ctx.Value(ContextKeyDeviceCertificate).(*x509.Certificate)
	err, _ := ctx.Value(ContextKeyDeviceCertificateVerifyError).(error)
	return cert, err
}

func populateDeviceCertificateFromSignRequestHeader(ctx context.Context, r *http.Request) context.Context {
	bodyReader := r.Body
	defer bodyReader.Close()

	// We can't gracefully bubble up errors from this function,
	// so we silently disregard them (terrible)
	body, _ := ioutil.ReadAll(r.Body)

	// Replace our body object with a fully buffered response
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	cert, err := verifySignature(r.Header.Get("Mdm-Signature"), body)
	ctx = context.WithValue(ctx, ContextKeyDeviceCertificate, cert)
	ctx = context.WithValue(ctx, ContextKeyDeviceCertificateVerifyError, err)

	return ctx
}

// TODO: If we ever use Go client cert auth we can use
//       r.TLS.PeerCertificates to return the client cert. Unnecessary
//       now as default config is uses Mdm-Signature header method instead
//       (for better compatilibity with proxies, etc.)
// func populateDeviceCertificateFromTLSPeerCertificates()

// Extract (raw) body bytes, parse property list
func mdmRequestBody(r *http.Request, s interface{}) ([]byte, error) {
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrap(err, "reading MDM acknowledge HTTP body")
	}

	err = plist.Unmarshal(body, s)
	if err != nil {
		return body, errors.Wrap(err, "unmarshal MDM acknowledge plist")
	}

	return body, nil
}

// Verify MDM header signature. Note: does NOT verify device certificate
func verifySignature(header string, body []byte) (*x509.Certificate, error) {
	if header == "" {
		return nil, errors.New("signature missing")
	}
	sig, err := base64.StdEncoding.DecodeString(header)
	if err != nil {
		return nil, errors.Wrap(err, "decode MDM SignMessage header")
	}
	p7, err := pkcs7.Parse(sig)
	if err != nil {
		return nil, errors.Wrap(err, "CMS parse decoded MDM SignMessage signature")
	}
	p7.Content = body
	if err := p7.Verify(); err != nil {
		return nil, errors.Wrap(err, "CMS verify MDM Signed Message")
	}
	cert := p7.GetOnlySigner()
	if cert == nil {
		return nil, errors.New("invalid or missing CMS signer")
	}
	return cert, nil
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

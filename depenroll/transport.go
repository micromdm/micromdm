package depenroll

import (
	"io/ioutil"
	"net/http"

	"github.com/fullsailor/pkcs7"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/groob/plist"
	"golang.org/x/net/context"
)

// ServiceHandler returns an HTTP Handler for the DEP enrollment service.
func ServiceHandler(ctx context.Context, svc Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}
	depEnrollmentHandler := kithttp.NewServer(
		ctx,
		makeDEPEnrollmentEndpoint(svc),
		decodeMDMEnrollmentRequest,
		encodeDEPEnrollmentResponse,
		opts...,
	)
	r := mux.NewRouter()
	r.Handle("/mdm/enroll/dep", depEnrollmentHandler).Methods("POST")
	return r
}

type errorer interface {
	error() error
}

// The enrollment request is PkCS7 signed.
// We'll ignore everything but the content for now
func decodeMDMEnrollmentRequest(_ context.Context, r *http.Request) (interface{}, error) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	p7, err := pkcs7.Parse(data)
	if err != nil {
		return nil, err
	}
	// TODO: We should verify but not currently possible. Apple
	// does no provide a cert for the CA.
	var request depEnrollmentRequest
	if err := plist.Unmarshal(p7.Content, &request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeDEPEnrollmentResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	resp := response.(depEnrollmentResponse)
	w.Write(resp.Profile)
	return nil
}

// encode errors from business-logic
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	// unwrap if the error is wrapped by kit http in it's own error type
	if httperr, ok := err.(kithttp.Error); ok {
		err = httperr.Err
	}
	switch err {
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	// w.Header().Set("Content-Type", "application/json; charset=utf-8")
	plist.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

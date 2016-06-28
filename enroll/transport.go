package enroll

import (
	"net/http"

	"golang.org/x/net/context"

	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/groob/plist"
)

// ServiceHandler returns an HTTP Handler for the enroll service
func ServiceHandler(ctx context.Context, svc Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
	}

	connectHandler := kithttp.NewServer(
		ctx,
		makeEnrollEndpoint(svc),
		decodeMDMEnrollRequest,
		encodeResponse,
		opts...,
	)
	r := mux.NewRouter()

	r.Handle("/mdm/enroll", connectHandler).Methods("GET")
	return r
}

func decodeMDMEnrollRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return r, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(mdmEnrollResponse)

	plistData, err := plist.Marshal(resp.Profile)
	if err != nil {
		return err
	}

	if len(plistData) != 0 {
		w.Write(plistData)
	}
	return nil
}

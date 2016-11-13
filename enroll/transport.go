package enroll

import (
	"net/http"

	"golang.org/x/net/context"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/groob/plist"
)

// ServiceHandler returns an HTTP Handler for the enroll service
func ServiceHandler(ctx context.Context, svc Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoints(svc)
	opts := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
	}

	r.Methods("GET").Path("/mdm/enroll").Handler(httptransport.NewServer(
		ctx,
		e.GetEnrollEndpoint,
		decodeMDMEnrollRequest,
		encodeResponse,
		opts...,
	))

	return r
}

func decodeMDMEnrollRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return r, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(mdmEnrollResponse)

	w.Header().Set("Content-Type", "application/x-apple-aspen-config")

	if err := plist.NewEncoder(w).Encode(resp); err != nil {
		return err
	}

	return nil
}

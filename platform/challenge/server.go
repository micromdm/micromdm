package challenge

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/micromdm/micromdm/pkg/httputil"
)

type Endpoints struct {
	ChallengeEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(s Service, outer endpoint.Middleware, others ...endpoint.Middleware) Endpoints {
	return Endpoints{
		ChallengeEndpoint: endpoint.Chain(outer, others...)(MakeChallengeEndpoint(s)),
	}
}

func RegisterHTTPHandlers(r *mux.Router, e Endpoints, options ...httptransport.ServerOption) {
	// POST   /v1/challenge    Generate and return SCEP challenge

	r.Methods("POST").Path("/v1/challenge").Handler(httptransport.NewServer(
		e.ChallengeEndpoint,
		decodeEmptyRequest,
		httputil.EncodeJSONResponse,
		options...,
	))
}

func decodeEmptyRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

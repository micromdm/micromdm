package depsync

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/micromdm/micromdm/pkg/httputil"
)

func NewService(syncer Syncer) *DEPSyncService {
	return &DEPSyncService{syncer: syncer}
}

type Endpoints struct {
	SyncNowEndpoint            endpoint.Endpoint
	ApplyAutoAssignerEndpoint  endpoint.Endpoint
	GetAutoAssignersEndpoint   endpoint.Endpoint
	RemoveAutoAssignerEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		SyncNowEndpoint:            MakeSyncNowEndpoint(s),
		ApplyAutoAssignerEndpoint:  MakeApplyAutoAssignerEndpoint(s),
		GetAutoAssignersEndpoint:   MakeGetAutoAssignersEndpoint(s),
		RemoveAutoAssignerEndpoint: MakeRemoveAutoAssignerEndpoint(s),
	}
}

func MakeHTTPHandler(e Endpoints, logger log.Logger) *mux.Router {
	r, options := httputil.NewRouter(logger)

	// POST		/v1/dep/syncnow			request a DEP sync operation to happen now
	// POST		/v1/dep/autoassigners	set a DEP auto-assigner
	// GET		/v1/dep/autoassigners	get list of DEP auto-assigners
	// DELETE	/v1/dep/autoassigners	remove a DEP auto-assigner

	r.Methods("POST").Path("/v1/dep/syncnow").Handler(httptransport.NewServer(
		e.SyncNowEndpoint,
		decodeEmptyRequest,
		encodeEmptyResponse,
		options...,
	))

	r.Methods("POST").Path("/v1/dep/autoassigners").Handler(httptransport.NewServer(
		e.ApplyAutoAssignerEndpoint,
		decodeApplyAutoAssignerRequest,
		httputil.EncodeJSONResponse,
		options...,
	))

	r.Methods("GET").Path("/v1/dep/autoassigners").Handler(httptransport.NewServer(
		e.GetAutoAssignersEndpoint,
		decodeEmptyRequest,
		httputil.EncodeJSONResponse,
		options...,
	))

	r.Methods("DELETE").Path("/v1/dep/autoassigners").Handler(httptransport.NewServer(
		e.RemoveAutoAssignerEndpoint,
		decodeRemoveAutoAssignerRequest,
		httputil.EncodeJSONResponse,
		options...,
	))

	return r
}

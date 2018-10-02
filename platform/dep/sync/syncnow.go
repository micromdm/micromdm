package sync

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"
)

func (s *DEPSyncService) SyncNow(_ context.Context) error {
	s.syncer.SyncNow()
	return nil
}

type syncNowResponse struct{}

func MakeSyncNowEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, _ interface{}) (interface{}, error) {
		err := s.SyncNow(ctx)
		return syncNowResponse{}, err
	}
}

func decodeEmptyRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func encodeEmptyResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return nil
}

func decodeEmptyResponse(ctx context.Context, r *http.Response) (interface{}, error) {
	return nil, nil
}

func (e Endpoints) SyncNow(ctx context.Context) error {
	_, err := e.SyncNowEndpoint(ctx, nil)
	return err
}

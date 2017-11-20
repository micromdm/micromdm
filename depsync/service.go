package depsync

import (
	"context"
)

type Service interface {
	Refresh(ctx context.Context)
}

type RefreshService struct {
	syncer Syncer
}

func (s *RefreshService) Refresh(_ context.Context) {
	s.syncer.SyncNow()
	return
}

func NewRPC(syncer Syncer) *RefreshService {
	return &RefreshService{syncer: syncer}
}

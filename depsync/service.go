package depsync

import (
	"context"
)

type Service interface {
	SyncNow(ctx context.Context)
}

type syncNowService struct {
	syncer Syncer
}

func (s *syncNowService) SyncNow(_ context.Context) {
	s.syncer.SyncNow()
	return
}

func NewRPC(syncer Syncer) *syncNowService {
	return &syncNowService{syncer: syncer}
}

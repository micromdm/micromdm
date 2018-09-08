package sync

import (
	"context"
)

type Service interface {
	SyncNow(context.Context) error
	ApplyAutoAssigner(context.Context, *AutoAssigner) error
	GetAutoAssigners(context.Context) ([]AutoAssigner, error)
	RemoveAutoAssigner(context.Context, string) error
}

type DB interface {
	SaveAutoAssigner(a *AutoAssigner) error
	LoadAutoAssigners() ([]AutoAssigner, error)
	DeleteAutoAssigner(filter string) error
}

type DEPSyncService struct {
	db     DB
	syncer Syncer
}

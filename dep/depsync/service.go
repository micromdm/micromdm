package depsync

import (
	"context"
)

type Service interface {
	SyncNow(context.Context) error
	ApplyAutoAssigner(context.Context, *AutoAssigner) error
	GetAutoAssigners(context.Context) ([]*AutoAssigner, error)
	RemoveAutoAssigner(context.Context, string) error
}

type DEPSyncService struct {
	syncer Syncer
}

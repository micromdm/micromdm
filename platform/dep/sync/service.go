package sync

import (
	"context"
	"time"
)

type Service interface {
	SyncNow(context.Context) error
	ApplyAutoAssigner(context.Context, *AutoAssigner) error
	GetAutoAssigners(context.Context) ([]AutoAssigner, error)
	RemoveAutoAssigner(context.Context, string) error
}

type DB interface {
	SaveAutoAssigner(ctx context.Context, a *AutoAssigner) error
	LoadAutoAssigners(ctx context.Context) ([]AutoAssigner, error)
	DeleteAutoAssigner(ctx context.Context, filter string) error
}

type DEPSyncService struct {
	db     DB
	syncer Syncer
}

type Cursor struct {
	Value     string    `json:"value" db:"value"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// A cursor is valid for a week.
func (c Cursor) Valid() bool {
	expiration := time.Now().Add(cursorValidDuration)
	return c.CreatedAt.Before(expiration)
}

type AutoAssigner struct {
	Filter      string `json:"filter"`
	ProfileUUID string `json:"profile_uuid"`
}

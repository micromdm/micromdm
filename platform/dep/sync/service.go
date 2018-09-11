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
	SaveAutoAssigner(a *AutoAssigner) error
	LoadAutoAssigners() ([]AutoAssigner, error)
	DeleteAutoAssigner(filter string) error
}

type DEPSyncService struct {
	db     DB
	syncer Syncer
}

type Cursor struct {
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
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

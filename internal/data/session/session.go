// Package session tracks user authenticated sessions.
package session

import (
	"context"
	"errors"
	"time"

	"micromdm.io/v2/pkg/id"
	"micromdm.io/v2/pkg/viewer"
)

type Session struct {
	ID         string
	UserID     string
	CreatedAt  time.Time
	AccessedAt time.Time
}

func columns() []string {
	return []string{
		"id",
		"user_id",
		"created_at",
		"accessed_at",
	}
}

func create(ctx context.Context) (*Session, error) {
	v, ok := viewer.FromContext(ctx)
	if !ok {
		return nil, errors.New("cannot create session without a viewer in context")
	}

	return &Session{
		ID:     id.New(),
		UserID: v.UserID,
	}, nil
}

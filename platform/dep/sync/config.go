package sync

import (
	"time"
)

type Cursor struct {
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
}

type cursor struct {
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
}

// A cursor is valid for a week.
func (c cursor) Valid() bool {
	expiration := time.Now().Add(cursorValidDuration)
	return c.CreatedAt.Before(expiration)
}

type AutoAssigner struct {
	Filter      string `json:"filter"`
	ProfileUUID string `json:"profile_uuid"`
}

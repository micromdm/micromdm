package id

import (
	"crypto/rand"

	"github.com/oklog/ulid"
)

func New() string {
	return ulid.MustNew(ulid.Now(), rand.Reader).String()
}

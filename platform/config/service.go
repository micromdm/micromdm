package config

import (
	"context"

	"github.com/pkg/errors"
)

type Service interface {
	SavePushCertificate(ctx context.Context, cert, key []byte) error
}

type ConfigService struct {
	store *DB
}

func NewService(db *DB) *ConfigService {
	return &ConfigService{store: db}
}

func (svc *ConfigService) SavePushCertificate(ctx context.Context, cert, key []byte) error {
	err := svc.store.SavePushCertificate(cert, key)
	return errors.Wrap(err, "save push certificate")
}

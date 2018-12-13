package config

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
)

type Service interface {
	SavePushCertificate(ctx context.Context, cert []byte, key []byte) error
	GetPushCertificate(ctx context.Context) ([]byte, error)
	ApplyDEPToken(ctx context.Context, P7MContent []byte) error
	GetDEPTokens(ctx context.Context) ([]DEPToken, []byte, error)
}

type Store interface {
	SavePushCertificate(ctx context.Context, cert []byte, key []byte) error
	GetPushCertificate(ctx context.Context) ([]byte, error)
	PushCertificate(ctx context.Context) (*tls.Certificate, error)
	PushTopic(ctx context.Context) (string, error)
	DEPKeypair(ctx context.Context) (key *rsa.PrivateKey, cert *x509.Certificate, err error)
	AddToken(ctx context.Context, consumerKey string, json []byte) error
	DEPTokens(ctx context.Context) ([]DEPToken, error)
}

type ConfigService struct {
	store Store
}

func New(store Store) *ConfigService {
	return &ConfigService{store: store}
}

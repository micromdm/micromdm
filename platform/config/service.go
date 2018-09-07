package config

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
)

type Service interface {
	SavePushCertificate(ctx context.Context, cert, key []byte) error
	GetPushCertificate(ctx context.Context) ([]byte, error)
	ApplyDEPToken(ctx context.Context, P7MContent []byte) error
	GetDEPTokens(ctx context.Context) ([]DEPToken, []byte, error)
}

type Store interface {
	SavePushCertificate(cert, key []byte) error
	GetPushCertificate() ([]byte, error)
	PushCertificate() (*tls.Certificate, error)
	PushTopic() (string, error)
	DEPKeypair() (key *rsa.PrivateKey, cert *x509.Certificate, err error)
	AddToken(consumerKey string, json []byte) error
	DEPTokens() ([]DEPToken, error)
}

type ConfigService struct {
	store Store
}

func New(store Store) *ConfigService {
	return &ConfigService{store: store}
}

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
	ApplyVPPToken(ctx context.Context, Content []byte) error
	GetVPPTokens(ctx context.Context) ([]VPPToken, error)
	RemoveVPPToken(ctx context.Context, UDID []byte) error
}

type Store interface {
	SavePushCertificate(cert, key []byte) error
	GetPushCertificate() ([]byte, error)
	PushCertificate() (*tls.Certificate, error)
	PushTopic() (string, error)
	DEPKeypair() (key *rsa.PrivateKey, cert *x509.Certificate, err error)
	AddDEPToken(consumerKey string, json []byte) error
	AddVPPToken(sToken string, json []byte) error
	DeleteVPPToken(sToken string) error
	DEPTokens() ([]DEPToken, error)
	VPPTokens() ([]VPPToken, error)
}

type ConfigService struct {
	store Store
}

func New(store Store) *ConfigService {
	return &ConfigService{store: store}
}

package scep

import (
	"crypto/rsa"
	"crypto/x509"
	"math/big"
)

/*
type Service interface {
	ApplyProfile(ctx context.Context, p *Profile) error
	GetProfiles(ctx context.Context, opt GetProfilesOption) ([]Profile, error)
	RemoveProfiles(ctx context.Context, ids []string) error
}

type GetProfilesOption struct {
	Identifier string `json:"id"`
}
*/

type Store interface {
	CA(pass []byte) ([]*x509.Certificate, *rsa.PrivateKey, error)
	Put(cn string, crt *x509.Certificate) error
	Serial() (*big.Int, error)
	HasCN(cn string, allowTime int, cert *x509.Certificate, revokeOldCertificate bool) (bool, error)
	CreateOrLoadKey(bits int) (*rsa.PrivateKey, error)
	CreateOrLoadCA(key *rsa.PrivateKey, years int, org, country string) (*x509.Certificate, error)
}

func New(store Store) *ScepService {
	return &ScepService{store: store}
}

type ScepService struct {
	store Store
}

func IsNotFound(err error) bool {
	type notFoundError interface {
		error
		NotFound() bool
	}

	_, ok := err.(notFoundError)
	return ok
}

package server

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"

	"github.com/micromdm/micromdm/mdm"
)

type ScepVerifyDepot interface {
	CA(pass []byte) ([]*x509.Certificate, *rsa.PrivateKey, error)
	HasCN(cn string, allowTime int, cert *x509.Certificate, revokeOldCertificate bool) (bool, error)
}

func VerifyCertificateMiddleware(validateSCEPIssuer bool, store ScepVerifyDepot, logger log.Logger) mdm.Middleware {
	return func(next mdm.Service) mdm.Service {
		return &verifyCertificateMiddleware{
			store:              store,
			next:               next,
			logger:             logger,
			validateSCEPIssuer: validateSCEPIssuer,
		}
	}
}

type verifyCertificateMiddleware struct {
	store              ScepVerifyDepot
	next               mdm.Service
	logger             log.Logger
	validateSCEPIssuer bool
}

func (mw *verifyCertificateMiddleware) verifyIssuer(devcert *x509.Certificate) error {
	expiration := devcert.NotAfter

	if time.Now().After(expiration) {
		return errors.New("device certificate is expired")
	}

	ca, _, err := mw.store.CA(nil)
	if err != nil {
		return errors.Wrap(err, "error retrieving CA")
	}

	roots := x509.NewCertPool()
	for _, cert := range ca {
		roots.AddCert(cert)
	}

	opts := x509.VerifyOptions{
		Roots: roots,
		KeyUsages: []x509.ExtKeyUsage{
			x509.ExtKeyUsageAny,
		},
	}

	if _, err := devcert.Verify(opts); err != nil {
		return errors.Wrap(err, "error verifying certificate")
	}

	return nil
}

func (mw *verifyCertificateMiddleware) Acknowledge(ctx context.Context, req mdm.AcknowledgeEvent) ([]byte, error) {
	devcert, err := mdm.DeviceCertificateFromContext(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error retrieving device certificate")
	}
	hasCN, err := mw.store.HasCN(devcert.Subject.CommonName, 0, devcert, false)
	if err != nil {
		return nil, errors.Wrap(err, "error checking device certificate")
	}

	unauth_err := errors.New("unauthorized client")
	if !hasCN && !mw.validateSCEPIssuer {
		_ = level.Info(mw.logger).Log("err", unauth_err, "issuer", devcert.Issuer.String(), "expiration", devcert.NotAfter)
		return nil, unauth_err
	}
	if !hasCN && mw.validateSCEPIssuer {
		err := mw.verifyIssuer(devcert)
		if err != nil {
			_ = level.Info(mw.logger).Log("err", err, "issuer", devcert.Issuer.String(), "expiration", devcert.NotAfter)
			return nil, unauth_err
		}
	}

	return mw.next.Acknowledge(ctx, req)
}

func (mw *verifyCertificateMiddleware) Checkin(ctx context.Context, req mdm.CheckinEvent) error {
	devcert, err := mdm.DeviceCertificateFromContext(ctx)
	if err != nil {
		return errors.Wrap(err, "error retrieving device certificate")
	}
	hasCN, err := mw.store.HasCN(devcert.Subject.CommonName, 0, devcert, false)
	if err != nil {
		return errors.Wrap(err, "error checking device certificate")
	}
	unauth_err := errors.New("unauthorized client")
	if !hasCN && !mw.validateSCEPIssuer {
		_ = level.Info(mw.logger).Log("err", unauth_err, "issuer", devcert.Issuer.String(), "expiration", devcert.NotAfter)
		return unauth_err
	}
	if !hasCN && mw.validateSCEPIssuer {
		err := mw.verifyIssuer(devcert)
		if err != nil {
			_ = level.Info(mw.logger).Log("err", err, "issuer", devcert.Issuer.String(), "expiration", devcert.NotAfter)
			return unauth_err
		}
	}
	return mw.next.Checkin(ctx, req)
}

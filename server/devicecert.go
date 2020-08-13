package server

import (
	"context"
	"crypto/x509"
	"fmt"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"

	"github.com/micromdm/micromdm/mdm"
	boltdepot "github.com/micromdm/scep/depot/bolt"
)

type ScepVerifyDepot interface {
	HasCN(cn string, allowTime int, cert *x509.Certificate, revokeOldCertificate bool) (bool, error)
}

func VerifyCertificateMiddleware(validateSCEPIssuer bool, scepIssuer string, scepDepot *boltdepot.Depot, store ScepVerifyDepot, logger log.Logger) mdm.Middleware {
	return func(next mdm.Service) mdm.Service {
		return &verifyCertificateMiddleware{
			store:              store,
			next:               next,
			logger:             logger,
			validateSCEPIssuer: validateSCEPIssuer,
			scepIssuer:         scepIssuer,
			scepDepot:          scepDepot,
		}
	}
}

type verifyCertificateMiddleware struct {
	store              ScepVerifyDepot
	next               mdm.Service
	logger             log.Logger
	validateSCEPIssuer bool
	scepIssuer         string
	scepDepot          *boltdepot.Depot
}

func verifyIssuer(devcert *x509.Certificate, scepIssuer string, scepDepot *boltdepot.Depot) (bool, error) {
	issuer := devcert.Issuer.String()
	expiration := devcert.NotAfter

	if time.Now().After(expiration) {
		err := errors.New("device certificate is expired")
		return false, err
	}

	ca, _, err := scepDepot.CA(nil)
	if err != nil {
		return false, errors.Wrap(err, "error retrieving CA")
	}

	roots := x509.NewCertPool()
	for _, cert := range ca {
		roots.AddCert(cert)
	}

	opts := x509.VerifyOptions{
		Roots: roots,
	}

	if _, err := devcert.Verify(opts); err != nil {
		return false, errors.Wrap(err, "error verifying certificate")
	}

	if issuer != scepIssuer {
		err := fmt.Errorf("device certificate not issued by %v", scepIssuer)
		return false, err
	}

	return true, nil
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
	if !hasCN {
		verifiedIssuer := false
		if mw.validateSCEPIssuer {
			verifiedIssuer, err = verifyIssuer(devcert, mw.scepIssuer, mw.scepDepot)
			if err != nil {
				_ = level.Info(mw.logger).Log("err", err, "issuer", devcert.Issuer.String(), "expiration", devcert.NotAfter)
				return nil, errors.Wrap(err, "error verifying CN")
			}
		}

		if !verifiedIssuer {
			err := errors.New("unauthorized client")
			_ = level.Info(mw.logger).Log("err", err, "issuer", devcert.Issuer.String(), "expiration", devcert.NotAfter)
			return nil, err
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
	if !hasCN {
		verifiedIssuer := false
		if mw.validateSCEPIssuer {
			verifiedIssuer, err = verifyIssuer(devcert, mw.scepIssuer, mw.scepDepot)
			if err != nil {
				_ = level.Info(mw.logger).Log("err", err, "issuer", devcert.Issuer.String(), "expiration", devcert.NotAfter)
				return errors.Wrap(err, "error verifying CN")
			}
		}

		if !verifiedIssuer {
			err := errors.New("unauthorized client")
			_ = level.Info(mw.logger).Log("err", err, "issuer", devcert.Issuer.String(), "expiration", devcert.NotAfter)
			return err
		}
	}
	return mw.next.Checkin(ctx, req)
}

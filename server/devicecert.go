package server

import (
	"context"
	"crypto/x509"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"

	"github.com/micromdm/micromdm/mdm"
)

type ScepVerifyDepot interface {
	HasCN(cn string, allowTime int, cert *x509.Certificate, revokeOldCertificate bool) (bool, error)
}

func VerifyCertificateMiddleware(store ScepVerifyDepot, logger log.Logger) mdm.Middleware {
	return func(next mdm.Service) mdm.Service {
		return &verifyCertificateMiddleware{
			store:  store,
			next:   next,
			logger: logger,
		}
	}
}

type verifyCertificateMiddleware struct {
	store  ScepVerifyDepot
	next   mdm.Service
	logger log.Logger
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
		err := errors.New("unauthorized client")
		level.Info(mw.logger).Log("err", err)
		return nil, err
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
		err := errors.New("unauthorized client")
		level.Info(mw.logger).Log("err", err)
		return err
	}
	return mw.next.Checkin(ctx, req)
}

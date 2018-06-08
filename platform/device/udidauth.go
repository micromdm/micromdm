package device

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"

	"github.com/micromdm/micromdm/mdm"
)

type UDIDCertAuthStore interface {
	SaveUDIDCertHash(udid, certHash []byte) error
	GetUDIDCertHash(udid []byte) ([]byte, error)
}

func UDIDCertAuthMiddleware(store UDIDCertAuthStore, logger log.Logger) mdm.Middleware {
	return func(next mdm.Service) mdm.Service {
		return &udidCertAuthMiddleware{
			store:  store,
			next:   next,
			logger: log.With(logger, "component", "udidcertauth"),
		}
	}
}

type udidCertAuthMiddleware struct {
	store  UDIDCertAuthStore
	next   mdm.Service
	logger log.Logger
}

func hashCertRaw(c []byte) []byte {
	ret_bytes := make([]byte, 32)
	sum := sha256.Sum256(c)
	copy(ret_bytes, sum[:])
	return ret_bytes
}

func (mw *udidCertAuthMiddleware) validateUDIDCertAuth(udid, certHash []byte) (bool, error) {
	dbCertHash, err := mw.store.GetUDIDCertHash(udid)
	if err != nil && !isNotFound(err) {
		return false, err
	} else if err != nil && isNotFound(err) {
		// TODO: we did not find any UDID at all. assume (but log) that
		// this device already existed/was enrolled and we need to store
		// its UDID-cert association. at some later late, when most/all
		// micromdm instances have stored udid-cert associations
		// this can be an outright failure.
		level.Info(mw.logger).Log("msg", "device cert hash not found, saving anyway", "udid", string(udid))
		if err := mw.store.SaveUDIDCertHash(udid, certHash); err != nil {
			return false, err
		}
		return true, nil
	}
	if 1 == subtle.ConstantTimeCompare(certHash, dbCertHash) {
		return true, nil
	}
	level.Info(mw.logger).Log("msg", "device cert hash mismatch", "udid", string(udid))
	return false, nil
}

func (mw *udidCertAuthMiddleware) Acknowledge(ctx context.Context, req mdm.AcknowledgeEvent) ([]byte, error) {
	devcert, err := mdm.DeviceCertificateFromContext(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error retrieving device certificate")
	}
	matched, err := mw.validateUDIDCertAuth([]byte(req.Response.UDID), hashCertRaw(devcert.Raw))
	if err != nil {
		return nil, err
	}
	if !matched {
		return nil, errors.New("device certifcate UDID mismatch")
	}
	return mw.next.Acknowledge(ctx, req)
}

func (mw *udidCertAuthMiddleware) Checkin(ctx context.Context, req mdm.CheckinEvent) error {
	devcert, err := mdm.DeviceCertificateFromContext(ctx)
	if err != nil {
		return errors.Wrap(err, "error retrieving device certificate")
	}
	switch req.Command.MessageType {
	case "Authenticate":
		// unconditionally save the certificate on the Authenticate step
		if err := mw.store.SaveUDIDCertHash([]byte(req.Command.UDID), hashCertRaw(devcert.Raw)); err != nil {
			return err
		}
		return mw.next.Checkin(ctx, req)
	case "TokenUpdate", "CheckOut":
		matched, err := mw.validateUDIDCertAuth([]byte(req.Command.UDID), hashCertRaw(devcert.Raw))
		if err != nil {
			return err
		}
		if !matched {
			return errors.New("device certifcate UDID mismatch")
		}
		return mw.next.Checkin(ctx, req)
	default:
		return errors.Errorf("unknown checkin message type %s", req.Command.MessageType)
	}
	return mw.next.Checkin(ctx, req)
}

package mw

import (
	"bytes"
	"context"
	"crypto/x509"
	"fmt"

	"github.com/pkg/errors"

	"github.com/micromdm/micromdm/mdm"
	devicebuiltin "github.com/micromdm/micromdm/platform/device/builtin"
)

func UDIDAuthMiddleware(store *devicebuiltin.DB) mdm.Middleware {
	return func(next mdm.Service) mdm.Service {
		return &udidAuthMiddleware{
			store: store,
			next:  next,
		}
	}
}

type udidAuthMiddleware struct {
	store *devicebuiltin.DB
	next  mdm.Service
}

func (mw *udidAuthMiddleware) checkCertUDIDMatch(cert *x509.Certificate, udid string) (bool, error) {
	devByUUID, err := mw.store.DeviceByUDID(udid)
	if err != nil {
		if devicebuiltin.IsNotFound(err) {
			return false, errors.New("device not found; no certificate UDID match")
		}
		return false, err
	}

	if devByUUID.DeviceCert == nil {
		// nil certificate. this could be a potentially dangerous thing
		// to allow but if, truely, there is no certificate then the
		// most likely thing is that this is a "legacy" device that
		// needs its certificate updated in the database.
		//
		// but for testing simply fail out.

		// upgrade any devices that don't have a certificate to have one
		devByUUID.DeviceCert = cert
		err = mw.store.Save(devByUUID)
		if err != nil {
			return false, err
		}
		fmt.Println("SAVED NEW!")
		return true, nil
	}
	if bytes.Compare(cert.Raw, devByUUID.DeviceCert.Raw) == 0 {
		fmt.Println("MATCH!")
		return true, nil
	}
	return false, nil
}

func (mw *udidAuthMiddleware) Acknowledge(ctx context.Context, req mdm.AcknowledgeEvent) ([]byte, error) {
	// fmt.Println("CERT<->UUID CHECKING HERE", req.Command.UDID, req.DeviceCert)
	return mw.next.Acknowledge(ctx, req)
}

func (mw *udidAuthMiddleware) Checkin(ctx context.Context, req mdm.CheckinEvent) error {
	switch req.Command.MessageType {
	case "Authenticate":
		// we'll pass this through as this will have been the first time
		// we've seen the device certificate and we need to pass the
		// pubsub message to save the checkin
		return mw.next.Checkin(ctx, req)
	case "TokenUpdate", "CheckOut":
		udidMatch, err := mw.checkCertUDIDMatch(req.DeviceCert, req.Command.UDID)
		if err != nil {
			return errors.Wrap(err, "checking device certificate UDID match")
		}
		if !udidMatch {
			return errors.New("device certifcate UDID mismatch")
		}
		return mw.next.Checkin(ctx, req)
	default:
		return errors.Errorf("unknown checkin message type %s", req.Command.MessageType)
	}
}

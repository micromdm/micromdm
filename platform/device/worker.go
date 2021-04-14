package device

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/micromdm/micromdm/mdm"
	"github.com/micromdm/micromdm/platform/dep/sync"
	"github.com/micromdm/micromdm/platform/pubsub"
)

type DeviceWorkerStore interface {
	Save(ctx context.Context, d *Device) error
	DeviceByUDID(ctx context.Context, udid string) (*Device, error)
	DeviceBySerial(ctx context.Context, serial string) (*Device, error)
}

type Worker struct {
	db     DeviceWorkerStore
	ps     pubsub.PublishSubscriber
	logger log.Logger
}

func NewWorker(db DeviceWorkerStore, ps pubsub.PublishSubscriber, logger log.Logger) *Worker {
	return &Worker{
		db:     db,
		ps:     ps,
		logger: logger,
	}
}

func (w *Worker) Run(ctx context.Context) error {
	const subscription = "devices_worker"
	authenticateEvents, err := w.ps.Subscribe(ctx, subscription, mdm.AuthenticateTopic)
	if err != nil {
		return errors.Wrapf(err, "subscribing %s to %s", subscription, mdm.AuthenticateTopic)
	}
	tokenUpdateEvents, err := w.ps.Subscribe(ctx, subscription, mdm.TokenUpdateTopic)
	if err != nil {
		return errors.Wrapf(err, "subscribing %s to %s", subscription, mdm.TokenUpdateTopic)
	}
	getBootstrapTokenEvents, err := w.ps.Subscribe(ctx, subscription, mdm.GetBootstrapTokenTopic)
	if err != nil {
		return errors.Wrapf(err, "subscribing %s to %s", subscription, mdm.GetBootstrapTokenTopic)
	}
	setBootstrapTokenEvents, err := w.ps.Subscribe(ctx, subscription, mdm.SetBootstrapTokenTopic)
	if err != nil {
		return errors.Wrapf(err, "subscribing %s to %s", subscription, mdm.SetBootstrapTokenTopic)
	}
	checkoutEvents, err := w.ps.Subscribe(ctx, subscription, mdm.CheckoutTopic)
	if err != nil {
		return errors.Wrapf(err, "subscribing %s to %s", subscription, mdm.CheckoutTopic)
	}
	depSyncEvents, err := w.ps.Subscribe(ctx, subscription, sync.SyncTopic)
	if err != nil {
		return errors.Wrapf(err, "subscribing %s to %s", subscription, sync.SyncTopic)
	}
	connectEvents, err := w.ps.Subscribe(ctx, subscription, mdm.ConnectTopic)
	if err != nil {
		return errors.Wrapf(err, "subscribing %s to %s", subscription, mdm.ConnectTopic)
	}

	for {
		var err error
		select {
		case <-ctx.Done():
			return ctx.Err()
		case ev := <-authenticateEvents:
			err = w.updateFromAuthenticate(ctx, ev.Message)
		case ev := <-tokenUpdateEvents:
			err = w.updateFromTokenUpdate(ctx, ev.Message)
		case ev := <-getBootstrapTokenEvents:
			err = w.updateFromGetBootstrapToken(ctx, ev.Message)
		case ev := <-setBootstrapTokenEvents:
			err = w.updateFromSetBootstrapToken(ctx, ev.Message)
		case ev := <-checkoutEvents:
			err = w.updateFromCheckout(ctx, ev.Message)
		case ev := <-depSyncEvents:
			err = w.updateFromDEPSync(ctx, ev.Message)
		case ev := <-connectEvents:
			err = w.updateFromAcknowledge(ctx, ev.Message)
		}
		if err != nil {
			level.Info(w.logger).Log(
				"msg", "update device from event",
				"err", err,
			)
			continue
		}
	}
}

func (w *Worker) updateFromDEPSync(ctx context.Context, message []byte) error {
	var ev sync.Event
	if err := sync.UnmarshalEvent(message, &ev); err != nil {
		return errors.Wrap(err, "unmarshal depsync event")
	}
	level.Debug(w.logger).Log(
		"msg", "updating devices from DEP",
		"device_count", len(ev.Devices),
	)

	for _, dd := range ev.Devices {
		dev, err := getOrCreateDeviceBySerial(ctx, w.db, dd.SerialNumber)
		if err != nil {
			return errors.Wrap(err, "get device by serial")
		}

		notSeenBefore := dev.UUID == ""
		logEvent := dd.OpType == "deleted" ||
			(notSeenBefore && dd.OpType == "added") ||
			(!notSeenBefore && dd.OpType == "modified")
		if logEvent {
			level.Debug(w.logger).Log(
				"msg", "updating devices from dep sync",
				"op_type", dd.OpType,
				"serial", dd.SerialNumber,
				"previously_known", !notSeenBefore,
			)
		}

		if dev.UUID == "" {
			dev.UUID = uuid.New().String()
		}

		dev.SerialNumber = dd.SerialNumber
		dev.Model = dd.Model
		dev.Description = dd.Description
		dev.Color = dd.Color
		dev.AssetTag = dd.AssetTag
		dev.DEPProfileStatus = DEPProfileStatus(dd.ProfileStatus)
		dev.DEPProfileUUID = dd.ProfileUUID
		dev.DEPProfileAssignTime = dd.ProfileAssignTime
		dev.DEPProfileAssignedDate = dd.DeviceAssignedDate
		dev.DEPProfileAssignedBy = dd.DeviceAssignedBy

		if err := w.db.Save(ctx, dev); err != nil {
			return errors.Wrap(err, "save device %s from DEP sync")
		}
	}

	return nil
}

func (w *Worker) updateFromAcknowledge(ctx context.Context, message []byte) error {
	var ev mdm.AcknowledgeEvent
	if err := mdm.UnmarshalAcknowledgeEvent(message, &ev); err != nil {
		return errors.Wrap(err, "unmarshal acknowledge event")
	}

	if ev.Response.EnrollmentID != nil {
		return nil
	}

	dev, err := w.db.DeviceByUDID(ctx, ev.Response.UDID)
	if err != nil {
		return errors.Wrapf(err, "retrieve device with udid %s", ev.Response.UDID)
	}
	dev.LastSeen = time.Now()

	err = w.db.Save(ctx, dev)
	return errors.Wrapf(err, "saving updated device for acknowledge event")

}

func (w *Worker) updateFromCheckout(ctx context.Context, message []byte) error {
	var ev mdm.CheckinEvent
	if err := mdm.UnmarshalCheckinEvent(message, &ev); err != nil {
		return errors.Wrap(err, "unmarshal checkin event")
	}

	if ev.Command.EnrollmentID != "" {
		return nil
	}

	dev, err := w.db.DeviceByUDID(ctx, ev.Command.UDID)
	if err != nil {
		return errors.Wrapf(err, "retrieve device with udid %s", ev.Command.UDID)
	}

	dev.Enrolled = false
	dev.LastSeen = time.Now()

	err = w.db.Save(ctx, dev)
	return errors.Wrapf(err, "saving updated device for checkout event")

}

func (w *Worker) updateFromGetBootstrapToken(ctx context.Context, message []byte) error {
	var ev mdm.CheckinEvent
	if err := mdm.UnmarshalCheckinEvent(message, &ev); err != nil {
		return errors.Wrap(err, "unmarshal checkin event")
	}

	dev, err := w.db.DeviceByUDID(ctx, ev.Command.UDID)
	if err != nil {
		return errors.Wrapf(err, "retrieve device with udid %s", ev.Command.UDID)
	}

	dev.AwaitingConfiguration = ev.Command.AwaitingConfiguration
	dev.LastSeen = time.Now()

	err = w.db.Save(ctx, dev)
	return errors.Wrapf(err, "saving updated device for GetBootstrapToken event")

}

func (w *Worker) updateFromSetBootstrapToken(ctx context.Context, message []byte) error {
	var ev mdm.CheckinEvent
	if err := mdm.UnmarshalCheckinEvent(message, &ev); err != nil {
		return errors.Wrap(err, "unmarshal checkin event")
	}

	dev, err := w.db.DeviceByUDID(ctx, ev.Command.UDID)
	if err != nil {
		return errors.Wrapf(err, "retrieve device with udid %s", ev.Command.UDID)
	}

	dev.BootstrapToken = ev.Command.BootstrapToken.String()
	dev.AwaitingConfiguration = ev.Command.AwaitingConfiguration
	dev.LastSeen = time.Now()

	err = w.db.Save(ctx, dev)
	return errors.Wrapf(err, "saving updated device for SetBootstrapToken event")

}

func (w *Worker) updateFromTokenUpdate(ctx context.Context, message []byte) error {
	var ev mdm.CheckinEvent
	if err := mdm.UnmarshalCheckinEvent(message, &ev); err != nil {
		return errors.Wrap(err, "unmarshal checkin event")
	}

	// do not process managed user, or user enrollment checkin events while updating device records.
	if ev.Command.UserID != "" || ev.Command.EnrollmentID != "" {
		return nil
	}

	dev, err := w.db.DeviceByUDID(ctx, ev.Command.UDID)
	if err != nil {
		return errors.Wrapf(err, "retrieve device with udid %s", ev.Command.UDID)
	}
	dev.Token = ev.Command.Token.String()
	dev.PushMagic = ev.Command.PushMagic
	dev.UnlockToken = ev.Command.UnlockToken.String()
	dev.AwaitingConfiguration = ev.Command.AwaitingConfiguration
	dev.LastSeen = time.Now()
	// first TokenUpdate event will have the enrollment status set to false.
	newlyEnrolled := !dev.Enrolled
	dev.Enrolled = true
	if err := w.db.Save(ctx, dev); err != nil {
		return errors.Wrapf(err, "saving updated device for Token event udid=%s", ev.Command.UDID)
	}

	if newlyEnrolled {
		// notify subscribers of a successful enrollment
		// TODO: The enrollment topic needs a custom event.
		err = w.ps.Publish(ctx, DeviceEnrolledTopic, message)
		return errors.Wrap(err, "publishing new enrollment message")
	}
	return nil
}

func (w *Worker) updateFromAuthenticate(ctx context.Context, message []byte) error {
	var ev mdm.CheckinEvent
	if err := mdm.UnmarshalCheckinEvent(message, &ev); err != nil {
		return errors.Wrap(err, "unmarshal checkin event")
	}

	if ev.Command.EnrollmentID != "" {
		return nil
	}

	device, reenrolling, err := getOrCreateDevice(ctx, w.db, ev.Command.SerialNumber, ev.Command.UDID)
	if err != nil {
		return errors.Wrap(err, "get device for authenticate event")
	}

	if reenrolling {
		level.Debug(w.logger).Log(
			"msg", "re-enrolling device",
			"serial", ev.Command.SerialNumber,
		)
	} else {
		level.Debug(w.logger).Log(
			"msg", "enrolling new device",
			"serial", ev.Command.SerialNumber,
		)
	}

	if device.UUID == "" {
		device.UUID = uuid.New().String()
	}
	device.UDID = ev.Command.UDID
	device.OSVersion = ev.Command.OSVersion
	device.BuildVersion = ev.Command.BuildVersion
	device.ProductName = ev.Command.ProductName
	device.SerialNumber = ev.Command.SerialNumber
	device.IMEI = ev.Command.IMEI
	device.MEID = ev.Command.MEID
	device.DeviceName = ev.Command.DeviceName
	device.Model = ev.Command.Model
	device.ModelName = ev.Command.ModelName
	device.LastSeen = time.Now()
	err = w.db.Save(ctx, device)
	return errors.Wrapf(err, "saving updated device for authenticate event")
}

func getOrCreateDevice(ctx context.Context, db DeviceWorkerStore, serial, udid string) (dev *Device, reenrolling bool, err error) {
	if udid != "" {
		// first try to fetch a device by UDID.
		// If the device was previously enrolled it will exist.
		// In case the device is known, set the enrolled status to false before returning.
		byUDID, err := db.DeviceByUDID(ctx, udid)
		if err == nil {
			byUDID.Enrolled = false
			return byUDID, true, nil
		}
		if err != nil && !isNotFound(err) {
			return nil, false, errors.Wrapf(err, "retrieve device with udid %s and serial %s", udid, serial)
		}
	}

	// next try to find the device by serial. If found, it's a DEP device, which contains only the
	// serials but not a udid.
	dev, err = getOrCreateDeviceBySerial(ctx, db, serial)
	return dev, false, err
}

func getOrCreateDeviceBySerial(ctx context.Context, db DeviceWorkerStore, serial string) (*Device, error) {
	bySerial, err := db.DeviceBySerial(ctx, serial)
	if err == nil && bySerial != nil {
		return bySerial, nil
	}
	if err != nil && !isNotFound(err) {
		return nil, errors.Wrapf(err, "retrieve device with serial number %s", serial)
	}

	dev := new(Device)
	return dev, nil
}

func isNotFound(err error) bool {
	err = errors.Cause(err)
	type notFoundErr interface {
		error
		NotFound() bool
	}

	e, ok := err.(notFoundErr)
	return ok && e.NotFound()
}

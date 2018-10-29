package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"

	"github.com/micromdm/micromdm/workflow/profile/device"
)

type ListWorkerDeviceStore interface {
	List(ctx context.Context) ([]device.Device, error)
}

type ListWorker struct {
	deviceDB ListWorkerDeviceStore
	logger   log.Logger
}

func NewWorker(logger log.Logger, db ListWorkerDeviceStore) *ListWorker {
	return &ListWorker{
		deviceDB: db,
		logger:   logger,
	}
}

func (w *ListWorker) Run(ctx context.Context) error {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	fn := func() {
		logger := level.Debug(w.logger)
		err := w.queueProfileListForAll(ctx)
		if err != nil {
			logger = level.Info(w.logger)
		}
		logger.Log("msg", "queue profile list", "err", err)
	}

	fn()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			fn()
		}
	}
}

func (w *ListWorker) queueProfileListForAll(ctx context.Context) error {
	cctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	list, err := w.deviceDB.List(cctx)
	if err != nil {
		return err
	}

	for _, dev := range list {
		if err := w.sendProfileListCommand(ctx, dev.UDID); err != nil {
			return errors.Wrap(err, "sending profile list command")
		}
	}

	return nil
}

func (w *ListWorker) sendProfileListCommand(ctx context.Context, udid string) error {
	level.Debug(w.logger).Log("msg", "sending list profiles command", "udid", udid)

	cctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var postBody bytes.Buffer
	var request = struct {
		UDID        string `json:"udid"`
		RequestType string `json:"request_type"`
	}{
		UDID:        udid,
		RequestType: "ProfileList",
	}
	if err := json.NewEncoder(&postBody).Encode(&request); err != nil {
		return err
	}
	req, err := http.NewRequest("POST", "https://3470b1a1.ngrok.io/v1/commands", &postBody)
	if err != nil {
		return err
	}
	req.SetBasicAuth("micromdm", "supersecret")
	req = req.WithContext(cctx)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

package apns

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/RobotsAndPencils/buford/payload"
	"github.com/RobotsAndPencils/buford/push"
	"github.com/go-kit/kit/endpoint"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/micromdm/micromdm/pkg/httputil"
)

func (svc *PushService) Push(ctx context.Context, deviceUDID string) (string, error) {
	info, err := svc.store.PushInfo(ctx, deviceUDID)
	if err != nil {
		return "", errors.Wrap(err, "retrieving PushInfo by UDID")
	}

	p := payload.MDM{Token: info.PushMagic}
	valid := push.IsDeviceTokenValid(info.Token)
	if !valid {
		return "", errors.New("invalid push token")
	}
	jsonPayload, err := json.Marshal(p)
	if err != nil {
		return "", errors.Wrap(err, "marshalling push notification payload")
	}
	result, err := svc.pushsvc.Push(info.Token, nil, jsonPayload)
	if err != nil && strings.HasSuffix(err.Error(), "remote error: tls: internal error") {
		// TODO: yuck, error substring searching. see:
		// https://github.com/micromdm/micromdm/issues/150
		return result, errors.Wrap(err, "push error: possibly expired or invalid APNs certificate")
	}
	return result, err
}

type pushRequest struct {
	UDID string
}

type pushResponse struct {
	Status string `json:"status,omitempty"`
	ID     string `json:"push_notification_id,omitempty"`
	Err    error  `json:"error,omitempty"`
}

func (r pushResponse) Failed() error { return r.Err }

func decodePushRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	udid, ok := vars["udid"]
	if !ok {
		return 0, errors.New("apns: bad route")
	}
	return pushRequest{
		UDID: udid,
	}, nil
}

func decodePushResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var resp pushResponse
	err := httputil.DecodeJSONResponse(r, &resp)
	return resp, err
}

func MakePushEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(pushRequest)
		id, err := svc.Push(ctx, req.UDID)
		if err != nil {
			return pushResponse{Err: err, Status: "failure"}, nil
		}
		return pushResponse{Status: "success", ID: id}, nil
	}
}

func (mw loggingMiddleware) Push(ctx context.Context, udid string) (id string, err error) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "Push",
			"udid", udid,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	id, err = mw.next.Push(ctx, udid)
	return
}

package mdm

import (
	"context"
	"net/http"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/gorilla/mux"
	"github.com/groob/plist"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (svc *MDMService) Checkin(ctx context.Context, event CheckinEvent) error {
	if event.Command.MessageType == "UserAuthenticate" {
		return &rejectUserAuth{}
	}

	msg, err := MarshalCheckinEvent(&event)
	if err != nil {
		return errors.Wrap(err, "marshal checkin event")
	}

	topic, err := topicFromMessage(event.Command.MessageType)
	if err != nil {
		return errors.Wrap(err, "get checkin topic from message")
	}

	err = svc.pub.Publish(ctx, topic, msg)
	return errors.Wrapf(err, "publish checkin on topic: %s", topic)
}

func topicFromMessage(messageType string) (string, error) {
	switch messageType {
	case "Authenticate":
		return AuthenticateTopic, nil
	case "TokenUpdate":
		return TokenUpdateTopic, nil
	case "CheckOut":
		return CheckoutTopic, nil
	default:
		return "", errors.Errorf("unknown checkin message type %s", messageType)
	}
}

type rejectUserAuth struct{}

func (e *rejectUserAuth) Error() string {
	return "reject user auth"
}
func (e *rejectUserAuth) UserAuthReject() bool {
	return true
}

func isRejectedUserAuth(err error) bool {
	type rejectUserAuthError interface {
		error
		UserAuthReject() bool
	}

	_, ok := errors.Cause(err).(rejectUserAuthError)
	return ok
}

type checkinRequest struct {
	Event CheckinEvent
}

type checkinResponse struct {
	Err error `plist:"error,omitempty"`
}

func (r checkinResponse) Failed() error { return r.Err }

func (d *requestDecoder) decodeCheckinRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	body, err := d.readBody(r)
	if err != nil {
		return nil, errors.Wrap(err, "read checkin request body")
	}

	var cmd CheckinCommand
	if err := plist.Unmarshal(body, &cmd); err != nil {
		return nil, errors.Wrap(err, "unmarshal MDM Checkin Request plist")
	}

	params := mux.Vars(r)
	event := CheckinEvent{
		ID:      uuid.NewV4().String(),
		Time:    time.Now().UTC(),
		Command: cmd,
		Params:  params,
		Raw:     body,
	}
	req := checkinRequest{Event: event}
	return req, nil
}

func MakeCheckinEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(checkinRequest)
		err := svc.Checkin(ctx, req.Event)
		return checkinResponse{Err: err}, nil
	}
}

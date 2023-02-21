package command

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/gorilla/mux"
	"github.com/groob/plist"
	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/micromdm/micromdm/mdm/mdm"
	"github.com/micromdm/micromdm/pkg/httputil"
)

const (
	// CommandTopic is a PubSub topic that events are published to.
	CommandTopic = "mdm.Command"

	// RawCommandTopic is a PubSub topic that events are published to.
	RawCommandTopic = "mdm.RawCommand"
)

func (svc *CommandService) NewCommand(ctx context.Context, request *mdm.CommandRequest) (*mdm.CommandPayload, error) {
	if request == nil {
		return nil, errors.New("empty CommandRequest")
	}
	payload, err := mdm.NewCommandPayload(request)
	if err != nil {
		return nil, errors.Wrap(err, "creating mdm payload")
	}
	event := NewEvent(payload, request.UDID)
	msg, err := MarshalEvent(event)
	if err != nil {
		return nil, errors.Wrap(err, "marshalling mdm command event")
	}
	if err := svc.publisher.Publish(context.TODO(), CommandTopic, msg); err != nil {
		return nil, errors.Wrapf(err, "publish mdm command on topic: %s", CommandTopic)
	}
	return payload, nil
}

func (svc *CommandService) NewRawCommand(ctx context.Context, cmd *RawCommand) error {
	if cmd == nil {
		return errors.New("empty RawCommand")
	}
	event := NewRawEvent(cmd)
	msg, err := MarshalRawEvent(event)
	if err != nil {
		return errors.Wrap(err, "marshalling raw mdm command event")
	}
	if err := svc.publisher.Publish(context.TODO(), RawCommandTopic, msg); err != nil {
		return errors.Wrapf(err, "publish raw mdm command on topic: %s", RawCommandTopic)
	}
	return nil
}

type newCommandRequest struct {
	mdm.CommandRequest
}

type newCommandResponse struct {
	Payload *mdm.CommandPayload `json:"payload,omitempty"`
	Err     error               `json:"error,omitempty"`
}

func (r newCommandResponse) Failed() error   { return r.Err }
func (r newCommandResponse) StatusCode() int { return http.StatusCreated }

func decodeNewCommandRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req newCommandRequest
	err := httputil.DecodeJSONRequest(r, &req)
	return req, err
}

var errEmptyRequest = errors.New("request must contain UDID of the device")

// MakeNewCommandEndpoint creates an endpoint which creates new MDM Commands.
func MakeNewCommandEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(newCommandRequest)
		if req.UDID == "" || req.RequestType == "" {
			return newCommandResponse{Err: errEmptyRequest}, nil
		}
		payload, err := svc.NewCommand(ctx, &req.CommandRequest)
		if err != nil {
			return newCommandResponse{Err: err}, nil
		}
		return newCommandResponse{Payload: payload}, nil
	}
}

type RawCommand struct {
	UDID        string `json:"udid" plist:"-"`
	CommandUUID string `json:"command_uuid"`
	Command     struct {
		RequestType string `json:"request_type"`
	} `json:"command"`
	Raw []byte `plist:"-" json:"payload"`
}

type newRawCommandRequest struct {
	RawCommand
}

type newRawCommandResponse struct {
	Payload *RawCommand `json:"payload,omitempty"`
	Err     error       `json:"error,omitempty"`
}

func (r newRawCommandResponse) Failed() error   { return r.Err }
func (r newRawCommandResponse) StatusCode() int { return http.StatusCreated }

func decodeNewRawCommandRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	udid, ok := mux.Vars(r)["udid"]
	if !ok {
		return nil, errors.New("empty udid")
	}

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("read payload body: %w", err)
	}

	// verify body is valid plist and parse CommandUUID and RequestType
	var req newRawCommandRequest
	if err := plist.NewXMLDecoder(bytes.NewBuffer(payload)).Decode(&req); err != nil {
		return nil, fmt.Errorf("parse payload as plist: %w", err)
	}

	req.UDID = udid
	req.Raw = payload
	return req, nil
}

// MakeNewRawCommandEndpoint creates an endpoint which creates new raw MDM Commands.
func MakeNewRawCommandEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(newRawCommandRequest)
		if req.CommandUUID == "" || req.Command.RequestType == "" {
			return newRawCommandResponse{Err: errEmptyRequest}, nil
		}
		if err := svc.NewRawCommand(ctx, &req.RawCommand); err != nil {
			return newRawCommandResponse{Err: err}, nil
		}
		return newRawCommandResponse{Payload: &req.RawCommand}, nil
	}
}

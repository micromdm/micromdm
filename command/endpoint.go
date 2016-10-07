package command

import (
	"errors"

	"golang.org/x/net/context"

	"github.com/go-kit/kit/endpoint"
	"github.com/micromdm/mdm"
)

var (
	// ErrEmptyRequest is returned if the request body is empty
	ErrEmptyRequest = errors.New("request must contain UDID of the device")
	errBadRouting   = errors.New("inconsistent mapping between route and handler (programmer error)")
)

// newCommandRequest represents an HTTP Request for a new MDM Command
type newCommandRequest struct {
	*mdm.CommandRequest
}

// newCommandResponse is a command reponse
type newCommandResponse struct {
	*mdm.Payload
	Err error `json:"error,omitempty"`
}

func (r newCommandResponse) error() error { return r.Err }

func makeNewCommandEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(newCommandRequest)
		if req.UDID == "" || req.RequestType == "" {
			return newCommandResponse{Err: ErrEmptyRequest}, nil
		}
		payload, err := svc.NewCommand(req.CommandRequest)
		if err != nil {
			return newCommandResponse{Err: err}, nil
		}
		return newCommandResponse{Payload: payload}, nil
	}
}

// NextCommandRequest is a request to return the next command in a device queue
type nextCommandRequest struct {
	UDID string
}

// NextCommandResponse is a response for the next command
type nextCommandResponse struct {
	Payload []byte `json:"command_payload"`
	Total   int    `json:"total_payloads"`
	Err     error  `json:"error,omitempty"`
}

func makeNextCommandEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(nextCommandRequest)
		if req.UDID == "" {
			return nextCommandResponse{Err: ErrEmptyRequest}, nil
		}
		payload, total, err := svc.NextCommand(req.UDID)
		if err != nil {
			return nextCommandResponse{Err: err}, nil
		}
		return nextCommandResponse{Payload: payload, Total: total}, nil
	}
}

// deleteCommandRequest is a request to delete a command
type deleteCommandRequest struct {
	// device UDID
	UDID string
	// command UUID
	UUID string
}

// deleteCommandResponse is a response for a delete request
type deleteCommandResponse struct {
	Total int   `json:"remaining_payloads"`
	Err   error `json:"error,omitempty"`
}

func makeDeleteCommandEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteCommandRequest)
		if req.UDID == "" {
			return deleteCommandResponse{Err: ErrEmptyRequest}, nil
		}
		if req.UUID == "" {
			return deleteCommandResponse{Err: ErrEmptyRequest}, nil
		}
		total, err := svc.DeleteCommand(req.UDID, req.UUID)
		if err != nil {
			return deleteCommandResponse{Err: err}, nil
		}
		return deleteCommandResponse{Total: total}, nil
	}
}

type getCommandsRequest struct {
	UDID string
}

type getCommandsResponse struct {
	Commands []mdm.Payload `json:"commands,omitempty"`
	Err      error         `json:"error,omitempty"`
}

func makeGetCommandsEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getCommandsRequest)
		if req.UDID == "" {
			return getCommandsResponse{Err: ErrEmptyRequest}, nil
		}
		commands, err := svc.Commands(req.UDID)
		if err != nil {
			return getCommandsResponse{Err: err}, nil
		}
		return getCommandsResponse{Commands: commands}, nil
	}
}

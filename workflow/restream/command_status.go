package restream

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/gorilla/mux"
	"github.com/groob/plist"
	"github.com/micromdm/micromdm/mdm"
	"github.com/micromdm/micromdm/platform/command"
	"github.com/pkg/errors"
)

func (svc RestreamService) CommandStatus(ctx context.Context, id string) (Event, error) {
	event, err := svc.db.Event(ctx, id)
	return event, err
}

type commandStatusRequest struct {
	ID string
}

type commandStatusResponse struct {
	ID            string `json:"id,omitempty"`
	Status        string `json:"status,omitempty"`
	RequestPlist  []byte `json:"request_payload,omitempty"`
	ResponsePlist []byte `json:"response_payload,omitempty"`
	Err           error  `json:"error,omitempty"`
}

func (r commandStatusResponse) Failed() error { return r.Err }

func decodeCommandStatusRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return 0, errors.New("restream: bad route")
	}

	return commandStatusRequest{ID: id}, nil
}

func MakeCommandStatusEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(commandStatusRequest)

		event, err := svc.CommandStatus(ctx, req.ID)
		if err != nil {
			return commandStatusResponse{Err: err}, nil
		}

		// req
		var ev command.Event
		if err := command.UnmarshalEvent(event.Request, &ev); err != nil {
			return commandStatusResponse{Err: err}, nil
		}
		requestPlist, err := plist.MarshalIndent(ev.Payload, "\t")
		if err != nil {
			return commandStatusResponse{Err: err}, nil
		}

		// ack'd response
		var responsePlist []byte
		if event.Response != nil {
			var resp mdm.AcknowledgeEvent
			if err := mdm.UnmarshalAcknowledgeEvent(event.Response, &resp); err != nil {
				return commandStatusResponse{Err: err}, nil
			}
			responsePlist = resp.Raw
		}

		return commandStatusResponse{
			ID:            event.ID,
			Status:        event.Status,
			RequestPlist:  requestPlist,
			ResponsePlist: responsePlist,
		}, nil
	}
}

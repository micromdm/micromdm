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

func (svc *MDMService) Acknowledge(ctx context.Context, req AcknowledgeEvent) (payload []byte, err error) {
	msg, err := MarshalAcknowledgeEvent(&req)
	if err != nil {
		return nil, errors.Wrap(err, "marshal acknowledge response to proto")
	}

	if err := svc.pub.Publish(ctx, ConnectTopic, msg); err != nil {
		return nil, errors.Wrap(err, "publish connect Response on pubsub")
	}

	payload, err = svc.queue.Next(ctx, req.Response)
	return payload, errors.Wrap(err, "calling Next with mdm response")

}

type acknowledgeRequest struct {
	Event AcknowledgeEvent
}

type acknowledgeResponse struct {
	Payload []byte
	Err     error `plist:"error,omitempty"`
}

func (r acknowledgeResponse) Response() []byte { return r.Payload }
func (r acknowledgeResponse) Failed() error    { return r.Err }

func (d *requestDecoder) decodeAcknowledgeRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	body, _, err := d.readBody(r)
	if err != nil {
		return nil, errors.Wrap(err, "read acknowledge request body")
	}

	var res Response
	err = plist.Unmarshal(body, &res)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal MDM Response plist")
	}

	params := mux.Vars(r)

	event := AcknowledgeEvent{
		ID:       uuid.NewV4().String(),
		Time:     time.Now().UTC(),
		Response: res,
		Params:   params,
		Raw:      body,
	}
	req := acknowledgeRequest{Event: event}
	return req, nil
}

func MakeAcknowledgeEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(acknowledgeRequest)
		payload, err := svc.Acknowledge(ctx, req.Event)
		return acknowledgeResponse{Payload: payload, Err: err}, nil
	}
}

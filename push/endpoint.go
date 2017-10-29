package push

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	PushEndpoint endpoint.Endpoint
}

type pushRequest struct {
	UDID string
}

type pushResponse struct {
	Status string `json:"status,omitempty"`
	ID     string `json:"push_notification_id,omitempty"`
	Err    error  `json:"error,omitempty"`
}

func (r pushResponse) error() error { return r.Err }

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

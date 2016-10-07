package management

import (
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
)

type pushRequest struct {
	UDID string
}

type pushResponse struct {
	Status string `json:"status,omitempty"`
	ID     string `json:"push_notification_id,omitempty"`
	Err    error  `json:"error,omitempty"`
}

func (r pushResponse) error() error { return r.Err }

func makePushEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(pushRequest)
		id, err := svc.Push(req.UDID)
		if err != nil {
			return pushResponse{Err: err, Status: "failure"}, nil
		}
		return pushResponse{Status: "success", ID: id}, nil
	}
}

package user

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/micromdm/micromdm/pkg/httputil"
	"github.com/pkg/errors"
)

func (svc *UserService) ApplyUser(ctx context.Context, u User) (*User, error) {
	toSave := &u
	if u.UUID == "" { //newUser
		usr, err := NewFromRequest(u)
		if err != nil {
			return nil, errors.Wrap(err, "create user from request")
		}
		toSave = usr
	}
	err := svc.store.Save(toSave)
	return toSave, errors.Wrap(err, "apply user")
}

type applyUserRequest struct {
	User User `json:"user"`
}

type applyUserResponse struct {
	User User  `json:"user"`
	Err  error `json:"err,omitempty"`
}

func (r applyUserResponse) Failed() error { return r.Err }

func decodeApplyUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req applyUserRequest
	err := httputil.DecodeJSONRequest(r, &req)
	return req, err
}

func decodeApplyUserResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var resp applyUserResponse
	err := httputil.DecodeJSONResponse(r, &resp)
	return resp, err
}

func MakeApplyUserEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(applyUserRequest)
		u, err := svc.ApplyUser(ctx, req.User)
		return applyUserResponse{
			User: *u,
			Err:  err,
		}, nil
	}
}

func (e Endpoints) ApplyUser(ctx context.Context, u User) (*User, error) {
	request := applyUserRequest{
		User: u,
	}
	resp, err := e.ApplyUserEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	usr := resp.(applyUserResponse).User
	return &usr, resp.(applyUserResponse).Err
}

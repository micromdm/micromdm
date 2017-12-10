package user

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/micromdm/micromdm/pkg/httputil"
	"github.com/pkg/errors"
)

func (svc *UserService) ListUsers(ctx context.Context, opts ListUsersOption) ([]User, error) {
	u, err := svc.store.List()
	return u, errors.Wrap(err, "list users from api request")
}

type getUsersRequest struct{ Opts ListUsersOption }
type getUsersResponse struct {
	Users []User `json:"users"`
	Err   error  `json:"err,omitempty"`
}

func (r getUsersResponse) Failed() error { return r.Err }

func decodeListUsersRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var opts ListUsersOption
	if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
		return nil, err
	}
	req := getUsersRequest{Opts: opts}
	return req, nil
}

func decodeListUsersResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var resp getUsersResponse
	err := httputil.DecodeJSONResponse(r, &resp)
	return resp, err
}

func MakeListUsersEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getUsersRequest)
		users, err := svc.ListUsers(ctx, req.Opts)
		return getUsersResponse{
			Users: users,
			Err:   err,
		}, nil
	}
}

func (e Endpoints) ListUsers(ctx context.Context, opts ListUsersOption) ([]User, error) {
	request := getUsersRequest{opts}
	response, err := e.ListUsersEndpoint(ctx, request.Opts)
	if err != nil {
		return nil, err
	}
	return response.(getUsersResponse).Users, response.(getUsersResponse).Err
}

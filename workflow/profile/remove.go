package profile

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-kit/kit/endpoint"
	"github.com/pkg/errors"
	sq "gopkg.in/Masterminds/squirrel.v1"
)

type RemoveOption struct {
	FilterPayloadIdentifier string
	FilterID                []string
}

func (svc *ProfileService) Remove(ctx context.Context, o RemoveOption) error {
	if len(o.FilterID) == 0 {
		return nil
	}

	query, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Delete(tableName).
		Where(sq.Eq{"id": o.FilterID}).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "building sql for delete by id")
	}

	_, err = svc.db.ExecContext(ctx, query, args...)
	return errors.Wrap(err, "delete profile by id")
}

type removeRequest struct {
	Opts RemoveOption
}

type removeResponse struct {
	Err error `json:"err,omitempty"`
}

func (r removeResponse) Failed() error   { return r.Err }
func (r removeResponse) StatusCode() int { return http.StatusNoContent }

func decodeRemoveRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	q := r.URL.Query()
	req := removeRequest{
		Opts: RemoveOption{
			FilterPayloadIdentifier: q.Get("payload_identifier"),
			FilterID:                strings.Split(q.Get("id"), ","),
		},
	}
	return req, nil
}

func MakeRemoveProfilesEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(removeRequest)
		err := svc.Remove(ctx, req.Opts)
		return removeResponse{Err: err}, nil
	}
}

func (e Endpoints) Remove(ctx context.Context, opts RemoveOption) error {
	request := removeRequest{Opts: opts}
	response, err := e.RemoveProfilesEndpoint(ctx, request)
	if err != nil {
		return err
	}
	return response.(removeResponse).Err
}

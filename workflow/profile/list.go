package profile

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-kit/kit/endpoint"
	"github.com/pkg/errors"
	sq "gopkg.in/Masterminds/squirrel.v1"
)

type ListOption struct {
	FilterPayloadIdentifier []string
	FilterUUID              string
}

func (svc *ProfileService) List(ctx context.Context, o ListOption) ([]Profile, error) {
	var p []Profile
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(columns()...).
		From(tableName).
		OrderBy("id DESC")
	switch {
	case len(o.FilterPayloadIdentifier) != 0 && o.FilterUUID != "":
		builder.Where(
			sq.And{
				sq.Eq{"payload_identifier": o.FilterPayloadIdentifier},
				sq.Eq{"payload_uuid": o.FilterUUID},
			},
		)
	case len(o.FilterPayloadIdentifier) != 0 && o.FilterUUID == "":
		builder.Where(sq.Eq{"payload_identifier": o.FilterPayloadIdentifier})
	case len(o.FilterPayloadIdentifier) == 0 && o.FilterUUID != "":
		builder.Where(sq.Eq{"payload_uuid": o.FilterUUID})
	}
	query, args, err := builder.ToSql()
	if err != nil {
		return p, errors.Wrap(err, "building sql to find profile by payload_identifier")
	}

	err = svc.db.SelectContext(ctx, &p, query, args...)
	if err == nil && len(p) == 0 {
		p = []Profile{}
	}
	return p, errors.Wrap(err, "finding profiles by payload_identifier")
}

type listRequest struct {
	Opts ListOption
}

type listResponse struct {
	Profiles []Profile `json:"profiles"`
	Err      error     `json:"err,omitempty"`
}

func (r listResponse) Failed() error { return r.Err }

func decodeListRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	q := r.URL.Query()
	req := listRequest{
		Opts: ListOption{
			FilterUUID:              q.Get("uuid"),
			FilterPayloadIdentifier: strings.Split(q.Get("payload_identifier"), ","),
		},
	}
	return req, nil
}

func MakeListProfilesEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listRequest)
		profiles, err := svc.List(ctx, req.Opts)
		return listResponse{Profiles: profiles, Err: err}, nil
	}
}

func (e Endpoints) List(ctx context.Context, opts ListOption) ([]Profile, error) {
	request := listRequest{Opts: opts}
	response, err := e.ListProfilesEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.(listResponse).Profiles, response.(listResponse).Err
}

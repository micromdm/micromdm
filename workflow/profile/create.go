package profile

import (
	"context"
	"net/http"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/groob/plist"
	"github.com/pkg/errors"
	sq "gopkg.in/Masterminds/squirrel.v1"

	"github.com/micromdm/micromdm/pkg/httputil"
	"github.com/micromdm/micromdm/pkg/id"
)

func (svc *ProfileService) Create(ctx context.Context, payload []byte) (Profile, error) {
	p, err := create(ctx, payload)
	if err != nil {
		return p, errors.Wrap(err, "create profile for insert")
	}

	query, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert(tableName).
		Columns(columns()...).
		Values(
			p.ID,
			p.PayloadIdentifier,
			p.PayloadUUID,
			p.MobileconfigData,
			p.CreatedAt,
			p.UpdatedAt,
			p.CreatedBy,
			p.UpdatedBy,
		).
		ToSql()
	if err != nil {
		return p, errors.Wrap(err, "building insert profile query")
	}

	_, err = svc.db.ExecContext(ctx, query, args...)
	if isAlreadyExistsPG(err) {
		return p, alreadyExistsErr{uuid: p.PayloadUUID, identifier: p.PayloadIdentifier}
	}
	return p, errors.Wrap(err, "exec insert profile in pg")
}

func create(ctx context.Context, payload []byte) (Profile, error) {
	now, id, author := time.Now().UTC(), id.New(), "default_user"
	p := Profile{
		ID:               id,
		MobileconfigData: payload,
		CreatedAt:        now,
		UpdatedAt:        now,
		CreatedBy:        author,
		UpdatedBy:        author,
	}
	if err := plist.Unmarshal(payload, &p); err != nil {
		return p, errors.Wrap(err, "unmarshal mobileconfig to create profile")
	}
	return p, nil
}

type createRequest struct {
	MobileconfigData []byte `json:"mobileconfig_data"`
}

type createResponse struct {
	Profile Profile `json:"profile"`
	Err     error   `json:"err,omitempty"`
}

func (r createResponse) Failed() error   { return r.Err }
func (r createResponse) StatusCode() int { return http.StatusCreated }

func decodeCreateRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req createRequest
	err := httputil.DecodeJSONRequest(r, &req)
	return req, errors.Wrap(err, "decoding create workflow profile request")
}

func MakeCreateProfileEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createRequest)
		profile, err := svc.Create(ctx, req.MobileconfigData)
		if err != nil {
			return createResponse{Err: err}, nil
		}
		return createResponse{Profile: profile}, nil
	}
}

func (e Endpoints) Create(ctx context.Context, payload []byte) (Profile, error) {
	request := createRequest{MobileconfigData: payload}
	response, err := e.CreateProfileEndpoint(ctx, request.MobileconfigData)
	if err != nil {
		return Profile{}, err
	}
	return response.(createResponse).Profile, response.(createResponse).Err
}

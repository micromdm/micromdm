package profile

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	sq "gopkg.in/Masterminds/squirrel.v1"

	"github.com/micromdm/micromdm/pkg/httputil"
)

type Profile struct {
	ID                string    `db:"id"`
	PayloadIdentifier string    `db:"payload_identifier"`
	PayloadUUID       string    `db:"payload_uuid"`
	MobileconfigData  []byte    `db:"mobileconfig_data"`
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
	CreatedBy         string    `db:"created_by"`
	UpdatedBy         string    `db:"updated_by"`
}

func columns() []string {
	return []string{
		"id",
		"payload_identifier",
		"payload_uuid",
		"mobileconfig_data",
		"created_at",
		"updated_at",
		"created_by",
		"updated_by",
	}
}

type Service interface {
	Create(ctx context.Context, payload []byte) (Profile, error)
	List(ctx context.Context, options ListOption) ([]Profile, error)
	Remove(ctx context.Context, options RemoveOption) error
}

type Option func(*ProfileService)

func WithLogger(logger log.Logger) Option {
	return func(svc *ProfileService) {
		svc.logger = logger
	}
}

type ProfileService struct {
	db     *sqlx.DB
	logger log.Logger
}

func New(db *sqlx.DB, opts ...Option) *ProfileService {
	svc := ProfileService{
		db:     db,
		logger: log.NewNopLogger(),
	}

	for _, opt := range opts {
		opt(&svc)
	}

	return &svc
}

type Endpoints struct {
	CreateProfileEndpoint  endpoint.Endpoint
	ListProfilesEndpoint   endpoint.Endpoint
	RemoveProfilesEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(svc Service, outer endpoint.Middleware, others ...endpoint.Middleware) Endpoints {
	return Endpoints{
		CreateProfileEndpoint:  endpoint.Chain(outer, others...)(MakeCreateProfileEndpoint(svc)),
		ListProfilesEndpoint:   endpoint.Chain(outer, others...)(MakeListProfilesEndpoint(svc)),
		RemoveProfilesEndpoint: endpoint.Chain(outer, others...)(MakeRemoveProfilesEndpoint(svc)),
	}
}

func RegisterHTTPHandlers(r *mux.Router, e Endpoints, options ...http.ServerOption) {
	// POST     /v1/workflow/profiles		Create a new profile.

	r.Methods("POST").Path("/v1/workflow/profiles").Handler(http.NewServer(
		e.CreateProfileEndpoint,
		decodeCreateRequest,
		httputil.EncodeJSONResponse,
		options...,
	))

	r.Methods("GET").Path("/v1/workflow/profiles").Handler(http.NewServer(
		e.ListProfilesEndpoint,
		decodeListRequest,
		httputil.EncodeJSONResponse,
		options...,
	))

	r.Methods("DELETE").Path("/v1/workflow/profiles").Handler(http.NewServer(
		e.RemoveProfilesEndpoint,
		decodeRemoveRequest,
		httputil.EncodeJSONResponse,
		options...,
	))
}

const (
	tableName     = "wf_profiles"
	defaultAuthor = "default_user"
)

func isAlreadyExistsPG(err error) bool {
	pqErr, ok := errors.Cause(err).(*pq.Error)
	return ok && pqErr.Code.Name() == "unique_violation"
}

func (svc *ProfileService) Find(ctx context.Context, id string) (Profile, error) {
	var p Profile
	query, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(columns()...).
		From(tableName).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return p, errors.Wrap(err, "building sql to find profile")
	}

	err = svc.db.QueryRowxContext(ctx, query, args...).StructScan(&p)
	if errors.Cause(err) == sql.ErrNoRows {
		return p, notFoundErr{}
	}
	return p, errors.Wrap(err, "finding profile by id")
}

func (svc *ProfileService) FindByPayloadUUID(ctx context.Context, uuid string) (Profile, error) {
	var p Profile
	query, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(columns()...).
		From(tableName).
		Where(sq.Eq{"payload_uuid": uuid}).
		ToSql()
	if err != nil {
		return p, errors.Wrap(err, "building sql to find profile by uuid")
	}

	err = svc.db.QueryRowxContext(ctx, query, args...).StructScan(&p)
	if errors.Cause(err) == sql.ErrNoRows {
		return p, notFoundErr{}
	}
	return p, errors.Wrap(err, "finding profile by payload_uuid")
}

func (svc *ProfileService) FindByPayloadIdentifier(ctx context.Context, identifier string) ([]Profile, error) {
	var p []Profile
	query, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(columns()...).
		From(tableName).
		Where(sq.Eq{"payload_identifier": identifier}).
		ToSql()
	if err != nil {
		return p, errors.Wrap(err, "building sql to find profile by payload_identifier")
	}

	err = svc.db.SelectContext(ctx, &p, query, args...)
	if errors.Cause(err) == sql.ErrNoRows {
		return p, notFoundErr{}
	}
	return p, errors.Wrap(err, "finding profiles by payload_identifier")
}

type notFoundErr struct{}

func (e notFoundErr) Error() string  { return "profile not found" }
func (e notFoundErr) NotFound() bool { return true }

type alreadyExistsErr struct{ uuid, identifier string }

func (e alreadyExistsErr) Error() string {
	return fmt.Sprintf("profile %q with uuid %q already exists", e.identifier, e.uuid)
}
func (e alreadyExistsErr) AlreadyExists() bool { return true }

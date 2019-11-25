package restream

import "context"

type Service interface {
	CommandStatus(ctx context.Context, id string) (Event, error)
}

type RestreamService struct {
	db Store
}

func NewService(db Store) RestreamService {
	return RestreamService{db: db}
}

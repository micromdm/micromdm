package blueprint

import (
	"context"
)

type GetBlueprintsOption struct {
	FilterName string
}

type Service interface {
	ApplyBlueprint(ctx context.Context, bp *Blueprint) error
	GetBlueprints(ctx context.Context, opt GetBlueprintsOption) ([]Blueprint, error)
	RemoveBlueprints(ctx context.Context, names []string) error
}

type Store interface {
	Save(*Blueprint) error
	BlueprintByName(name string) (*Blueprint, error)
	List() ([]Blueprint, error)
	Delete(string) error
}

type BlueprintService struct {
	store Store
}

func New(store Store) *BlueprintService {
	return &BlueprintService{store: store}
}

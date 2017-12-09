package remove

import (
	"context"

	"github.com/micromdm/micromdm/platform/remove"
)

type Service interface {
	UnblockDevice(ctx context.Context, udid string) error
}

type RemoveService struct {
	*remove.RemoveService
}

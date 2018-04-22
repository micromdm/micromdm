package checkin

import (
	"github.com/micromdm/mdm"
	"golang.org/x/net/context"
)

// Service defines methods for and MDM Check-in service.
type Service interface {
	Authenticate(ctx context.Context, cmd mdm.CheckinCommand, id string) error
	TokenUpdate(ctx context.Context, cmd mdm.CheckinCommand, id string) error
	CheckOut(ctx context.Context, cmd mdm.CheckinCommand, id string) error
}

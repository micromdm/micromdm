package checkin

import (
	"github.com/micromdm/mdm"
	"golang.org/x/net/context"
)

// Service defines methods for and MDM Check-in service.
type Service interface {
	Authenticate(ctx context.Context, cmd mdm.CheckinCommand) error
	TokenUpdate(ctx context.Context, cmd mdm.CheckinCommand) error
	CheckOut(ctx context.Context, cmd mdm.CheckinCommand) error
}

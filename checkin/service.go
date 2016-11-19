package checkin

import "github.com/micromdm/mdm"

// Service defines methods for and MDM Checkin service
type Service interface {
	Authenticate(mdm.CheckinCommand) error
	TokenUpdate(mdm.CheckinCommand) error
	Checkout(mdm.CheckinCommand) error
}

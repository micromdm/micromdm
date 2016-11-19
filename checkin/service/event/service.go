// Package nsq implements the checkin.Service by writing checkin events to
// an event publisher.
package event

import (
	"github.com/micromdm/mdm"
	"github.com/micromdm/micromdm/event"
)

// Service implements the checkin.Service.
type Service struct {
	publisher *event.Publisher
}

func NewCheckinService(publisher *event.Publisher) *Service {
	return &Service{
		publisher: publisher,
	}
}

func (svc *Service) Authenticate(command mdm.CheckinCommand) error {
	return svc.createAndPublish("mdm.Authenticate", command)
}

func (svc *Service) TokenUpdate(command mdm.CheckinCommand) error {
	return svc.createAndPublish("mdm.TokenUpdate", command)
}

func (svc *Service) Checkout(command mdm.CheckinCommand) error {
	return svc.createAndPublish("mdm.Checkout", command)
}

func (svc *Service) createAndPublish(topic string, command mdm.CheckinCommand) error {
	event, err := event.CreateCheckinEvent(command)
	if err != nil {
		return err
	}
	return svc.publisher.Publish(topic, event)
}

package event

import (
	"time"

	"github.com/groob/plist"
	"github.com/micromdm/mdm"
	nsq "github.com/nsqio/go-nsq"
	uuid "github.com/satori/go.uuid"
)

type CheckinEvent struct {
	ID    string
	Time  int64
	Event mdm.CheckinCommand
}

func CreateCheckinEvent(command mdm.CheckinCommand) ([]byte, error) {
	event := CheckinEvent{
		ID:    uuid.NewV4().String(),
		Time:  time.Now().UnixNano(),
		Event: command,
	}
	return plist.Marshal(&event)
}

type Publisher struct {
	Producer *nsq.Producer
}

func (p *Publisher) Publish(topic string, event []byte) error {
	return p.Producer.Publish(topic, event)
}

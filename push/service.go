package push

import (
	"encoding/json"
	"fmt"

	"github.com/RobotsAndPencils/buford/payload"
	"github.com/RobotsAndPencils/buford/push"
	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/micromdm/micromdm/pubsub"
	"github.com/micromdm/micromdm/queue"
)

type Push struct {
	db      *DB
	pushsvc *push.Service
}

func New(db *DB, push *push.Service, sub pubsub.Subscriber) (*Push, error) {
	pushSvc := Push{db, push}
	if err := pushSvc.startQueuedSubscriber(push, sub); err != nil {
		return &pushSvc, err
	}
	return &pushSvc, nil
}

func (svc *Push) startQueuedSubscriber(push *push.Service, sub pubsub.Subscriber) error {
	commandQueuedEvents, err := sub.Subscribe("push-info", queue.CommandQueuedTopic)
	if err != nil {
		return errors.Wrapf(err,
			"subscribing push to %s topic", queue.CommandQueuedTopic)
	}
	go func() {
		for {
			select {
			case event := <-commandQueuedEvents:
				cq, err := queue.UnmarshalQueuedCommand(event.Message)
				if err != nil {
					fmt.Println(err)
				}
				_, err = svc.Push(nil, cq.DeviceUDID)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}()

	return nil
}

func (svc *Push) Push(ctx context.Context, deviceUDID string) (string, error) {
	info, err := svc.db.PushInfo(deviceUDID)
	if err != nil {
		return "", errors.Wrap(err, "retrieving PushInfo by UDID")
	}

	p := payload.MDM{Token: info.PushMagic}
	valid := push.IsDeviceTokenValid(info.Token)
	if !valid {
		return "", errors.New("invalid push token")
	}
	jsonPayload, err := json.Marshal(p)
	if err != nil {
		return "", errors.Wrap(err, "marshalling push notification payload")
	}
	return svc.pushsvc.Push(info.Token, nil, jsonPayload)
}

package inmem

import (
	"container/list"
	"context"

	"github.com/micromdm/micromdm/mdm"
	"github.com/micromdm/micromdm/platform/command"
	"github.com/micromdm/micromdm/platform/pubsub"
	boltqueue "github.com/micromdm/micromdm/platform/queue"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/groob/plist"
)

// QueueInMem represents an in-memory command queue
type QueueInMem struct {
	logger log.Logger
	queue  map[string]*list.List
}

type queuedCommand struct {
	uuid    string
	payload []byte
	notNow  bool
}

// New creates a new in-memory command queue
func New(pubsub pubsub.PublishSubscriber, logger log.Logger) *QueueInMem {
	q := &QueueInMem{
		logger: logger,
		queue:  make(map[string]*list.List),
	}
	q.startPolling(pubsub)
	return q
}

func (q *QueueInMem) clearList(udid string) {
	delete(q.queue, udid)
	return
}

func (q *QueueInMem) getList(udid string) *list.List {
	if _, ok := q.queue[udid]; !ok {
		q.queue[udid] = list.New()
	}
	return q.queue[udid]
}

func (q *QueueInMem) enqueue(l *list.List, uuid string, payload []byte) {
	l.PushBack(&queuedCommand{
		uuid:    uuid,
		payload: payload,
	})
}

func (q *QueueInMem) findCommandByUUID(l *list.List, uuid string) (*queuedCommand, *list.Element) {
	for e := l.Front(); e != nil; e = e.Next() {
		qCmd := e.Value.(*queuedCommand)
		if qCmd.uuid == uuid {
			return qCmd, e
		}
	}
	return nil, nil
}

func (q *QueueInMem) nextCommandPayload(l *list.List, skipNotNow bool) []byte {
	for e := l.Front(); e != nil; e = e.Next() {
		qCmd := e.Value.(*queuedCommand)
		if !(skipNotNow && qCmd.notNow) {
			return qCmd.payload
		}
	}
	return nil
}

// Next delivers the next command from the command queue for the enrollment in resp
func (q *QueueInMem) Next(_ context.Context, resp mdm.Response) ([]byte, error) {
	udid := resp.UDID
	if resp.UserID != nil {
		udid = *resp.UserID
	}
	if resp.EnrollmentID != nil {
		udid = *resp.EnrollmentID
	}

	l := q.getList(udid)

	switch resp.Status {
	case "NotNow":
		qCmd, _ := q.findCommandByUUID(l, resp.CommandUUID)
		qCmd.notNow = true
	case "Acknowledged", "Error", "CommandFormatError":
		_, e := q.findCommandByUUID(l, resp.CommandUUID)
		if e != nil {
			l.Remove(e)
			if l.Len() == 0 {
				q.clearList(udid)
			}
		}
	}

	cmdBytes := q.nextCommandPayload(l, resp.Status == "NotNow")

	return cmdBytes, nil
}

// Clear clears a command queue for the enrollment in event
func (q *QueueInMem) Clear(_ context.Context, event mdm.CheckinEvent) error {
	udid := event.Command.UDID
	if event.Command.UserID != "" {
		udid = event.Command.UserID
	}
	if event.Command.EnrollmentID != "" {
		udid = event.Command.EnrollmentID
	}

	q.clearList(udid)
	return nil
}

func (q *QueueInMem) startPolling(pubsub pubsub.PublishSubscriber) error {
	events, err := pubsub.Subscribe(context.TODO(), "command-queue", command.CommandTopic)
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case event := <-events:
				var cmdEvent command.Event
				if err := command.UnmarshalEvent(event.Message, &cmdEvent); err != nil {
					level.Info(q.logger).Log(
						"msg", "unmarshal command event from pubsub",
						"err", err,
					)
					continue
				}
				rawCmdPlist, err := plist.Marshal(cmdEvent.Payload)
				if err != nil {
					level.Info(q.logger).Log(
						"msg", "marshal command plist",
						"err", err,
					)
					continue
				}
				q.enqueue(
					q.getList(cmdEvent.DeviceUDID),
					cmdEvent.Payload.CommandUUID,
					rawCmdPlist,
				)
				level.Info(q.logger).Log(
					"msg", "queued command for device",
					"device_udid", cmdEvent.DeviceUDID,
					"command_uuid", cmdEvent.Payload.CommandUUID,
					"request_type", cmdEvent.Payload.Command.RequestType,
				)

				err = boltqueue.PublishCommandQueued(pubsub, cmdEvent.DeviceUDID, cmdEvent.Payload.CommandUUID)
				if err != nil {
					level.Info(q.logger).Log(
						"msg", "publish command to queued topic",
						"err", err,
					)
				}
			}
		}
	}()
	return nil
}

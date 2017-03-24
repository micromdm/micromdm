package queue

import (
	"github.com/gogo/protobuf/proto"
	"github.com/micromdm/micromdm/queue/internal/commandqueuedproto"
)

func MarshalQueuedCommand(udid, uuid string) ([]byte, error) {
	return proto.Marshal(&commandqueued.CommandQueued{
		DeviceUdid:  udid,
		CommandUuid: uuid,
	})
}

func UnmarshalQueuedCommand(data []byte) (string, string, error) {
	cmdQueued := commandqueued.CommandQueued{}
	if err := proto.Unmarshal(data, &cmdQueued); err != nil {
		return cmdQueued.DeviceUdid, cmdQueued.CommandUuid, err
	}
	return cmdQueued.DeviceUdid, cmdQueued.CommandUuid, nil
}

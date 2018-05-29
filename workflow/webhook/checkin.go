package webhook

import (
	"github.com/pkg/errors"

	"github.com/micromdm/micromdm/mdm"
)

type CheckinEvent struct {
	UDID       string            `json:"udid"`
	Params     map[string]string `json:"url_params"`
	RawPayload []byte            `json:"raw_payload"`
}

func checkinEvent(topic string, data []byte) (*Event, error) {
	var ev mdm.CheckinEvent
	if err := mdm.UnmarshalCheckinEvent(data, &ev); err != nil {
		return nil, errors.Wrap(err, "unmarshal checkin event for webhook")
	}

	webhookEvent := Event{
		Topic:     topic,
		EventID:   ev.ID,
		CreatedAt: ev.Time,

		CheckinEvent: &CheckinEvent{
			UDID:       ev.Command.UDID,
			Params:     ev.Params,
			RawPayload: ev.Raw,
		},
	}

	return &webhookEvent, nil
}

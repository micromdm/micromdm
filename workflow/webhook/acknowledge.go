package webhook

import (
	"github.com/pkg/errors"

	"github.com/micromdm/micromdm/mdm"
)

type AcknowledgeEvent struct {
	UDID        string            `json:"udid"`
	Status      string            `json:"status"`
	CommandUUID string            `json:"command_uuid"`
	Params      map[string]string `json:"url_params"`
	RawPayload  []byte            `json:"raw_payload"`
}

func acknowledgeEvent(topic string, data []byte) (*Event, error) {
	var ev mdm.AcknowledgeEvent
	if err := mdm.UnmarshalAcknowledgeEvent(data, &ev); err != nil {
		return nil, errors.Wrap(err, "unmarshal acknowledge event for webhook")
	}
	webhookEvent := Event{
		Topic:     topic,
		EventID:   ev.ID,
		CreatedAt: ev.Time,

		AcknowledgeEvent: &AcknowledgeEvent{
			UDID:        ev.Response.UDID,
			Status:      ev.Response.Status,
			CommandUUID: ev.Response.CommandUUID,
			Params:      ev.Params,
			RawPayload:  ev.Raw,
		},
	}

	return &webhookEvent, nil
}

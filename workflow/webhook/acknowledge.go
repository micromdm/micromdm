package webhook

import (
	"github.com/pkg/errors"

	"github.com/micromdm/micromdm/mdm"
)

type AcknowledgeEvent struct {
	UDID         string            `json:"udid,omitempty"`
	EnrollmentID string            `json:"enrollment_id,omitempty"`
	Status       string            `json:"status"`
	CommandUUID  string            `json:"command_uuid,omitempty"`
	Params       map[string]string `json:"url_params,omitempty"`
	RawPayload   []byte            `json:"raw_payload"`
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
	if ev.Response.EnrollmentID != nil {
		webhookEvent.AcknowledgeEvent.EnrollmentID = *ev.Response.EnrollmentID
	}

	return &webhookEvent, nil
}

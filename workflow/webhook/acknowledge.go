package webhook

import (
	"context"

	"github.com/go-kit/kit/log/level"
	"github.com/micromdm/micromdm/mdm"
)

type AcknowledgeEvent struct {
	UDID        string            `json:"udid"`
	Status      string            `json:"status"`
	CommandUUID string            `json:"command_uuid"`
	Params      map[string]string `json:"url_params"`
	RawPayload  []byte            `json:"raw_payload"`
}

func (w *Worker) acknowledgeEvent(ctx context.Context, topic string, data []byte) {
	var ev mdm.AcknowledgeEvent
	if err := mdm.UnmarshalAcknowledgeEvent(data, &ev); err != nil {
		level.Info(w.logger).Log(
			"msg", "unmarshal pubsub event",
			"err", err,
			"topic", topic,
		)
		return
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

	w.post(ctx, &webhookEvent)
	return
}

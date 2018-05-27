package webhook

import (
	"context"

	"github.com/go-kit/kit/log/level"
	"github.com/micromdm/micromdm/mdm"
)

type CheckinEvent struct {
	UDID       string            `json:"udid"`
	Params     map[string]string `json:"url_params"`
	RawPayload []byte            `json:"raw_payload"`
}

func (w *Worker) checkinEvent(ctx context.Context, topic string, data []byte) {
	var ev mdm.CheckinEvent
	if err := mdm.UnmarshalCheckinEvent(data, &ev); err != nil {
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

		CheckinEvent: &CheckinEvent{
			UDID:       ev.Command.UDID,
			Params:     ev.Params,
			RawPayload: ev.Raw,
		},
	}

	w.post(ctx, &webhookEvent)
	return
}

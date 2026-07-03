package mdm

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/micromdm/micromdm/pkg/httputil"
	"github.com/micromdm/micromdm/platform/pubsub"
	"github.com/pkg/errors"
)

// RegisterMigrationSyncHandler registers a Basic-Auth-protected endpoint that publishes a
// checkin event directly onto the pubsub bus, bypassing both Mdm-Signature verification and
// the cert-auth middleware around the real Service.Checkin path (which panics on a missing
// device cert — one this route never has).
//
// It exists only for the NanoMDM migration: a dual-writer (mdmdirector) mirrors
// Authenticate/TokenUpdate/CheckOut events from NanoMDM into this MicroMDM's device store
// (via the same workers a real checkin feeds), keeping it fresh as a rollback target.
// Gated behind the -migration flag
func RegisterMigrationSyncHandler(r *mux.Router, pub pubsub.Publisher, username, apiKey string) {
	r.Methods(http.MethodPut).Path("/mdm/migration-checkin").Handler(
		httputil.RequireBasicAuth(migrationCheckinHandler(pub), username, apiKey, "micromdm"),
	)
}

func migrationCheckinHandler(pub pubsub.Publisher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var cmd CheckinCommand
		body, err := mdmRequestBody(r, &cmd)
		if err != nil {
			http.Error(w, errors.Wrap(err, "read migration-checkin body").Error(), http.StatusBadRequest)
			return
		}

		event := CheckinEvent{
			ID:      uuid.New().String(),
			Time:    time.Now().UTC(),
			Command: cmd,
			Raw:     body,
		}

		msg, err := MarshalCheckinEvent(&event)
		if err != nil {
			http.Error(w, errors.Wrap(err, "marshal migration checkin event").Error(), http.StatusInternalServerError)
			return
		}

		topic, err := topicFromMessage(cmd.MessageType)
		if err != nil {
			http.Error(w, errors.Wrap(err, "get migration checkin topic").Error(), http.StatusBadRequest)
			return
		}

		if err := pub.Publish(context.Background(), topic, msg); err != nil {
			http.Error(w, errors.Wrap(err, "publish migration checkin").Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

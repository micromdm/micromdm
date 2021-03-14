package inmem

import (
	"fmt"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/micromdm/micromdm/mdm"
	"github.com/micromdm/micromdm/platform/pubsub/inmem"
)

func TestQueue(t *testing.T) {
	q := New(inmem.NewPubSub(), log.NewNopLogger())
	udid := "ABCD-EFGH"
	l := q.getList(udid)

	q.enqueue(l, "CMD-001", []byte("CMD-001"))
	q.enqueue(l, "CMD-002", []byte("CMD-002"))
	q.enqueue(l, "CMD-003", []byte("CMD-003"))

	for i, test := range []struct {
		nextUUID        string
		nextStatus      string
		expectedContent string
		expectedLength  int
	}{
		{"", "Idle", "CMD-001", 3},
		{"CMD-001", "Acknowledged", "CMD-002", 2},
		{"CMD-002", "NotNow", "CMD-003", 2},
		{"CMD-003", "NotNow", "", 2},
		{"", "Idle", "CMD-002", 2},
		{"CMD-002", "Acknowledged", "CMD-003", 1},
		{"CMD-003", "Acknowledged", "", 0},
		{"", "Idle", "", 0},
	} {
		t.Run(fmt.Sprintf("QueueTest%d-%s", i, test.nextStatus), func(t *testing.T) {
			resp, err := q.Next(nil, mdm.Response{
				UDID:        udid,
				CommandUUID: test.nextUUID,
				Status:      test.nextStatus,
			})
			if err != nil {
				t.Fatal(err)
			}
			if have, want, msg := l.Len(), test.expectedLength, "queue length"; have != want {
				t.Errorf("%v; have: %v, want: %v", msg, have, want)
			}
			if have, want, msg := string(resp), test.expectedContent, "response content"; have != want {
				t.Errorf("%v; have: %v, want: %v", msg, have, want)
			}
		})
	}
}

package inmem

import (
	"fmt"
	"testing"
	"time"
)

func TestPubSub(t *testing.T) {
	inmem := NewPubSub()
	tests := []string{"a", "b", "c"}
	for _, tt := range tests {
		if err := inmem.Publish(tt, []byte(tt+tt+tt)); err != nil {
			t.Fatal(err)
		}
		if err := inmem.Publish(tt, []byte(tt+tt)); err != nil {
			t.Fatal(err)
		}
	}

	subA, err := inmem.Subscribe("asub", "a")
	if err != nil {
		t.Fatal(err)
	}
	subA1, err := inmem.Subscribe("asub1", "a")
	if err != nil {
		t.Fatal(err)
	}

	subB, err := inmem.Subscribe("bsub", "b")
	if err != nil {
		t.Fatal(err)
	}

	for {
		select {
		case s := <-subA:
			fmt.Println("asub:", s)
		case s := <-subA1:
			fmt.Println("asub1:", s)
		case s := <-subB:
			fmt.Println("bsub:", s)
		case <-time.After(10 * time.Millisecond):
			return
		}
	}
}

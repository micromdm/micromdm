package vpp

import (
	"context"
	"encoding/json"
	"log"

	"github.com/micromdm/micromdm/platform/config"
	"github.com/micromdm/micromdm/platform/pubsub"
)

func (svc *VPPService) watchTokenUpdates(pubsub pubsub.Subscriber, serverURL string) error {
	tokenAdded, err := pubsub.Subscribe(context.TODO(), "list-token-events", config.VPPTokenTopic)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case event := <-tokenAdded:
				var token config.VPPToken
				if err := json.Unmarshal(event.Message, &token); err != nil {
					log.Printf("unmarshalling tokenAdded to token: %s\n", err)
					continue
				}

				client, err := token.Client(serverURL)
				if err != nil {
					log.Printf("creating new VPP client: %s\n", err)
					continue
				}

				svc.mtx.Lock()
				svc.client = client
				svc.mtx.Unlock()
			}
		}
	}()

	return nil
}

package apply

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"sync"

	"github.com/micromdm/dep"

	"github.com/micromdm/micromdm/platform/appstore"
	"github.com/micromdm/micromdm/platform/config"
	"github.com/micromdm/micromdm/platform/pubsub"
)

type Service interface {
	UploadApp(ctx context.Context, manifestName string, manifest io.Reader, pkgName string, pkg io.Reader) error
	DEPService
}

type ApplyService struct {
	mtx       sync.RWMutex
	DEPClient dep.Client

	Apps appstore.AppStore
}

func (svc *ApplyService) UploadApp(ctx context.Context, manifestName string, manifest io.Reader, pkgName string, pkg io.Reader) error {
	if manifestName != "" {
		if err := svc.Apps.SaveFile(manifestName, manifest); err != nil {
			return err
		}
	}

	if pkgName != "" {
		if err := svc.Apps.SaveFile(pkgName, pkg); err != nil {
			return err
		}
	}

	return nil
}

func (svc *ApplyService) WatchTokenUpdates(pubsub pubsub.Subscriber) error {
	tokenAdded, err := pubsub.Subscribe(context.TODO(), "apply-token-events", config.DEPTokenTopic)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case event := <-tokenAdded:
				var token config.DEPToken
				if err := json.Unmarshal(event.Message, &token); err != nil {
					log.Printf("unmarshalling tokenAdded to token: %s\n", err)
					continue
				}

				client, err := token.Client()
				if err != nil {
					log.Printf("creating new DEP client: %s\n", err)
					continue
				}

				svc.mtx.Lock()
				svc.DEPClient = client
				svc.mtx.Unlock()
			}
		}
	}()

	return nil
}

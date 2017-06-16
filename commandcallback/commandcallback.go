package commandcallback

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"github.com/micromdm/micromdm/pubsub"
)

const CommandResponseJSONTopic = "mdm.CommandResponseJSON"

func NewCommandCallback(sub pubsub.Subscriber, url string) error {
	cmdRespJSONEvents, err := sub.Subscribe(context.TODO(), "CommandResponseJSON", CommandResponseJSONTopic)
	if err != nil {
		return errors.Wrapf(err,
			"subscribing to %s topic", CommandResponseJSONTopic)
	}

	go func() {
		for {
			select {
			case e := <-cmdRespJSONEvents:
				go func() {
					c := http.Client{}
					_, err := c.Post(url, "application/json", bytes.NewBuffer(e.Message))
					if err != nil {
						fmt.Printf("error sending JSON command response: %s\n", err)
					}
				}()
			}
		}
	}()

	return nil
}

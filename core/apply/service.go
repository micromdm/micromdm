package apply

import (
	"context"
	"log"
	"sync"

	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/textproto"

	"github.com/fullsailor/pkcs7"

	"github.com/micromdm/dep"
	"github.com/micromdm/micromdm/blueprint"
	"github.com/micromdm/micromdm/deptoken"
	"github.com/micromdm/micromdm/profile"
	"github.com/micromdm/micromdm/pubsub"
)

type Service interface {
	ApplyBlueprint(ctx context.Context, bp *blueprint.Blueprint) error
	ApplyDEPToken(ctx context.Context, P7MContent []byte) error
	ApplyProfile(ctx context.Context, p *profile.Profile) error
	DEPService
}

type ApplyService struct {
	mtx       sync.RWMutex
	DEPClient dep.Client

	Blueprints *blueprint.DB
	Profiles   *profile.DB
	Tokens     *deptoken.DB
}

func (svc *ApplyService) WatchTokenUpdates(pubsub pubsub.Subscriber) error {
	tokenAdded, err := pubsub.Subscribe("apply-token-events", deptoken.DEPTokenTopic)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case event := <-tokenAdded:
				var token deptoken.DEPToken
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

func (svc *ApplyService) ApplyBlueprint(ctx context.Context, bp *blueprint.Blueprint) error {
	return svc.Blueprints.Save(bp)
}

// unwrapSMIME removes the S/MIME-like wrapper around raw CMS/PKCS7 data
func unwrapSMIME(smime []byte) ([]byte, error) {
	tr := textproto.NewReader(bufio.NewReader(bytes.NewReader(smime)))
	if _, err := tr.ReadMIMEHeader(); err != nil {
		return nil, err
	}
	dec := base64.NewDecoder(base64.StdEncoding, tr.DotReader())
	buf := new(bytes.Buffer)
	io.Copy(buf, dec)
	return buf.Bytes(), nil
}

// unwrapTokenJSON removes the MIME-like headers and text surrounding the DEP token JSON
func unwrapTokenJSON(wrapped []byte) ([]byte, error) {
	tr := textproto.NewReader(bufio.NewReader(bytes.NewReader(wrapped)))
	if _, err := tr.ReadMIMEHeader(); err != nil {
		return nil, err
	}
	tokenJSON := new(bytes.Buffer)
	for {
		line, err := tr.ReadLineBytes()
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		line = bytes.Trim(line, "-----BEGIN MESSAGE-----")
		line = bytes.Trim(line, "-----END MESSAGE-----")
		if _, err := tokenJSON.Write(line); err != nil {
			return nil, err
		}
	}
	return tokenJSON.Bytes(), nil
}

func (svc *ApplyService) ApplyDEPToken(ctx context.Context, P7MContent []byte) error {
	unwrapped, err := unwrapSMIME(P7MContent)
	if err != nil {
		return err
	}
	key, cert, err := svc.Tokens.DEPKeypair()
	if err != nil {
		return err
	}
	p7, err := pkcs7.Parse(unwrapped)
	if err != nil {
		return err
	}
	decrypted, err := p7.Decrypt(cert, key)
	if err != nil {
		return err
	}
	tokenJSON, err := unwrapTokenJSON(decrypted)
	if err != nil {
		return err
	}
	var depToken deptoken.DEPToken
	err = json.Unmarshal(tokenJSON, &depToken)
	if err != nil {
		return err
	}
	err = svc.Tokens.AddToken(depToken.ConsumerKey, tokenJSON)
	if err != nil {
		return err
	}
	log.Println("stored DEP token with ck", depToken.ConsumerKey)
	return nil
}

func (svc *ApplyService) ApplyProfile(ctx context.Context, p *profile.Profile) error {
	return svc.Profiles.Save(p)
}

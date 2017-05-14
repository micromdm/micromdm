package apply

import (
	"context"
	"fmt"

	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/fullsailor/pkcs7"
	"io"
	"net/textproto"

	"github.com/boltdb/bolt"
	"github.com/micromdm/micromdm/blueprint"
	"github.com/micromdm/micromdm/core/list"
	"github.com/micromdm/micromdm/profile"
)

type Service interface {
	ApplyBlueprint(ctx context.Context, bp *blueprint.Blueprint) error
	ApplyDEPToken(ctx context.Context, P7MContent []byte) error
	ApplyProfile(ctx context.Context, p *profile.Profile) error
}

type ApplyService struct {
	Blueprints *blueprint.DB
	Profiles   *profile.DB
	DB         *bolt.DB // TODO: replace with reference to DEP token svc/pkg
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

// TODO: move into seperate svc/pkg
const (
	depTokenBucket = "mdm.DEPToken"
)

func PutDEPToken(db *bolt.DB, consumerKey string, json []byte) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(depTokenBucket))
		if err != nil {
			return err
		}
		return b.Put([]byte(consumerKey), json)
	})
	return err
}

func (svc *ApplyService) ApplyDEPToken(ctx context.Context, P7MContent []byte) error {
	unwrapped, err := unwrapSMIME(P7MContent)
	if err != nil {
		return err
	}
	key, cert, err := list.GetDEPKeypair(svc.DB)
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
	var depToken list.DEPToken
	err = json.Unmarshal(tokenJSON, &depToken)
	if err != nil {
		return err
	}
	err = PutDEPToken(svc.DB, depToken.ConsumerKey, tokenJSON)
	if err != nil {
		return err
	}
	fmt.Println("stored DEP token with ck", depToken.ConsumerKey)
	return nil
}

func (svc *ApplyService) ApplyProfile(ctx context.Context, p *profile.Profile) error {
	return svc.Profiles.Save(p)
}

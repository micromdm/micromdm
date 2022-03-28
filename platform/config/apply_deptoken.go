package config

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/textproto"

	"github.com/fullsailor/pkcs7"
	"github.com/go-kit/kit/endpoint"
	"github.com/micromdm/micromdm/pkg/httputil"
)

func (svc *ConfigService) ApplyDEPToken(ctx context.Context, P7MContent []byte) error {
	unwrapped, err := unwrapSMIME(P7MContent)
	if err != nil {
		return err
	}
	key, cert, err := svc.store.DEPKeypair()
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
	var depToken DEPToken
	err = json.Unmarshal(tokenJSON, &depToken)
	if err != nil {
		return err
	}
	err = svc.store.AddDEPToken(depToken.ConsumerKey, tokenJSON)
	if err != nil {
		return err
	}
	fmt.Println("stored DEP token with ck", depToken.ConsumerKey)
	return nil
}

type applyDEPTokenRequest struct {
	P7MContent []byte `json:"p7m_content"`
}

type applyDEPTokenResponse struct {
	Err error `json:"err,omitempty"`
}

func (r applyDEPTokenResponse) Failed() error { return r.Err }

func decodeApplyDEPTokensRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req applyDEPTokenRequest
	err := httputil.DecodeJSONRequest(r, &req)
	return req, err
}

func decodeApplyDEPTokensResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var resp applyDEPTokenResponse
	err := httputil.DecodeJSONResponse(r, &resp)
	return resp, err
}

func MakeApplyDEPTokensEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(applyDEPTokenRequest)
		err = svc.ApplyDEPToken(ctx, req.P7MContent)
		return applyDEPTokenResponse{
			Err: err,
		}, nil
	}
}

func (e Endpoints) ApplyDEPToken(ctx context.Context, P7MContent []byte) error {
	req := applyDEPTokenRequest{P7MContent: P7MContent}
	resp, err := e.ApplyDEPTokensEndpoint(ctx, req)
	if err != nil {
		return err
	}
	return resp.(applyDEPTokenResponse).Err
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

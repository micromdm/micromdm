package scepserver

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"golang.org/x/net/context"
)

// ServiceHandler is an HTTP Handler for a SCEP endpoint.
func ServiceHandler(ctx context.Context, svc Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerBefore(updateContext),
	}

	scepHandler := kithttp.NewServer(
		ctx,
		makeSCEPEndpoint(svc),
		decodeSCEPRequest,
		encodeSCEPResponse,
		opts...,
	)

	mux := http.NewServeMux()
	mux.Handle("/scep", scepHandler)
	return mux
}

func updateContext(ctx context.Context, r *http.Request) context.Context {
	q := r.URL.Query()
	if _, ok := q["operation"]; ok {
		ctx = context.WithValue(ctx, "operation", q.Get("operation"))
	}
	return ctx
}

// EncodeSCEPRequest encodes a SCEP http request
func EncodeSCEPRequest(ctx context.Context, r *http.Request, request interface{}) error {
	req := request.(SCEPRequest)
	params := r.URL.Query()
	params.Set("operation", req.Operation)
	switch r.Method {
	case "GET":
		if len(req.Message) > 0 {
			return errors.New("only POSTPKIOperation supported")
		}
	case "POST":
		var buf bytes.Buffer
		_, err := buf.Write(req.Message)
		if err != nil {
			return err
		}
		r.Body = ioutil.NopCloser(&buf)
	default:
		return errors.New("method not supported")
	}
	r.URL.RawQuery = params.Encode()
	return nil
}

// DecodeSCEPRequest decodes an HTTP request to the SCEP server
// extracting the Operation and Message.
func decodeSCEPRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	msg, err := message(r)
	if err != nil {
		return nil, err
	}

	request := SCEPRequest{
		Message: msg,
	}

	return request, nil
}

// extract message from request
func message(r *http.Request) ([]byte, error) {
	switch r.Method {
	case "GET":
		var msg string
		q := r.URL.Query()
		if _, ok := q["message"]; ok {
			msg = q.Get("message")
		}
		return []byte(msg), nil
	case "POST":
		return ioutil.ReadAll(r.Body)
	default:
		return nil, errors.New("method not supported")
	}
}

// EncodeSCEPResponse writes a SCEP response back to the SCEP client.
func encodeSCEPResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(SCEPResponse)
	if resp.Err != nil {
		fmt.Println(resp.Err)
		return resp.Err
	}
	w.Header().Set("Content-Type", contentHeader(ctx, resp.CACertNum))
	w.Write(resp.Data)
	return nil
}

// DecodeSCEPResponse decodes a SCEP response
func DecodeSCEPResponse(ctx context.Context, r *http.Response) (interface{}, error) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	resp := SCEPResponse{
		Data: data,
	}
	header := r.Header.Get("Content-Type")
	if header == certChainHeader {
		// TODO decode the response instead of just passing []byte around
		// 0 or 1
		resp.CACertNum = 2
	}
	return resp, nil
}

const (
	certChainHeader = "application/x-x509-ca-ra-cert"
	leafHeader      = "application/x-x509-ca-cert"
	pkiOpHeader     = "application/x-pki-message"
)

func contentHeader(ctx context.Context, certNum int) string {
	op := ctx.Value("operation")
	switch op {
	case "GetCACert":
		if certNum > 1 {
			return certChainHeader
		}
		return leafHeader
	case "PKIOperation":
		return pkiOpHeader
	default:
		return "text/plain"
	}
}

func makeSCEPEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		op := ctx.Value("operation")
		if op == nil {
			return SCEPResponse{Err: errors.New("unknown operation")}, nil
		}
		req := request.(SCEPRequest)
		switch op {
		case "GetCACaps":
			caps, err := svc.GetCACaps(ctx)
			if err != nil {
				return SCEPResponse{Err: err}, nil
			}
			return SCEPResponse{Data: caps}, nil
		case "GetCACert":
			cert, certNum, err := svc.GetCACert(ctx)
			if err != nil {
				return SCEPResponse{Err: err, CACertNum: certNum}, nil
			}
			return SCEPResponse{Data: cert}, nil
		case "PKIOperation":
			resp, err := svc.PKIOperation(ctx, req.Message)
			if err != nil {
				return SCEPResponse{Err: err}, nil
			}
			return SCEPResponse{Data: resp}, nil
		default:
			return nil, errors.New("operation not implemented")
		}
	}
}

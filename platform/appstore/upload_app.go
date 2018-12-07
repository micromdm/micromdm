package appstore

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/pkg/errors"

	"github.com/micromdm/micromdm/pkg/httputil"
)

func (svc *AppService) UploadApp(ctx context.Context, manifestName string, manifest io.Reader, pkgName string, pkg io.Reader) error {
	if manifestName != "" {
		if err := svc.store.SaveFile(manifestName, manifest); err != nil {
			return err
		}
	}

	if pkgName != "" {
		if err := svc.store.SaveFile(pkgName, pkg); err != nil {
			return err
		}
	}

	return nil
}

type appUploadRequest struct {
	ManifestName string
	ManifestFile io.Reader

	PKGFilename string
	PKGFile     io.Reader
}

type appUploadResponse struct {
	Err error `json:"err,omitempty"`
}

func (r appUploadResponse) Failed() error { return r.Err }

func decodeAppUploadRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	appManifestFilename := r.FormValue("app_manifest_filename")
	manifestFile, _, err := r.FormFile("app_manifest_filedata")
	if err != nil && err != http.ErrMissingFile {
		return nil, errors.Wrap(err, "manifest file")
	}
	pkgFilename := r.FormValue("pkg_name")
	pkgFile, _, err := r.FormFile("pkg_filedata")
	if err != nil && err != http.ErrMissingFile {
		return nil, err
	}

	return appUploadRequest{
		ManifestName: appManifestFilename,
		ManifestFile: manifestFile,
		PKGFilename:  pkgFilename,
		PKGFile:      pkgFile,
	}, nil
}

func encodeUploadAppRequest(_ context.Context, r *http.Request, request interface{}) error {
	req := request.(appUploadRequest)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	if req.ManifestName != "" {
		partManifest, err := writer.CreateFormFile("app_manifest_filedata", req.ManifestName)
		if err != nil {
			return err
		}
		_, err = io.Copy(partManifest, req.ManifestFile)
		if err != nil {
			return errors.Wrap(err, "copying appmanifest file to multipart writer")
		}
		writer.WriteField("app_manifest_filename", req.ManifestName)
	}

	if req.PKGFilename != "" {
		partPkg, err := writer.CreateFormFile("pkg_filedata", req.PKGFilename)
		if err != nil {
			return err
		}
		_, err = io.Copy(partPkg, req.PKGFile)
		if err != nil {
			return errors.Wrap(err, "copying pkg file to multipart writer")
		}
		writer.WriteField("pkg_name", req.PKGFilename)
	}
	if err := writer.Close(); err != nil {
		return errors.Wrap(err, "closing multipart writer")
	}

	r.Header.Set("Content-Type", writer.FormDataContentType())
	r.Body = ioutil.NopCloser(body)
	return nil
}

func decodeUploadAppResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var resp appUploadResponse
	err := httputil.DecodeJSONResponse(r, &resp)
	return resp, err
}

func MakeUploadAppEndpiont(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(appUploadRequest)
		err = svc.UploadApp(ctx, req.ManifestName, req.ManifestFile, req.PKGFilename, req.PKGFile)
		return &appUploadResponse{
			Err: err,
		}, nil
	}
}

func (e Endpoints) UploadApp(ctx context.Context, manifestName string, manifest io.Reader, pkgName string, pkg io.Reader) error {
	request := appUploadRequest{
		ManifestName: manifestName,
		ManifestFile: manifest,
		PKGFilename:  pkgName,
		PKGFile:      pkg,
	}
	resp, err := e.AppUploadEndpoint(ctx, request)
	if err != nil {
		return err
	}
	return resp.(appUploadResponse).Err
}

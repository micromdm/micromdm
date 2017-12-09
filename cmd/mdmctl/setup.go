package main

import (
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/micromdm/micromdm/platform/api/server/apply"
	"github.com/micromdm/micromdm/platform/api/server/list"
	"github.com/micromdm/micromdm/platform/api/server/remove"
	"github.com/micromdm/micromdm/platform/profile"
)

type remoteServices struct {
	profilesvc profile.Service
	applysvc   apply.Service
	list       list.Service
	remove     remove.Service
}

func setupClient(logger log.Logger) (*remoteServices, error) {
	cfg, err := LoadServerConfig()
	if err != nil {
		return nil, err
	}
	applysvc, err := apply.NewClient(
		cfg.ServerURL, logger, cfg.APIToken,
		httptransport.SetClient(skipVerifyHTTPClient(cfg.SkipVerify)))
	if err != nil {
		return nil, err
	}

	profilesvc, err := profile.NewHTTPClient(
		cfg.ServerURL, cfg.APIToken, logger,
		httptransport.SetClient(skipVerifyHTTPClient(cfg.SkipVerify)))
	if err != nil {
		return nil, err
	}

	listsvc, err := list.NewClient(
		cfg.ServerURL, logger, cfg.APIToken,
		httptransport.SetClient(skipVerifyHTTPClient(cfg.SkipVerify)))
	if err != nil {
		return nil, err
	}

	rmsvc, err := remove.NewClient(
		cfg.ServerURL, logger, cfg.APIToken,
		httptransport.SetClient(skipVerifyHTTPClient(cfg.SkipVerify)))
	if err != nil {
		return nil, err
	}

	return &remoteServices{
		profilesvc: profilesvc,
		applysvc:   applysvc,
		list:       listsvc,
		remove:     rmsvc,
	}, nil
}

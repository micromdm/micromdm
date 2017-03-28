package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/user"
)

func NewClientConfig() (*ClientConfig, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}
	cfgData, err := ioutil.ReadFile(usr.HomeDir + "/.micromdm/default.json")
	if err != nil {
		return nil, fmt.Errorf("unable to load default config file: %s", err)
	}
	var cfg ClientConfig
	err = json.Unmarshal(cfgData, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal ~/.micromdm/default.json %s", err)
	}
	return &cfg, nil
}

type ClientConfig struct {
	APIToken  string `json:"api_token"`
	ServerURL string `json:"server_url"`
}

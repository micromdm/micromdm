package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/user"
	"path"
	"strings"
)

type configCommand struct {
	config *ClientConfig
}

func (cmd *configCommand) Run(args []string) error {
	if len(args) < 1 {
		cmd.Usage()
		os.Exit(1)
	}

	var config *ClientConfig
	if cfg, err := LoadClientConfig(); err == nil {
		config = cfg
	} else {
		config = new(ClientConfig)
	}
	var run func(*ClientConfig, []string) error
	switch strings.ToLower(args[0]) {
	case "set":
		run = setCmd
	default:
		cmd.Usage()
		os.Exit(1)
	}

	return run(config, args[1:])
}

func (cmd *configCommand) Usage() error {
	const help = `
mdmctl config set -h
`
	fmt.Println(help)
	return nil
}

func setCmd(cfg *ClientConfig, args []string) error {
	flagset := flag.NewFlagSet("set", flag.ExitOnError)
	var (
		flToken     = flagset.String("api-token", "", "api token to connect to micromdm server")
		flServerURL = flagset.String("server-url", "", "server url of micromdm server")
	)

	flagset.Usage = usageFor(flagset, "mdmctl config set [flags]")
	if err := flagset.Parse(args); err != nil {
		return err
	}

	if *flToken != "" {
		cfg.APIToken = *flToken
	}

	if *flServerURL != "" {
		if !strings.HasPrefix(*flServerURL, "http") ||
			!strings.HasPrefix(*flServerURL, "https") {
			*flServerURL = "https://" + *flServerURL
		}
		u, err := url.Parse(*flServerURL)
		if err != nil {
			return err
		}
		u.Scheme = "https"
		u.Path = "/"
		cfg.ServerURL = u.String()
	}

	return SaveClientConfig(cfg)
}

func clientConfigDirPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return path.Join(usr.HomeDir, ".micromdm"), nil
}

func clientConfigPath() (string, error) {
	dir, err := clientConfigDirPath()
	return path.Join(dir, "default.json"), err
}

func SaveClientConfig(cfg *ClientConfig) error {
	path, err := clientConfigDirPath()
	if err != nil {
		return err
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0777)
		if err != nil {
			return err
		}
	}
	path, _ = clientConfigPath()
	if err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	if cfg == nil {
		cfg = new(ClientConfig)
	}
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(cfg)
}

func LoadClientConfig() (*ClientConfig, error) {
	path, err := clientConfigPath()
	if err != nil {
		return nil, err
	}
	cfgData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("unable to load default config file: %s", err)
	}
	var cfg ClientConfig
	err = json.Unmarshal(cfgData, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal %s : %s", path, err)
	}
	return &cfg, nil
}

type ClientConfig struct {
	APIToken  string `json:"api_token"`
	ServerURL string `json:"server_url"`
}

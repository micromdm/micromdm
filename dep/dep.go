package dep

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/garyburd/go-oauth/oauth"
)

const (
	libraryVersion = "0.0.2"
	defaultBaseURL = "https://mdmenrollment.apple.com"
	userAgent      = "micromdm/" + libraryVersion
	mediaType      = "application/json;charset=UTF8"
)

// Client interacts with DEP
type Client interface {
	AccountService
	DeviceService
	ProfileService
}

// Config is a configuration struct for DEP
type Config struct {
	ConsumerKey    string //given by apple
	ConsumerSecret string //given by apple
	AccessToken    string //given by apple
	AccessSecret   string //given by apple

	AuthSessionToken string //requested from DEP using above credentials
	sessionExpires   time.Time
	url              *url.URL
	debug            bool
}

func (c *Config) session() error {
	if c.AuthSessionToken == "" {
		err := c.newSession()
		if err != nil {
			return err
		}
	}

	if time.Now().After(c.sessionExpires) {
		err := c.newSession()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Config) newSession() error {
	var authSessionToken struct {
		AuthSessionToken string `json:"auth_session_token"`
	}
	consumerCredentials := oauth.Credentials{
		Token:  c.ConsumerKey,
		Secret: c.ConsumerSecret,
	}

	accessCredentials := &oauth.Credentials{
		Token:  c.AccessToken,
		Secret: c.AccessSecret,
	}
	form := url.Values{}

	// session url, relative to basepath
	rel, err := url.Parse("/session")
	if err != nil {
		return err
	}
	sessionURL := c.url.ResolveReference(rel)

	oauthClient := oauth.Client{
		SignatureMethod: oauth.HMACSHA1,
		TokenRequestURI: sessionURL.String(),
		Credentials:     consumerCredentials,
	}

	// create request
	req, err := http.NewRequest("GET", oauthClient.TokenRequestURI, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}

	// set Authorization Header
	if err := oauthClient.SetAuthorizationHeader(req.Header, accessCredentials, "GET", req.URL, form); err != nil {
		return err
	}
	// add headers
	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("Content-Type", mediaType)
	req.Header.Add("Accept", mediaType)
	req.Header.Add("X-Server-Protocol-Version", "2")

	// get Authorization Header
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// check resp statuscode
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error establishing DEP session: %v", resp.Status)
	}

	// decode token from response
	if err = decodeJSON(c.debug, resp.Body, &authSessionToken); err != nil {
		return err
	}

	// set token and expiration value
	c.AuthSessionToken = authSessionToken.AuthSessionToken
	c.sessionExpires = time.Now().Add(time.Minute * 3)
	return nil
}

type depClient struct {
	client *http.Client

	// Base URL for API requests.
	BaseURL *url.URL

	// User agent for client
	UserAgent string

	Config *Config

	accountService
	deviceService
	profileService
}

// NewClient creates a new HTTP client for communicating with DEP
func NewClient(config *Config, options ...func(*Config) error) (Client, error) {
	for _, option := range options {
		if err := option(config); err != nil {
			return nil, err
		}
	}

	if config.url == nil {
		config.url, _ = url.Parse(defaultBaseURL)
	}
	c := &depClient{client: http.DefaultClient, BaseURL: config.url, UserAgent: userAgent, Config: config}
	c.accountService = accountService{client: c}
	c.deviceService = deviceService{client: c}
	c.profileService = profileService{client: c}
	return c, nil
}

// ServerURL allows the user to provide a URL for DEP server other than the default
// useful for testing with depsim
func ServerURL(baseURL string) func(*Config) error {
	return func(c *Config) error {
		var err error
		c.url, err = url.Parse(baseURL)
		return err
	}
}

// Debug will preint responses from DEP to stdout.
func Debug() func(*Config) error {
	return func(c *Config) error {
		c.debug = true
		return nil
	}
}

// NewRequest creates a DEP request
func (c *depClient) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)
	buf := new(bytes.Buffer)
	if body != nil {
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("Content-Type", mediaType)
	req.Header.Add("Accept", mediaType)
	req.Header.Add("X-Server-Protocol-Version", "2")
	return req, nil
}

// Do sends an API request and returns the API response.
func (c *depClient) Do(req *http.Request, into interface{}) error {
	// set/check session token
	err := c.Config.session()
	if err != nil {
		return err
	}
	req.Header.Add("X-ADM-Auth-Session", c.Config.AuthSessionToken)

	// perform request
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("DEP API Error: %v", string(body))
	}

	return decodeJSON(c.Config.debug, resp.Body, into)
}

func decodeJSON(debug bool, body io.Reader, into interface{}) error {
	var dec *json.Decoder
	if debug {
		dec = json.NewDecoder(io.TeeReader(body, os.Stdout))
	} else {
		dec = json.NewDecoder(body)
	}

	return dec.Decode(into)
}

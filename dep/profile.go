package dep

import "github.com/pkg/errors"

const (
	defineProfilePath = "profile"
	assignProfilePath = "profile/devices"
)

type Profile struct {
	ProfileName           string   `json:"profile_name"`
	ProfileUUID           string   `json:"profile_uuid,omitempty"`
	URL                   string   `json:"url"`
	AllowPairing          bool     `json:"allow_pairing,omitempty"`
	IsSupervised          bool     `json:"is_supervised,omitempty"`
	IsMultiUser           bool     `json:"is_multi_user,omitempty"`
	IsMandatory           bool     `json:"is_mandatory,omitempty"`
	AwaitDeviceConfigured bool     `json:"await_device_configured,omitempty"`
	IsMDMRemovable        bool     `json:"is_mdm_removable"`
	SupportPhoneNumber    string   `json:"support_phone_number,omitempty"`
	AutoAdvanceSetup      bool     `json:"auto_advance_setup,omitempty"`
	SupportEmailAddress   string   `json:"support_email_address,omitempty"`
	OrgMagic              string   `json:"org_magic"`
	AnchorCerts           []string `json:"anchor_certs,omitempty"`
	SupervisingHostCerts  []string `json:"supervising_host_certs,omitempty"`
	SkipSetupItems        []string `json:"skip_setup_items,omitempty"`
	Department            string   `json:"department,omitempty"`
	Devices               []string `json:"devices"`
	Language              string   `json:"language,omitempty"`
	Region                string   `json:"region,omitempty"`
	ConfigurationWebURL   string   `json:"configuration_web_url,omitempty"`
}

type ProfileResponse struct {
	ProfileUUID string            `json:"profile_uuid"`
	Devices     map[string]string `json:"devices"`
}

func (c *Client) DefineProfile(request *Profile) (*ProfileResponse, error) {
	var response ProfileResponse
	req, err := c.newRequest("POST", defineProfilePath, request)
	if err != nil {
		return nil, errors.Wrap(err, "create define profile request")
	}
	err = c.do(req, &response)
	return &response, errors.Wrap(err, "define profile")
}

func (c *Client) AssignProfile(uuid string, serials ...string) (*ProfileResponse, error) {
	var response ProfileResponse
	var request = struct {
		ProfileUUID string   `json:"profile_uuid"`
		Devices     []string `json:"devices"`
	}{
		ProfileUUID: uuid,
		Devices:     serials,
	}

	req, err := c.newRequest("PUT", assignProfilePath, request)
	if err != nil {
		return nil, errors.Wrap(err, "create assign profile request")
	}
	err = c.do(req, &response)
	return &response, errors.Wrap(err, "assign profile")
}

func (c *Client) FetchProfile(uuid string) (*Profile, error) {
	var response Profile
	req, err := c.newRequest("GET", defineProfilePath, nil)
	if err != nil {
		return nil, errors.Wrap(err, "create fetch profile request")
	}
	query := req.URL.Query()
	query.Add("profile_uuid", uuid)
	req.URL.RawQuery = query.Encode()

	err = c.do(req, &response)
	return &response, errors.Wrap(err, "fetch profile")
}

func (c *Client) RemoveProfile(serials ...string) (map[string]string, error) {
	var response struct {
		Devices map[string]string `json:"devices"`
	}
	var request = struct {
		Devices []string `json:"devices"`
	}{Devices: serials}

	req, err := c.newRequest("DELETE", assignProfilePath, &request)
	if err != nil {
		return nil, errors.Wrap(err, "create fetch profile request")
	}
	err = c.do(req, &response)
	return response.Devices, errors.Wrap(err, "remove profile")
}

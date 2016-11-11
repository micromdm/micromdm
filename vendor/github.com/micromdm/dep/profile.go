package dep

const (
	defineProfilePath = "profile"
	assignProfilePath = "profile/devices"
)

// ProfileService allows Defining, Assigning and Fetching profiles from DEP
type ProfileService interface {
	DefineProfile(profile *Profile) (*ProfileResponse, error)
	AssignProfile(profileUUID string, devices []string) (*ProfileResponse, error)
	FetchProfile(profileUUID string) (*Profile, error)
}

type profileService struct {
	client *depClient
}

// Profile is a DEP setup profile.
// The profile can be defined, assigned and fetched.
type Profile struct {
	ProfileName           string   `json:"profile_name"`
	URL                   string   `json:"url"`
	AllowPairing          bool     `json:"allow_pairing,omitempty"`
	IsSupervised          bool     `json:"is_supervised,omitempty"`
	IsMultiUser           bool     `json:"is_multi_user,omitempty"`
	IsMandatory           bool     `json:"is_mandatory,omitempty"`
	AwaitDeviceConfigured bool     `json:"await_device_configured,omitempty"`
	IsMDMRemovable        bool     `json:"is_mdm_removable"`
	SupportPhoneNumber    string   `json:"support_phone_number,omitempty"`
	SupportEmailAddress   string   `json:"support_email_address,omitempty"`
	OrgMagic              string   `json:"org_magic"`
	AnchorCerts           []string `json:"anchor_certs,omitempty"`
	SupervisingHostCerts  []string `json:"supervising_host_certs,omitempty"`
	SkipSetupItems        []string `json:"skip_setup_items,omitempty"`
	Department            string   `json:"deparment,omitempty"`
	Devices               []string `json:"devices"`
}

// ProfileResponse is the response body for Define Profile
type ProfileResponse struct {
	ProfileUUID string            `json:"profile_uuid"`
	Devices     map[string]string `json:"devices"`
}

// DefineProfile is a Define Profile request to DEP
func (s profileService) DefineProfile(request *Profile) (*ProfileResponse, error) {
	var response ProfileResponse
	req, err := s.client.NewRequest("POST", defineProfilePath, request)
	if err != nil {
		return nil, err
	}
	err = s.client.Do(req, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// AssignProfile is an assign profile request to DEP
func (s profileService) AssignProfile(profileUUID string, devices []string) (*ProfileResponse, error) {
	var response ProfileResponse
	var request = struct {
		ProfileUUID string   `json:"profile_uuid"`
		Devices     []string `json:"devices"`
	}{profileUUID, devices}
	req, err := s.client.NewRequest("PUT", assignProfilePath, request)
	if err != nil {
		return nil, err
	}
	err = s.client.Do(req, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// FetchProfile is a Fetch Profile request to DEP
func (s profileService) FetchProfile(profileUUID string) (*Profile, error) {
	var response Profile
	req, err := s.client.NewRequest("GET", defineProfilePath, nil)
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	query.Add("profile_uuid", profileUUID)
	req.URL.RawQuery = query.Encode()
	err = s.client.Do(req, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

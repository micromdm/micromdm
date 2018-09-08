package dep

import (
	"time"

	"github.com/pkg/errors"
)

const (
	fetchDevicesPath  = "server/devices"
	syncDevicesPath   = "devices/sync"
	deviceDetailsPath = "devices"
)

type Device struct {
	SerialNumber       string    `json:"serial_number"`
	Model              string    `json:"model"`
	Description        string    `json:"description"`
	Color              string    `json:"color"`
	AssetTag           string    `json:"asset_tag,omitempty"`
	ProfileStatus      string    `json:"profile_status"`
	ProfileUUID        string    `json:"profile_uuid,omitempty"`
	ProfileAssignTime  time.Time `json:"profile_assign_time,omitempty"`
	ProfilePushTime    time.Time `json:"profile_push_time,omitempty"`
	DeviceAssignedDate time.Time `json:"device_assigned_date,omitempty"`
	DeviceAssignedBy   string    `json:"device_assigned_by,omitempty"`
	OS                 string    `json:"os,omitempty"`
	DeviceFamily       string    `json:"device_family,omitempty"`
	// sync fields
	OpType string    `json:"op_type,omitempty"`
	OpDate time.Time `json:"op_date,omitempty"`
	// details fields
	ResponseStatus string `json:"response_status,omitempty"`
}

// DeviceRequestOption is an optional parameter for the DeviceService API.
// The option can be used to set Cursor or Limit options for the request.
type DeviceRequestOption func(*deviceRequestOpts) error

type deviceRequestOpts struct {
	Cursor string `json:"cursor,omitempty"`
	Limit  int    `json:"limit,omitempty"`
}

// Cursor is an optional argument that can be added to FetchDevices
func Cursor(cursor string) DeviceRequestOption {
	return func(opts *deviceRequestOpts) error {
		opts.Cursor = cursor
		return nil
	}
}

// Limit is an optional argument that can be passed to FetchDevices and SyncDevices
func Limit(limit int) DeviceRequestOption {
	return func(opts *deviceRequestOpts) error {
		if limit > 1000 {
			return errors.New("limit must not be higher than 1000")
		}
		opts.Limit = limit
		return nil
	}
}

type DeviceResponse struct {
	Devices      []Device  `json:"devices"`
	Cursor       string    `json:"cursor"`
	FetchedUntil time.Time `json:"fetched_until"`
	MoreToFollow bool      `json:"more_to_follow"`
}

func (c *Client) FetchDevices(opts ...DeviceRequestOption) (*DeviceResponse, error) {
	request := &deviceRequestOpts{}
	for _, option := range opts {
		if err := option(request); err != nil {
			return nil, err
		}
	}
	var response DeviceResponse
	req, err := c.newRequest("POST", fetchDevicesPath, request)
	if err != nil {
		return nil, errors.Wrap(err, "create fetch devices request")
	}
	err = c.do(req, &response)
	return &response, errors.Wrap(err, "fetch devices")
}

func (c *Client) SyncDevices(cursor string, opts ...DeviceRequestOption) (*DeviceResponse, error) {
	request := &deviceRequestOpts{Cursor: cursor}
	for _, option := range opts {
		if err := option(request); err != nil {
			return nil, err
		}
	}
	var response DeviceResponse
	req, err := c.newRequest("POST", syncDevicesPath, request)
	if err != nil {
		return nil, errors.Wrap(err, "create sync devices request")
	}
	err = c.do(req, &response)
	return &response, errors.Wrap(err, "sync devices")
}

type DeviceDetailsResponse struct {
	Devices map[string]Device `json:"devices"`
}

func (c *Client) DeviceDetails(serials ...string) (*DeviceDetailsResponse, error) {
	request := struct {
		Devices []string `json:"devices"`
	}{
		Devices: serials,
	}

	var response DeviceDetailsResponse
	req, err := c.newRequest("POST", deviceDetailsPath, request)
	if err != nil {
		return nil, errors.Wrap(err, "create device details request")
	}
	err = c.do(req, &response)
	return &response, errors.Wrap(err, "get device details")
}

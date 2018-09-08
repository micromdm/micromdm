package dep

import (
	"fmt"
	"time"
)

const (
	fetchDevicesPath  = "server/devices"
	syncDevicesPath   = "devices/sync"
	deviceDetailsPath = "devices"
)

// DeviceService allows fetching and syncing devices, as well as requesting device details
/*
Use (Cursor() and Limit() as optional arguments for Fetch Devices, example:
	fetchResponse, err := client.FetchDevices(dep.Limit(100))
	if err != nil {
		// handle err
	}
	fmt.Println(fetchResponse.Devices)
*/
type DeviceService interface {
	FetchDevices(opts ...DeviceRequestOption) (*DeviceResponse, error)
	SyncDevices(cursor string, opts ...DeviceRequestOption) (*DeviceResponse, error)
	DeviceDetails(devices []string) (*DeviceDetailsResponse, error)
}

type deviceService struct {
	client *depClient
}

// Device is a DEP device
type Device struct {
	SerialNumber       string    `json:"serial_number"`
	Model              string    `json:"model"`
	Description        string    `json:"description"`
	Color              string    `json:"color"`
	AssetTag           string    `json:"asset_tag"`
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
			return fmt.Errorf("Limit must not be higher than 1000")
		}
		opts.Limit = limit
		return nil
	}
}

// DeviceResponse is a DEP FetchDevices response
type DeviceResponse struct {
	Devices      []Device  `json:"devices"`
	Cursor       string    `json:"cursor"`
	FetchedUntil time.Time `json:"fetched_until"`
	MoreToFollow bool      `json:"more_to_follow"`
}

// FetchDevices returns the result of a Fetch Devices request from DEP
func (s deviceService) FetchDevices(opts ...DeviceRequestOption) (*DeviceResponse, error) {
	request := &deviceRequestOpts{}
	for _, option := range opts {
		if err := option(request); err != nil {
			return nil, err
		}
	}
	var response DeviceResponse
	req, err := s.client.NewRequest("POST", fetchDevicesPath, request)
	if err != nil {
		return nil, err
	}
	err = s.client.Do(req, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// SyncDevices returns the result of a Sync Devices request from DEP
func (s deviceService) SyncDevices(cursor string, opts ...DeviceRequestOption) (*DeviceResponse, error) {
	request := &deviceRequestOpts{Cursor: cursor}
	for _, option := range opts {
		if err := option(request); err != nil {
			return nil, err
		}
	}
	var response DeviceResponse
	req, err := s.client.NewRequest("POST", syncDevicesPath, request)
	if err != nil {
		return nil, err
	}
	err = s.client.Do(req, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// DeviceDetailsResponse is a response for a DeviceDetails request
type DeviceDetailsResponse struct {
	Devices map[string]Device `json:"devices"`
}

// DeviceDetails returns the result of a Sync Devices request from DEP
func (s deviceService) DeviceDetails(devices []string) (*DeviceDetailsResponse, error) {
	request := struct {
		Devices []string `json:"devices"`
	}{
		Devices: devices,
	}
	var response DeviceDetailsResponse
	req, err := s.client.NewRequest("POST", deviceDetailsPath, request)
	if err != nil {
		return nil, err
	}
	err = s.client.Do(req, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

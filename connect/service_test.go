package connect

import (
	"github.com/micromdm/mdm"
	"github.com/micromdm/micromdm/applications"
	"github.com/micromdm/micromdm/device"
	"testing"
	"time"
)

type MockDevices struct{}

func (md MockDevices) New(src string, d *device.Device) (string, error) {
	return "", nil
}
func (md MockDevices) GetDeviceByUDID(udid string, fields ...string) (*device.Device, error) {
	return &device.Device{}, nil
}
func (md MockDevices) GetDeviceByUUID(uuid string, fields ...string) (*device.Device, error) {
	return &device.Device{}, nil
}
func (md MockDevices) Devices(params ...interface{}) ([]device.Device, error) {
	return []device.Device{
		{
			UUID: "ABCD-EFGH-IJKL",
		},
	}, nil
}
func (md MockDevices) Save(msg string, dev *device.Device) error {
	return nil
}

type MockApps struct{}

func (ma MockApps) New(a *applications.Application) (string, error) {
	return "", nil
}
func (ma MockApps) Applications(params ...interface{}) ([]applications.Application, error) {
	return []applications.Application{}, nil
}
func (ma MockApps) GetApplicationsByDeviceUUID(deviceUUID string) ([]applications.Application, error) {
	return []applications.Application{}, nil
}
func (ma MockApps) SaveApplicationByDeviceUUID(deviceUUID string, app *applications.Application) error {
	return nil
}

type MockCmd struct{}

func (mc MockCmd) NewCommand(*mdm.CommandRequest) (*mdm.Payload, error) {
	return &mdm.Payload{}, nil
}
func (mc MockCmd) NextCommand(udid string) ([]byte, int, error) {
	return []byte{}, 0, nil
}
func (mc MockCmd) DeleteCommand(deviceUDID, commandUUID string) (int, error) {
	return 0, nil
}

type MockContext struct{}

func (mc MockContext) Done() <-chan struct{} {
	ch := make(chan struct{})

	return ch
}
func (mc MockContext) Err() error {
	return nil
}
func (mc MockContext) Deadline() (deadline time.Time, ok bool) {
	return time.Now(), true
}
func (mc MockContext) Value(key interface{}) interface{} {
	return nil
}

func TestAckQueryResponses(t *testing.T) {
	response := mdm.Response{
		UDID:           "00000000-1111-2222-3333-444455556666",
		Status:         "Acknowledged",
		CommandUUID:    "10000000-1111-2222-3333-444455556666",
		RequestType:    "DeviceInformation",
		QueryResponses: mdm.QueryResponses{},
	}

	mockDevices := MockDevices{}
	mockApps := MockApps{}
	mockCmd := MockCmd{}

	svc := NewService(mockDevices, mockApps, mockCmd)
	svc.Acknowledge(MockContext{}, response)
}

func TestAckInstalledApplicationList(t *testing.T) {

	response := mdm.Response{
		UDID:        "00000000-1111-2222-3333-444455556666",
		Status:      "Acknowledged",
		CommandUUID: "10000000-1111-2222-3333-444455556666",
		RequestType: "InstalledApplicationList",
		InstalledApplicationList: []mdm.InstalledApplicationListItem{
			{
				Name:       "Wireless Network Utility",
				BundleSize: 2416111,
			},
			{
				Name:         "Keychain Access",
				Identifier:   "com.apple.keychainaccess",
				ShortVersion: "9.0",
				Version:      "9.0",
				BundleSize:   14166172,
			},
		},
	}

	mockDevices := MockDevices{}
	mockApps := MockApps{}
	mockCmd := MockCmd{}

	svc := NewService(mockDevices, mockApps, mockCmd)
	svc.Acknowledge(MockContext{}, response)
}

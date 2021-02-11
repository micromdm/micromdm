package dep

import "github.com/pkg/errors"

const (
	activationLockPath = "device/activationlock"
)

type ActivationLockRequest struct {
	Device string `json:"device"`

	//  If the escrow key is not provided, the device will be locked to the person who created the MDM server in the portal.
	// https://developer.apple.com/documentation/devicemanagement/device_assignment/activation_lock_a_device/creating_and_using_bypass_codes
	// The EscrowKey is a hex-encoded PBKDF2 derivation of the bypass code. See activationlock.BypassCode.
	EscrowKey string `json:"escrow_key"`

	LostMessage string `json:"lost_message"`
}

type ActivationLockResponse struct {
	SerialNumber string `json:"serial_number"`
	Status       string `json:"response_status"`
}

func (c *Client) ActivationLock(alr ActivationLockRequest) (*ActivationLockResponse, error) {
	req, err := c.newRequest("POST", activationLockPath, &alr)
	if err != nil {
		return nil, errors.Wrap(err, "create activation lock request")
	}

	var response ActivationLockResponse
	err = c.do(req, &response)
	return &response, errors.Wrap(err, "activation lock")
}

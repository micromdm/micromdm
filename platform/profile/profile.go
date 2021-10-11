package profile

import (
	"github.com/micromdm/micromdm/platform/profile/internal/profileproto"

	"github.com/gogo/protobuf/proto"
	"github.com/groob/plist"
	"github.com/pkg/errors"
	"go.mozilla.org/pkcs7"
)

type Mobileconfig []byte

// only used to parse plists to get the PayloadIdentifier
type payloadIdentifier struct {
	PayloadIdentifier string
}

func (mc *Mobileconfig) GetPayloadIdentifier() (string, error) {
	mcBytes := *mc
	if len(mcBytes) > 5 && string(mcBytes[0:5]) != "<?xml" {
		p7, err := pkcs7.Parse(mcBytes)
		if err != nil {
			return "", errors.Wrapf(err, "Mobileconfig is not XML nor PKCS7 parseable")
		}
		err = p7.Verify()
		if err != nil {
			return "", err
		}
		mcBytes = Mobileconfig(p7.Content)
	}
	var pId payloadIdentifier
	err := plist.Unmarshal(mcBytes, &pId)
	if err != nil {
		return "", err
	}
	if pId.PayloadIdentifier == "" {
		return "", errors.New("empty PayloadIdentifier in profile")
	}
	return pId.PayloadIdentifier, err
}

type Profile struct {
	Identifier   string
	Mobileconfig Mobileconfig
}

// Validate checks the internal consistency and validity of a Profile structure
func (p *Profile) Validate() error {
	if p.Identifier == "" {
		return errors.New("Profile struct must have Identifier")
	}
	if len(p.Mobileconfig) < 1 {
		return errors.New("no Mobileconfig data")
	}
	payloadId, err := p.Mobileconfig.GetPayloadIdentifier()
	if err != nil {
		return err
	}
	if payloadId != p.Identifier {
		return errors.New("payload Identifier does not match Profile")
	}
	return nil
}

func MarshalProfile(p *Profile) ([]byte, error) {
	protobp := profileproto.Profile{
		Id:           p.Identifier,
		Mobileconfig: p.Mobileconfig,
	}
	return proto.Marshal(&protobp)
}

func UnmarshalProfile(data []byte, p *Profile) error {
	var pb profileproto.Profile
	if err := proto.Unmarshal(data, &pb); err != nil {
		return err
	}
	p.Identifier = pb.GetId()
	p.Mobileconfig = pb.GetMobileconfig()
	return nil
}

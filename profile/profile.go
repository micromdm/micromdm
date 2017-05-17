package profile

import (
	"errors"

	"github.com/gogo/protobuf/proto"
	"github.com/groob/plist"
	"github.com/micromdm/micromdm/profile/internal/profileproto"
)

type Mobileconfig []byte

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
	payloadId, err := GetMobileconfigIdentifier(p.Mobileconfig)
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

// only used to parse plists to get the PayloadIdentifier
type payloadIdentifier struct {
	PayloadIdentifier string
}

func GetMobileconfigIdentifier(mc Mobileconfig) (string, error) {
	// TODO: support CMS signed profiles
	var pId payloadIdentifier
	err := plist.Unmarshal(mc, &pId)
	if err != nil {
		return "", err
	}
	if pId.PayloadIdentifier == "" {
		return "", errors.New("empty PayloadIdentifier in profile")
	}
	return pId.PayloadIdentifier, err
}

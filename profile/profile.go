package profile

import (
	"github.com/gogo/protobuf/proto"
	"github.com/micromdm/micromdm/profile/internal/profileproto"
)

type Mobileconfig []byte

type Profile struct {
	Identifier string
	Mobileconfig
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

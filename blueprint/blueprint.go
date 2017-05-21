package blueprint

import (
	"github.com/gogo/protobuf/proto"
	"github.com/micromdm/micromdm/blueprint/internal/blueprintproto"
)

type Mobileconfig []byte

type Blueprint struct {
	UUID            string         `json:"uuid"`
	Name            string         `json:"name"`
	When            string         `json:"when"`
	ApplicationURLs []string       `json:"install_application_manifest_urls"`
	Profiles        []Mobileconfig `json:"profiles"`
}

func MarshalBlueprint(bp *Blueprint) ([]byte, error) {
	var profiles [][]byte
	for _, p := range bp.Profiles {
		profiles = append(profiles, []byte(p))
	}
	protobp := blueprintproto.Blueprint{
		Uuid:          bp.UUID,
		Name:          bp.Name,
		When:          bp.When,
		ManifestUrls:  bp.ApplicationURLs,
		Mobileconfigs: profiles,
	}
	return proto.Marshal(&protobp)
}

func UnmarshalBlueprint(data []byte, bp *Blueprint) error {
	var pb blueprintproto.Blueprint
	if err := proto.Unmarshal(data, &pb); err != nil {
		return err
	}
	var profiles []Mobileconfig
	for _, p := range pb.GetMobileconfigs() {
		profiles = append(profiles, Mobileconfig(p))
	}
	bp.UUID = pb.GetUuid()
	bp.Name = pb.GetName()
	bp.When = pb.GetWhen()
	bp.ApplicationURLs = pb.GetManifestUrls()
	bp.Profiles = profiles
	return nil
}

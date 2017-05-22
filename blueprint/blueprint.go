package blueprint

import (
	"errors"

	"github.com/gogo/protobuf/proto"
	"github.com/micromdm/micromdm/blueprint/internal/blueprintproto"
)

// ApplyAt is a case-insensitive string that specifies at which point the
// system should apply a Blueprint to devices. For example if a Blueprint has
// an ApplyAt of "Enroll" then that profile will be applied immediately after
// a device's enrollment in the MDM system. Currently "Enroll" is the only
// supported value but more are planned.
const (
	ApplyAtEnroll string = "Enroll"
)

type Blueprint struct {
	UUID               string   `json:"uuid"`
	Name               string   `json:"name"`
	ApplicationURLs    []string `json:"install_application_manifest_urls"`
	ProfileIdentifiers []string `json:"profile_ids"`
	ApplyAt            []string `json:"apply_at"`
}

func (bp *Blueprint) Verify() error {
	if bp.Name == "" || bp.UUID == "" {
		return errors.New("Blueprint must have Name and UUID")
	}
	return nil
}

func MarshalBlueprint(bp *Blueprint) ([]byte, error) {
	protobp := blueprintproto.Blueprint{
		Uuid:         bp.UUID,
		Name:         bp.Name,
		ManifestUrls: bp.ApplicationURLs,
		ProfileIds:   bp.ProfileIdentifiers,
		ApplyAt:      bp.ApplyAt,
	}
	return proto.Marshal(&protobp)
}

func UnmarshalBlueprint(data []byte, bp *Blueprint) error {
	var pb blueprintproto.Blueprint
	if err := proto.Unmarshal(data, &pb); err != nil {
		return err
	}
	bp.UUID = pb.GetUuid()
	bp.Name = pb.GetName()
	bp.ApplicationURLs = pb.GetManifestUrls()
	bp.ProfileIdentifiers = pb.GetProfileIds()
	bp.ApplyAt = pb.GetApplyAt()
	return nil
}

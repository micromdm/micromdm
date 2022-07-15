package apns

import (
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"

	"github.com/micromdm/micromdm/platform/apns/internal/pushproto"
)

type PushInfo struct {
	UDID      string `db:"udid"`
	PushMagic string `db:"push_magic"`
	Token     string `db:"token"`
	MDMTopic  string `db:"mdm_topic"`
}

func MarshalPushInfo(p *PushInfo) ([]byte, error) {
	protopush := pushproto.PushInfo{
		Udid:      p.UDID,
		PushMagic: p.PushMagic,
		Token:     p.Token,
		MdmTopic:  p.MDMTopic,
	}
	return proto.Marshal(&protopush)
}

func UnmarshalPushInfo(data []byte, p *PushInfo) error {
	var pb pushproto.PushInfo
	if err := proto.Unmarshal(data, &pb); err != nil {
		return errors.Wrap(err, "unmarshal proto to PushInfo")
	}
	p.UDID = pb.GetUdid()
	p.Token = pb.GetToken()
	p.PushMagic = pb.GetPushMagic()
	p.MDMTopic = pb.GetMdmTopic()
	return nil
}

package config

import (
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"

	"github.com/micromdm/micromdm/platform/config/internal/configproto"
)

// ServerConfig holds the configuration of the MDM Server.
type ServerConfig struct {
	PushCertificate []byte
	PrivateKey      []byte
}

func MarshalServerConfig(conf *ServerConfig) ([]byte, error) {
	pb := configproto.ServerConfig{
		PushCertificate:    conf.PushCertificate,
		PushCertificateKey: conf.PrivateKey,
	}
	data, err := proto.Marshal(&pb)
	return data, errors.Wrap(err, "marshal server config to proto")
}

func UnmarshalServerConfig(data []byte, conf *ServerConfig) error {
	var pb configproto.ServerConfig
	if err := proto.Unmarshal(data, &pb); err != nil {
		return errors.Wrap(err, "unmarshal server config from proto")
	}
	conf.PushCertificate = pb.GetPushCertificate()
	conf.PrivateKey = pb.GetPushCertificateKey()
	return nil
}

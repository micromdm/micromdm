// Package user provides utilites for managing users with MDM.
package user

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"

	"github.com/micromdm/micromdm/platform/user/internal/userproto"
)

type User struct {
	UUID          string `json:"uuid"`
	UDID          string `json:"udid"`
	UserID        string `json:"user_id"`
	UserShortname string `json:"user_shortname"`
	UserLongname  string `json:"user_longname"`
	AuthToken     string `json:"auth_token"`
	PasswordHash  []byte `json:"password_hash"`
	Hidden        bool   `json:"hidden"`
}

func NewFromRequest(u User) (*User, error) {
	newUser := User{
		UUID:          uuid.New().String(),
		UserShortname: u.UserShortname,
		UserLongname:  u.UserLongname,
		PasswordHash:  u.PasswordHash,
		Hidden:        u.Hidden,
	}
	return &newUser, nil
}

func MarshalUser(u *User) ([]byte, error) {
	pb := userproto.User{
		Uuid:          u.UUID,
		Udid:          u.UDID,
		UserId:        u.UserID,
		UserShortname: u.UserShortname,
		UserLongname:  u.UserLongname,
		AuthToken:     u.AuthToken,
		PasswordHash:  u.PasswordHash,
		Hidden:        u.Hidden,
	}
	return proto.Marshal(&pb)
}

func UnmarshalUser(data []byte, u *User) error {
	var pb userproto.User
	if err := proto.Unmarshal(data, &pb); err != nil {
		return errors.Wrap(err, "unmarshal proto to user")
	}
	u.UUID = pb.GetUuid()
	u.UDID = pb.GetUdid()
	u.UserID = pb.GetUserId()
	u.UserShortname = pb.GetUserShortname()
	u.UserLongname = pb.GetUserLongname()
	u.AuthToken = pb.GetAuthToken()
	u.PasswordHash = pb.GetPasswordHash()
	u.Hidden = pb.GetHidden()
	return nil
}

package mdm

import (
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/micromdm/micromdm/mdm/internal/connectproto"
)

type AcknowledgeEvent struct {
	ID       string
	Time     time.Time
	Response Response
	Params   map[string]string
	Raw      []byte
}

type Response struct {
	RequestType  string `json:"request_type,omitempty" plist:",omitempty"`
	UDID         string
	UserID       *string `json:"user_id,omitempty" plist:"UserID,omitempty"`
	EnrollmentID *string `json:"enrollment_id,omitempty" plist:"EnrollmentID,omitempty"`
	Status       string
	CommandUUID  string
	ErrorChain   []ErrorChainItem `json:"error_chain" plist:",omitempty"`
}

type ErrorChainItem struct {
	ErrorCode            int    `json:"error_code,omitempty"`
	ErrorDomain          string `json:"error_domain,omitempty"`
	LocalizedDescription string `json:"localized_description,omitempty"`
	USEnglishDescription string `json:"us_english_description,omitempty"`
}

func MarshalAcknowledgeEvent(e *AcknowledgeEvent) ([]byte, error) {
	response := &connectproto.Response{
		CommandUuid: e.Response.CommandUUID,
		Udid:        e.Response.UDID,
		Status:      e.Response.Status,
		RequestType: e.Response.RequestType,
	}
	if e.Response.UserID != nil {
		response.UserId = *e.Response.UserID
	}
	if e.Response.EnrollmentID != nil {
		response.EnrollmentId = *e.Response.EnrollmentID
	}

	return proto.Marshal(&connectproto.Event{
		Id:       e.ID,
		Time:     e.Time.UnixNano(),
		Response: response,
		Params:   e.Params,
		Raw:      e.Raw,
	})
}

func UnmarshalAcknowledgeEvent(data []byte, e *AcknowledgeEvent) error {
	var pb connectproto.Event
	if err := proto.Unmarshal(data, &pb); err != nil {
		return err
	}
	e.ID = pb.Id
	e.Time = time.Unix(0, pb.Time).UTC()
	if pb.Response == nil {
		return nil
	}
	r := pb.GetResponse()
	e.Response = Response{
		UDID:         r.GetUdid(),
		UserID:       strPtr(r.GetUserId()),
		EnrollmentID: strPtr(r.GetEnrollmentId()),
		Status:       r.GetStatus(),
		RequestType:  r.GetRequestType(),
		CommandUUID:  r.GetCommandUuid(),
	}
	e.Raw = pb.GetRaw()
	e.Params = pb.GetParams()
	return nil
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

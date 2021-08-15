// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.3
// source: checkin.proto

package checkinproto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Event struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id      string            `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Time    int64             `protobuf:"varint,2,opt,name=time,proto3" json:"time,omitempty"`
	Command *Command          `protobuf:"bytes,3,opt,name=command,proto3" json:"command,omitempty"`
	Raw     []byte            `protobuf:"bytes,4,opt,name=raw,proto3" json:"raw,omitempty"`
	Params  map[string]string `protobuf:"bytes,5,rep,name=params,proto3" json:"params,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Event) Reset() {
	*x = Event{}
	if protoimpl.UnsafeEnabled {
		mi := &file_checkin_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Event) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Event) ProtoMessage() {}

func (x *Event) ProtoReflect() protoreflect.Message {
	mi := &file_checkin_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Event.ProtoReflect.Descriptor instead.
func (*Event) Descriptor() ([]byte, []int) {
	return file_checkin_proto_rawDescGZIP(), []int{0}
}

func (x *Event) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Event) GetTime() int64 {
	if x != nil {
		return x.Time
	}
	return 0
}

func (x *Event) GetCommand() *Command {
	if x != nil {
		return x.Command
	}
	return nil
}

func (x *Event) GetRaw() []byte {
	if x != nil {
		return x.Raw
	}
	return nil
}

func (x *Event) GetParams() map[string]string {
	if x != nil {
		return x.Params
	}
	return nil
}

type Command struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MessageType  string        `protobuf:"bytes,1,opt,name=message_type,json=messageType,proto3" json:"message_type,omitempty"`
	Topic        string        `protobuf:"bytes,2,opt,name=topic,proto3" json:"topic,omitempty"`
	Udid         string        `protobuf:"bytes,3,opt,name=udid,proto3" json:"udid,omitempty"`
	Authenticate *Authenticate `protobuf:"bytes,4,opt,name=authenticate,proto3" json:"authenticate,omitempty"`
	TokenUpdate  *TokenUpdate  `protobuf:"bytes,5,opt,name=token_update,json=tokenUpdate,proto3" json:"token_update,omitempty"`
	EnrollmentId string        `protobuf:"bytes,6,opt,name=enrollment_id,json=enrollmentId,proto3" json:"enrollment_id,omitempty"`
}

func (x *Command) Reset() {
	*x = Command{}
	if protoimpl.UnsafeEnabled {
		mi := &file_checkin_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Command) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Command) ProtoMessage() {}

func (x *Command) ProtoReflect() protoreflect.Message {
	mi := &file_checkin_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Command.ProtoReflect.Descriptor instead.
func (*Command) Descriptor() ([]byte, []int) {
	return file_checkin_proto_rawDescGZIP(), []int{1}
}

func (x *Command) GetMessageType() string {
	if x != nil {
		return x.MessageType
	}
	return ""
}

func (x *Command) GetTopic() string {
	if x != nil {
		return x.Topic
	}
	return ""
}

func (x *Command) GetUdid() string {
	if x != nil {
		return x.Udid
	}
	return ""
}

func (x *Command) GetAuthenticate() *Authenticate {
	if x != nil {
		return x.Authenticate
	}
	return nil
}

func (x *Command) GetTokenUpdate() *TokenUpdate {
	if x != nil {
		return x.TokenUpdate
	}
	return nil
}

func (x *Command) GetEnrollmentId() string {
	if x != nil {
		return x.EnrollmentId
	}
	return ""
}

type Authenticate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OsVersion    string `protobuf:"bytes,1,opt,name=os_version,json=osVersion,proto3" json:"os_version,omitempty"`
	BuildVersion string `protobuf:"bytes,2,opt,name=build_version,json=buildVersion,proto3" json:"build_version,omitempty"`
	ProductName  string `protobuf:"bytes,3,opt,name=product_name,json=productName,proto3" json:"product_name,omitempty"`
	SerialNumber string `protobuf:"bytes,4,opt,name=serial_number,json=serialNumber,proto3" json:"serial_number,omitempty"`
	Imei         string `protobuf:"bytes,5,opt,name=imei,proto3" json:"imei,omitempty"`
	Meid         string `protobuf:"bytes,6,opt,name=meid,proto3" json:"meid,omitempty"`
	DeviceName   string `protobuf:"bytes,7,opt,name=device_name,json=deviceName,proto3" json:"device_name,omitempty"`
	Challenge    []byte `protobuf:"bytes,8,opt,name=challenge,proto3" json:"challenge,omitempty"`
	Model        string `protobuf:"bytes,9,opt,name=model,proto3" json:"model,omitempty"`
	ModelName    string `protobuf:"bytes,10,opt,name=model_name,json=modelName,proto3" json:"model_name,omitempty"`
}

func (x *Authenticate) Reset() {
	*x = Authenticate{}
	if protoimpl.UnsafeEnabled {
		mi := &file_checkin_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Authenticate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Authenticate) ProtoMessage() {}

func (x *Authenticate) ProtoReflect() protoreflect.Message {
	mi := &file_checkin_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Authenticate.ProtoReflect.Descriptor instead.
func (*Authenticate) Descriptor() ([]byte, []int) {
	return file_checkin_proto_rawDescGZIP(), []int{2}
}

func (x *Authenticate) GetOsVersion() string {
	if x != nil {
		return x.OsVersion
	}
	return ""
}

func (x *Authenticate) GetBuildVersion() string {
	if x != nil {
		return x.BuildVersion
	}
	return ""
}

func (x *Authenticate) GetProductName() string {
	if x != nil {
		return x.ProductName
	}
	return ""
}

func (x *Authenticate) GetSerialNumber() string {
	if x != nil {
		return x.SerialNumber
	}
	return ""
}

func (x *Authenticate) GetImei() string {
	if x != nil {
		return x.Imei
	}
	return ""
}

func (x *Authenticate) GetMeid() string {
	if x != nil {
		return x.Meid
	}
	return ""
}

func (x *Authenticate) GetDeviceName() string {
	if x != nil {
		return x.DeviceName
	}
	return ""
}

func (x *Authenticate) GetChallenge() []byte {
	if x != nil {
		return x.Challenge
	}
	return nil
}

func (x *Authenticate) GetModel() string {
	if x != nil {
		return x.Model
	}
	return ""
}

func (x *Authenticate) GetModelName() string {
	if x != nil {
		return x.ModelName
	}
	return ""
}

type TokenUpdate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token                 []byte `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	PushMagic             string `protobuf:"bytes,2,opt,name=push_magic,json=pushMagic,proto3" json:"push_magic,omitempty"`
	UnlockToken           []byte `protobuf:"bytes,3,opt,name=unlock_token,json=unlockToken,proto3" json:"unlock_token,omitempty"`
	AwaitingConfiguration bool   `protobuf:"varint,4,opt,name=awaiting_configuration,json=awaitingConfiguration,proto3" json:"awaiting_configuration,omitempty"`
	UserId                string `protobuf:"bytes,5,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	UserLongName          string `protobuf:"bytes,6,opt,name=user_long_name,json=userLongName,proto3" json:"user_long_name,omitempty"`
	UserShortName         string `protobuf:"bytes,7,opt,name=user_short_name,json=userShortName,proto3" json:"user_short_name,omitempty"`
	NotOnConsole          bool   `protobuf:"varint,8,opt,name=not_on_console,json=notOnConsole,proto3" json:"not_on_console,omitempty"`
}

func (x *TokenUpdate) Reset() {
	*x = TokenUpdate{}
	if protoimpl.UnsafeEnabled {
		mi := &file_checkin_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TokenUpdate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TokenUpdate) ProtoMessage() {}

func (x *TokenUpdate) ProtoReflect() protoreflect.Message {
	mi := &file_checkin_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TokenUpdate.ProtoReflect.Descriptor instead.
func (*TokenUpdate) Descriptor() ([]byte, []int) {
	return file_checkin_proto_rawDescGZIP(), []int{3}
}

func (x *TokenUpdate) GetToken() []byte {
	if x != nil {
		return x.Token
	}
	return nil
}

func (x *TokenUpdate) GetPushMagic() string {
	if x != nil {
		return x.PushMagic
	}
	return ""
}

func (x *TokenUpdate) GetUnlockToken() []byte {
	if x != nil {
		return x.UnlockToken
	}
	return nil
}

func (x *TokenUpdate) GetAwaitingConfiguration() bool {
	if x != nil {
		return x.AwaitingConfiguration
	}
	return false
}

func (x *TokenUpdate) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *TokenUpdate) GetUserLongName() string {
	if x != nil {
		return x.UserLongName
	}
	return ""
}

func (x *TokenUpdate) GetUserShortName() string {
	if x != nil {
		return x.UserShortName
	}
	return ""
}

func (x *TokenUpdate) GetNotOnConsole() bool {
	if x != nil {
		return x.NotOnConsole
	}
	return false
}

var File_checkin_proto protoreflect.FileDescriptor

var file_checkin_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x69, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x0c, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x69, 0x6e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xe2, 0x01,
	0x0a, 0x05, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x2f, 0x0a, 0x07, 0x63,
	0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x63,
	0x68, 0x65, 0x63, 0x6b, 0x69, 0x6e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x6f, 0x6d, 0x6d,
	0x61, 0x6e, 0x64, 0x52, 0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x12, 0x10, 0x0a, 0x03,
	0x72, 0x61, 0x77, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x72, 0x61, 0x77, 0x12, 0x37,
	0x0a, 0x06, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1f,
	0x2e, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x69, 0x6e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x2e, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52,
	0x06, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x1a, 0x39, 0x0a, 0x0b, 0x50, 0x61, 0x72, 0x61, 0x6d,
	0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02,
	0x38, 0x01, 0x22, 0xf9, 0x01, 0x0a, 0x07, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x12, 0x21,
	0x0a, 0x0c, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x64, 0x69, 0x64, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x75, 0x64, 0x69, 0x64, 0x12, 0x3e, 0x0a, 0x0c, 0x61,
	0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x69, 0x6e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x41, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x65, 0x52, 0x0c, 0x61,
	0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x65, 0x12, 0x3c, 0x0a, 0x0c, 0x74,
	0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x19, 0x2e, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x69, 0x6e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x0b, 0x74, 0x6f,
	0x6b, 0x65, 0x6e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x65, 0x6e, 0x72,
	0x6f, 0x6c, 0x6c, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0c, 0x65, 0x6e, 0x72, 0x6f, 0x6c, 0x6c, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x22, 0xb6,
	0x02, 0x0a, 0x0c, 0x41, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x65, 0x12,
	0x1d, 0x0a, 0x0a, 0x6f, 0x73, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x09, 0x6f, 0x73, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x23,
	0x0a, 0x0d, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x56, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x12, 0x21, 0x0a, 0x0c, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x5f, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x70, 0x72, 0x6f, 0x64, 0x75,
	0x63, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c,
	0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x73,
	0x65, 0x72, 0x69, 0x61, 0x6c, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x69,
	0x6d, 0x65, 0x69, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x69, 0x6d, 0x65, 0x69, 0x12,
	0x12, 0x0a, 0x04, 0x6d, 0x65, 0x69, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6d,
	0x65, 0x69, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65,
	0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x63, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67,
	0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x63, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e,
	0x67, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x18, 0x09, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x12, 0x1d, 0x0a, 0x0a, 0x6d, 0x6f, 0x64, 0x65,
	0x6c, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6d, 0x6f,
	0x64, 0x65, 0x6c, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0xa9, 0x02, 0x0a, 0x0b, 0x54, 0x6f, 0x6b, 0x65,
	0x6e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x1d, 0x0a,
	0x0a, 0x70, 0x75, 0x73, 0x68, 0x5f, 0x6d, 0x61, 0x67, 0x69, 0x63, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x70, 0x75, 0x73, 0x68, 0x4d, 0x61, 0x67, 0x69, 0x63, 0x12, 0x21, 0x0a, 0x0c,
	0x75, 0x6e, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x0b, 0x75, 0x6e, 0x6c, 0x6f, 0x63, 0x6b, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12,
	0x35, 0x0a, 0x16, 0x61, 0x77, 0x61, 0x69, 0x74, 0x69, 0x6e, 0x67, 0x5f, 0x63, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x15, 0x61, 0x77, 0x61, 0x69, 0x74, 0x69, 0x6e, 0x67, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69,
	0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12,
	0x24, 0x0a, 0x0e, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x6c, 0x6f, 0x6e, 0x67, 0x5f, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x75, 0x73, 0x65, 0x72, 0x4c, 0x6f, 0x6e,
	0x67, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x26, 0x0a, 0x0f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x73, 0x68,
	0x6f, 0x72, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d,
	0x75, 0x73, 0x65, 0x72, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x24, 0x0a,
	0x0e, 0x6e, 0x6f, 0x74, 0x5f, 0x6f, 0x6e, 0x5f, 0x63, 0x6f, 0x6e, 0x73, 0x6f, 0x6c, 0x65, 0x18,
	0x08, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0c, 0x6e, 0x6f, 0x74, 0x4f, 0x6e, 0x43, 0x6f, 0x6e, 0x73,
	0x6f, 0x6c, 0x65, 0x42, 0x38, 0x5a, 0x36, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x6d, 0x69, 0x63, 0x72, 0x6f, 0x6d, 0x64, 0x6d, 0x2f, 0x6d, 0x69, 0x63, 0x72, 0x6f,
	0x6d, 0x64, 0x6d, 0x2f, 0x6d, 0x64, 0x6d, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c,
	0x2f, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x69, 0x6e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_checkin_proto_rawDescOnce sync.Once
	file_checkin_proto_rawDescData = file_checkin_proto_rawDesc
)

func file_checkin_proto_rawDescGZIP() []byte {
	file_checkin_proto_rawDescOnce.Do(func() {
		file_checkin_proto_rawDescData = protoimpl.X.CompressGZIP(file_checkin_proto_rawDescData)
	})
	return file_checkin_proto_rawDescData
}

var file_checkin_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_checkin_proto_goTypes = []interface{}{
	(*Event)(nil),        // 0: checkinproto.Event
	(*Command)(nil),      // 1: checkinproto.Command
	(*Authenticate)(nil), // 2: checkinproto.Authenticate
	(*TokenUpdate)(nil),  // 3: checkinproto.TokenUpdate
	nil,                  // 4: checkinproto.Event.ParamsEntry
}
var file_checkin_proto_depIdxs = []int32{
	1, // 0: checkinproto.Event.command:type_name -> checkinproto.Command
	4, // 1: checkinproto.Event.params:type_name -> checkinproto.Event.ParamsEntry
	2, // 2: checkinproto.Command.authenticate:type_name -> checkinproto.Authenticate
	3, // 3: checkinproto.Command.token_update:type_name -> checkinproto.TokenUpdate
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_checkin_proto_init() }
func file_checkin_proto_init() {
	if File_checkin_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_checkin_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Event); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_checkin_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Command); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_checkin_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Authenticate); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_checkin_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TokenUpdate); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_checkin_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_checkin_proto_goTypes,
		DependencyIndexes: file_checkin_proto_depIdxs,
		MessageInfos:      file_checkin_proto_msgTypes,
	}.Build()
	File_checkin_proto = out.File
	file_checkin_proto_rawDesc = nil
	file_checkin_proto_goTypes = nil
	file_checkin_proto_depIdxs = nil
}

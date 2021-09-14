// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: device.proto

package deviceproto

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

type Device struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uuid                   string `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Udid                   string `protobuf:"bytes,2,opt,name=udid,proto3" json:"udid,omitempty"`
	SerialNumber           string `protobuf:"bytes,3,opt,name=serial_number,json=serialNumber,proto3" json:"serial_number,omitempty"`
	OsVersion              string `protobuf:"bytes,4,opt,name=os_version,json=osVersion,proto3" json:"os_version,omitempty"`
	BuildVersion           string `protobuf:"bytes,5,opt,name=build_version,json=buildVersion,proto3" json:"build_version,omitempty"`
	ProductName            string `protobuf:"bytes,6,opt,name=product_name,json=productName,proto3" json:"product_name,omitempty"`
	Imei                   string `protobuf:"bytes,7,opt,name=imei,proto3" json:"imei,omitempty"`
	Meid                   string `protobuf:"bytes,8,opt,name=meid,proto3" json:"meid,omitempty"`
	Token                  string `protobuf:"bytes,9,opt,name=token,proto3" json:"token,omitempty"`
	PushMagic              string `protobuf:"bytes,10,opt,name=push_magic,json=pushMagic,proto3" json:"push_magic,omitempty"`
	MdmTopic               string `protobuf:"bytes,11,opt,name=mdm_topic,json=mdmTopic,proto3" json:"mdm_topic,omitempty"`
	UnlockToken            string `protobuf:"bytes,12,opt,name=unlock_token,json=unlockToken,proto3" json:"unlock_token,omitempty"`
	Enrolled               bool   `protobuf:"varint,13,opt,name=enrolled,proto3" json:"enrolled,omitempty"`
	AwaitingConfiguration  bool   `protobuf:"varint,14,opt,name=awaiting_configuration,json=awaitingConfiguration,proto3" json:"awaiting_configuration,omitempty"`
	DeviceName             string `protobuf:"bytes,15,opt,name=device_name,json=deviceName,proto3" json:"device_name,omitempty"`
	Model                  string `protobuf:"bytes,16,opt,name=model,proto3" json:"model,omitempty"`
	ModelName              string `protobuf:"bytes,17,opt,name=model_name,json=modelName,proto3" json:"model_name,omitempty"`
	Description            string `protobuf:"bytes,18,opt,name=description,proto3" json:"description,omitempty"`
	Color                  string `protobuf:"bytes,19,opt,name=color,proto3" json:"color,omitempty"`
	AssetTag               string `protobuf:"bytes,20,opt,name=asset_tag,json=assetTag,proto3" json:"asset_tag,omitempty"`
	DepDevice              bool   `protobuf:"varint,21,opt,name=dep_device,json=depDevice,proto3" json:"dep_device,omitempty"`
	DepProfileStatus       string `protobuf:"bytes,22,opt,name=dep_profile_status,json=depProfileStatus,proto3" json:"dep_profile_status,omitempty"`
	DepProfileUuid         string `protobuf:"bytes,23,opt,name=dep_profile_uuid,json=depProfileUuid,proto3" json:"dep_profile_uuid,omitempty"`
	DepProfileAssignTime   int64  `protobuf:"varint,24,opt,name=dep_profile_assign_time,json=depProfileAssignTime,proto3" json:"dep_profile_assign_time,omitempty"`
	DepProfilePushTime     int64  `protobuf:"varint,25,opt,name=dep_profile_push_time,json=depProfilePushTime,proto3" json:"dep_profile_push_time,omitempty"`
	DepProfileAssignedDate int64  `protobuf:"varint,26,opt,name=dep_profile_assigned_date,json=depProfileAssignedDate,proto3" json:"dep_profile_assigned_date,omitempty"`
	DepProfileAssignedBy   string `protobuf:"bytes,27,opt,name=dep_profile_assigned_by,json=depProfileAssignedBy,proto3" json:"dep_profile_assigned_by,omitempty"`
	LastSeen               int64  `protobuf:"varint,28,opt,name=last_seen,json=lastSeen,proto3" json:"last_seen,omitempty"`
	LastQueryResponse      []byte `protobuf:"bytes,29,opt,name=last_query_response,json=lastQueryResponse,proto3" json:"last_query_response,omitempty"`
	BootstrapToken         string `protobuf:"bytes,30,opt,name=bootstrap_token,json=bootstrapToken,proto3" json:"bootstrap_token,omitempty"`
}

func (x *Device) Reset() {
	*x = Device{}
	if protoimpl.UnsafeEnabled {
		mi := &file_device_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Device) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Device) ProtoMessage() {}

func (x *Device) ProtoReflect() protoreflect.Message {
	mi := &file_device_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Device.ProtoReflect.Descriptor instead.
func (*Device) Descriptor() ([]byte, []int) {
	return file_device_proto_rawDescGZIP(), []int{0}
}

func (x *Device) GetUuid() string {
	if x != nil {
		return x.Uuid
	}
	return ""
}

func (x *Device) GetUdid() string {
	if x != nil {
		return x.Udid
	}
	return ""
}

func (x *Device) GetSerialNumber() string {
	if x != nil {
		return x.SerialNumber
	}
	return ""
}

func (x *Device) GetOsVersion() string {
	if x != nil {
		return x.OsVersion
	}
	return ""
}

func (x *Device) GetBuildVersion() string {
	if x != nil {
		return x.BuildVersion
	}
	return ""
}

func (x *Device) GetProductName() string {
	if x != nil {
		return x.ProductName
	}
	return ""
}

func (x *Device) GetImei() string {
	if x != nil {
		return x.Imei
	}
	return ""
}

func (x *Device) GetMeid() string {
	if x != nil {
		return x.Meid
	}
	return ""
}

func (x *Device) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *Device) GetPushMagic() string {
	if x != nil {
		return x.PushMagic
	}
	return ""
}

func (x *Device) GetMdmTopic() string {
	if x != nil {
		return x.MdmTopic
	}
	return ""
}

func (x *Device) GetUnlockToken() string {
	if x != nil {
		return x.UnlockToken
	}
	return ""
}

func (x *Device) GetEnrolled() bool {
	if x != nil {
		return x.Enrolled
	}
	return false
}

func (x *Device) GetAwaitingConfiguration() bool {
	if x != nil {
		return x.AwaitingConfiguration
	}
	return false
}

func (x *Device) GetDeviceName() string {
	if x != nil {
		return x.DeviceName
	}
	return ""
}

func (x *Device) GetModel() string {
	if x != nil {
		return x.Model
	}
	return ""
}

func (x *Device) GetModelName() string {
	if x != nil {
		return x.ModelName
	}
	return ""
}

func (x *Device) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Device) GetColor() string {
	if x != nil {
		return x.Color
	}
	return ""
}

func (x *Device) GetAssetTag() string {
	if x != nil {
		return x.AssetTag
	}
	return ""
}

func (x *Device) GetDepDevice() bool {
	if x != nil {
		return x.DepDevice
	}
	return false
}

func (x *Device) GetDepProfileStatus() string {
	if x != nil {
		return x.DepProfileStatus
	}
	return ""
}

func (x *Device) GetDepProfileUuid() string {
	if x != nil {
		return x.DepProfileUuid
	}
	return ""
}

func (x *Device) GetDepProfileAssignTime() int64 {
	if x != nil {
		return x.DepProfileAssignTime
	}
	return 0
}

func (x *Device) GetDepProfilePushTime() int64 {
	if x != nil {
		return x.DepProfilePushTime
	}
	return 0
}

func (x *Device) GetDepProfileAssignedDate() int64 {
	if x != nil {
		return x.DepProfileAssignedDate
	}
	return 0
}

func (x *Device) GetDepProfileAssignedBy() string {
	if x != nil {
		return x.DepProfileAssignedBy
	}
	return ""
}

func (x *Device) GetLastSeen() int64 {
	if x != nil {
		return x.LastSeen
	}
	return 0
}

func (x *Device) GetLastQueryResponse() []byte {
	if x != nil {
		return x.LastQueryResponse
	}
	return nil
}

func (x *Device) GetBootstrapToken() string {
	if x != nil {
		return x.BootstrapToken
	}
	return ""
}

var File_device_proto protoreflect.FileDescriptor

var file_device_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b,
	0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa0, 0x08, 0x0a, 0x06,
	0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x75, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x75, 0x75, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x64,
	0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x75, 0x64, 0x69, 0x64, 0x12, 0x23,
	0x0a, 0x0d, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x4e, 0x75, 0x6d,
	0x62, 0x65, 0x72, 0x12, 0x1d, 0x0a, 0x0a, 0x6f, 0x73, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6f, 0x73, 0x56, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x12, 0x23, 0x0a, 0x0d, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x5f, 0x76, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x62, 0x75, 0x69, 0x6c, 0x64,
	0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x21, 0x0a, 0x0c, 0x70, 0x72, 0x6f, 0x64, 0x75,
	0x63, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x70,
	0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x69, 0x6d,
	0x65, 0x69, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x69, 0x6d, 0x65, 0x69, 0x12, 0x12,
	0x0a, 0x04, 0x6d, 0x65, 0x69, 0x64, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6d, 0x65,
	0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x09, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x75, 0x73, 0x68,
	0x5f, 0x6d, 0x61, 0x67, 0x69, 0x63, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x70, 0x75,
	0x73, 0x68, 0x4d, 0x61, 0x67, 0x69, 0x63, 0x12, 0x1b, 0x0a, 0x09, 0x6d, 0x64, 0x6d, 0x5f, 0x74,
	0x6f, 0x70, 0x69, 0x63, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6d, 0x64, 0x6d, 0x54,
	0x6f, 0x70, 0x69, 0x63, 0x12, 0x21, 0x0a, 0x0c, 0x75, 0x6e, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x74,
	0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x75, 0x6e, 0x6c, 0x6f,
	0x63, 0x6b, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x65, 0x6e, 0x72, 0x6f, 0x6c,
	0x6c, 0x65, 0x64, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x65, 0x6e, 0x72, 0x6f, 0x6c,
	0x6c, 0x65, 0x64, 0x12, 0x35, 0x0a, 0x16, 0x61, 0x77, 0x61, 0x69, 0x74, 0x69, 0x6e, 0x67, 0x5f,
	0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x0e, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x15, 0x61, 0x77, 0x61, 0x69, 0x74, 0x69, 0x6e, 0x67, 0x43, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1f, 0x0a, 0x0b, 0x64, 0x65,
	0x76, 0x69, 0x63, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0a, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x6d,
	0x6f, 0x64, 0x65, 0x6c, 0x18, 0x10, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6d, 0x6f, 0x64, 0x65,
	0x6c, 0x12, 0x1d, 0x0a, 0x0a, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x11, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x4e, 0x61, 0x6d, 0x65,
	0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18,
	0x12, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x18, 0x13, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x12, 0x1b, 0x0a, 0x09, 0x61, 0x73, 0x73, 0x65,
	0x74, 0x5f, 0x74, 0x61, 0x67, 0x18, 0x14, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x61, 0x73, 0x73,
	0x65, 0x74, 0x54, 0x61, 0x67, 0x12, 0x1d, 0x0a, 0x0a, 0x64, 0x65, 0x70, 0x5f, 0x64, 0x65, 0x76,
	0x69, 0x63, 0x65, 0x18, 0x15, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x64, 0x65, 0x70, 0x44, 0x65,
	0x76, 0x69, 0x63, 0x65, 0x12, 0x2c, 0x0a, 0x12, 0x64, 0x65, 0x70, 0x5f, 0x70, 0x72, 0x6f, 0x66,
	0x69, 0x6c, 0x65, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x16, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x10, 0x64, 0x65, 0x70, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x12, 0x28, 0x0a, 0x10, 0x64, 0x65, 0x70, 0x5f, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c,
	0x65, 0x5f, 0x75, 0x75, 0x69, 0x64, 0x18, 0x17, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x64, 0x65,
	0x70, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x55, 0x75, 0x69, 0x64, 0x12, 0x35, 0x0a, 0x17,
	0x64, 0x65, 0x70, 0x5f, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x61, 0x73, 0x73, 0x69,
	0x67, 0x6e, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x18, 0x20, 0x01, 0x28, 0x03, 0x52, 0x14, 0x64,
	0x65, 0x70, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x54,
	0x69, 0x6d, 0x65, 0x12, 0x31, 0x0a, 0x15, 0x64, 0x65, 0x70, 0x5f, 0x70, 0x72, 0x6f, 0x66, 0x69,
	0x6c, 0x65, 0x5f, 0x70, 0x75, 0x73, 0x68, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x19, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x12, 0x64, 0x65, 0x70, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x50, 0x75,
	0x73, 0x68, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x39, 0x0a, 0x19, 0x64, 0x65, 0x70, 0x5f, 0x70, 0x72,
	0x6f, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x61, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x5f, 0x64,
	0x61, 0x74, 0x65, 0x18, 0x1a, 0x20, 0x01, 0x28, 0x03, 0x52, 0x16, 0x64, 0x65, 0x70, 0x50, 0x72,
	0x6f, 0x66, 0x69, 0x6c, 0x65, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x44, 0x61, 0x74,
	0x65, 0x12, 0x35, 0x0a, 0x17, 0x64, 0x65, 0x70, 0x5f, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65,
	0x5f, 0x61, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x5f, 0x62, 0x79, 0x18, 0x1b, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x14, 0x64, 0x65, 0x70, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x41, 0x73,
	0x73, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x42, 0x79, 0x12, 0x1b, 0x0a, 0x09, 0x6c, 0x61, 0x73, 0x74,
	0x5f, 0x73, 0x65, 0x65, 0x6e, 0x18, 0x1c, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x6c, 0x61, 0x73,
	0x74, 0x53, 0x65, 0x65, 0x6e, 0x12, 0x2e, 0x0a, 0x13, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x71, 0x75,
	0x65, 0x72, 0x79, 0x5f, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x18, 0x1d, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x11, 0x6c, 0x61, 0x73, 0x74, 0x51, 0x75, 0x65, 0x72, 0x79, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x27, 0x0a, 0x0f, 0x62, 0x6f, 0x6f, 0x74, 0x73, 0x74, 0x72,
	0x61, 0x70, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x1e, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e,
	0x62, 0x6f, 0x6f, 0x74, 0x73, 0x74, 0x72, 0x61, 0x70, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x42, 0x0f,
	0x5a, 0x0d, 0x2f, 0x3b, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_device_proto_rawDescOnce sync.Once
	file_device_proto_rawDescData = file_device_proto_rawDesc
)

func file_device_proto_rawDescGZIP() []byte {
	file_device_proto_rawDescOnce.Do(func() {
		file_device_proto_rawDescData = protoimpl.X.CompressGZIP(file_device_proto_rawDescData)
	})
	return file_device_proto_rawDescData
}

var file_device_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_device_proto_goTypes = []interface{}{
	(*Device)(nil), // 0: deviceproto.Device
}
var file_device_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_device_proto_init() }
func file_device_proto_init() {
	if File_device_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_device_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Device); i {
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
			RawDescriptor: file_device_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_device_proto_goTypes,
		DependencyIndexes: file_device_proto_depIdxs,
		MessageInfos:      file_device_proto_msgTypes,
	}.Build()
	File_device_proto = out.File
	file_device_proto_rawDesc = nil
	file_device_proto_goTypes = nil
	file_device_proto_depIdxs = nil
}

// Code generated by protoc-gen-go. DO NOT EDIT.
// source: command.proto

/*
Package commandproto is a generated protocol buffer package.

It is generated from these files:
	command.proto

It has these top-level messages:
	Event
	Payload
	Command
	ScheduleOSUpdate
	OSUpdate
	ScheduleOSUpdateScan
	AccountConfiguration
	AutoSetupAdminAccounts
	DeviceInformation
	InstallProfile
	RemoveProfile
	DeleteUser
	InstallApplication
	Settings
	Setting
	DeviceNameSetting
	HostnameSetting
*/
package commandproto

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Event struct {
	Id         string   `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Time       int64    `protobuf:"varint,2,opt,name=time" json:"time,omitempty"`
	Payload    *Payload `protobuf:"bytes,3,opt,name=payload" json:"payload,omitempty"`
	DeviceUdid string   `protobuf:"bytes,4,opt,name=device_udid,json=deviceUdid" json:"device_udid,omitempty"`
}

func (m *Event) Reset()                    { *m = Event{} }
func (m *Event) String() string            { return proto.CompactTextString(m) }
func (*Event) ProtoMessage()               {}
func (*Event) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Event) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Event) GetTime() int64 {
	if m != nil {
		return m.Time
	}
	return 0
}

func (m *Event) GetPayload() *Payload {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *Event) GetDeviceUdid() string {
	if m != nil {
		return m.DeviceUdid
	}
	return ""
}

type Payload struct {
	CommandUuid string   `protobuf:"bytes,1,opt,name=command_uuid,json=commandUuid" json:"command_uuid,omitempty"`
	Command     *Command `protobuf:"bytes,2,opt,name=command" json:"command,omitempty"`
}

func (m *Payload) Reset()                    { *m = Payload{} }
func (m *Payload) String() string            { return proto.CompactTextString(m) }
func (*Payload) ProtoMessage()               {}
func (*Payload) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Payload) GetCommandUuid() string {
	if m != nil {
		return m.CommandUuid
	}
	return ""
}

func (m *Payload) GetCommand() *Command {
	if m != nil {
		return m.Command
	}
	return nil
}

type Command struct {
	RequestType          string                `protobuf:"bytes,1,opt,name=request_type,json=requestType" json:"request_type,omitempty"`
	DeviceInformation    *DeviceInformation    `protobuf:"bytes,2,opt,name=device_information,json=deviceInformation" json:"device_information,omitempty"`
	InstallProfile       *InstallProfile       `protobuf:"bytes,3,opt,name=install_profile,json=installProfile" json:"install_profile,omitempty"`
	InstallApplication   *InstallApplication   `protobuf:"bytes,4,opt,name=install_application,json=installApplication" json:"install_application,omitempty"`
	AccountConfiguration *AccountConfiguration `protobuf:"bytes,5,opt,name=account_configuration,json=accountConfiguration" json:"account_configuration,omitempty"`
	ScheduleOsUpdate     *ScheduleOSUpdate     `protobuf:"bytes,6,opt,name=schedule_os_update,json=scheduleOsUpdate" json:"schedule_os_update,omitempty"`
	ScheduleOsUpdateScan *ScheduleOSUpdateScan `protobuf:"bytes,7,opt,name=schedule_os_update_scan,json=scheduleOsUpdateScan" json:"schedule_os_update_scan,omitempty"`
	RemoveProfile        *RemoveProfile        `protobuf:"bytes,8,opt,name=remove_profile,json=removeProfile" json:"remove_profile,omitempty"`
	DeleteUser           *DeleteUser           `protobuf:"bytes,9,opt,name=delete_user,json=deleteUser" json:"delete_user,omitempty"`
	Settings             *Settings             `protobuf:"bytes,10,opt,name=settings" json:"settings,omitempty"`
}

func (m *Command) Reset()                    { *m = Command{} }
func (m *Command) String() string            { return proto.CompactTextString(m) }
func (*Command) ProtoMessage()               {}
func (*Command) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Command) GetRequestType() string {
	if m != nil {
		return m.RequestType
	}
	return ""
}

func (m *Command) GetDeviceInformation() *DeviceInformation {
	if m != nil {
		return m.DeviceInformation
	}
	return nil
}

func (m *Command) GetInstallProfile() *InstallProfile {
	if m != nil {
		return m.InstallProfile
	}
	return nil
}

func (m *Command) GetInstallApplication() *InstallApplication {
	if m != nil {
		return m.InstallApplication
	}
	return nil
}

func (m *Command) GetAccountConfiguration() *AccountConfiguration {
	if m != nil {
		return m.AccountConfiguration
	}
	return nil
}

func (m *Command) GetScheduleOsUpdate() *ScheduleOSUpdate {
	if m != nil {
		return m.ScheduleOsUpdate
	}
	return nil
}

func (m *Command) GetScheduleOsUpdateScan() *ScheduleOSUpdateScan {
	if m != nil {
		return m.ScheduleOsUpdateScan
	}
	return nil
}

func (m *Command) GetRemoveProfile() *RemoveProfile {
	if m != nil {
		return m.RemoveProfile
	}
	return nil
}

func (m *Command) GetDeleteUser() *DeleteUser {
	if m != nil {
		return m.DeleteUser
	}
	return nil
}

func (m *Command) GetSettings() *Settings {
	if m != nil {
		return m.Settings
	}
	return nil
}

type ScheduleOSUpdate struct {
	Updates []*OSUpdate `protobuf:"bytes,1,rep,name=updates" json:"updates,omitempty"`
}

func (m *ScheduleOSUpdate) Reset()                    { *m = ScheduleOSUpdate{} }
func (m *ScheduleOSUpdate) String() string            { return proto.CompactTextString(m) }
func (*ScheduleOSUpdate) ProtoMessage()               {}
func (*ScheduleOSUpdate) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *ScheduleOSUpdate) GetUpdates() []*OSUpdate {
	if m != nil {
		return m.Updates
	}
	return nil
}

type OSUpdate struct {
	ProductKey    string `protobuf:"bytes,1,opt,name=product_key,json=productKey" json:"product_key,omitempty"`
	InstallAction string `protobuf:"bytes,2,opt,name=install_action,json=installAction" json:"install_action,omitempty"`
}

func (m *OSUpdate) Reset()                    { *m = OSUpdate{} }
func (m *OSUpdate) String() string            { return proto.CompactTextString(m) }
func (*OSUpdate) ProtoMessage()               {}
func (*OSUpdate) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *OSUpdate) GetProductKey() string {
	if m != nil {
		return m.ProductKey
	}
	return ""
}

func (m *OSUpdate) GetInstallAction() string {
	if m != nil {
		return m.InstallAction
	}
	return ""
}

type ScheduleOSUpdateScan struct {
	Force bool `protobuf:"varint,1,opt,name=force" json:"force,omitempty"`
}

func (m *ScheduleOSUpdateScan) Reset()                    { *m = ScheduleOSUpdateScan{} }
func (m *ScheduleOSUpdateScan) String() string            { return proto.CompactTextString(m) }
func (*ScheduleOSUpdateScan) ProtoMessage()               {}
func (*ScheduleOSUpdateScan) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *ScheduleOSUpdateScan) GetForce() bool {
	if m != nil {
		return m.Force
	}
	return false
}

type AccountConfiguration struct {
	SkipPrimarySetupAccountCreation     bool                      `protobuf:"varint,1,opt,name=skip_primary_setup_account_creation,json=skipPrimarySetupAccountCreation" json:"skip_primary_setup_account_creation,omitempty"`
	SetPrimarySetupAccountAsRegularUser bool                      `protobuf:"varint,2,opt,name=set_primary_setup_account_as_regular_user,json=setPrimarySetupAccountAsRegularUser" json:"set_primary_setup_account_as_regular_user,omitempty"`
	AutoSetupAdminAccounts              []*AutoSetupAdminAccounts `protobuf:"bytes,3,rep,name=auto_setup_admin_accounts,json=autoSetupAdminAccounts" json:"auto_setup_admin_accounts,omitempty"`
}

func (m *AccountConfiguration) Reset()                    { *m = AccountConfiguration{} }
func (m *AccountConfiguration) String() string            { return proto.CompactTextString(m) }
func (*AccountConfiguration) ProtoMessage()               {}
func (*AccountConfiguration) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *AccountConfiguration) GetSkipPrimarySetupAccountCreation() bool {
	if m != nil {
		return m.SkipPrimarySetupAccountCreation
	}
	return false
}

func (m *AccountConfiguration) GetSetPrimarySetupAccountAsRegularUser() bool {
	if m != nil {
		return m.SetPrimarySetupAccountAsRegularUser
	}
	return false
}

func (m *AccountConfiguration) GetAutoSetupAdminAccounts() []*AutoSetupAdminAccounts {
	if m != nil {
		return m.AutoSetupAdminAccounts
	}
	return nil
}

type AutoSetupAdminAccounts struct {
	ShortName    string `protobuf:"bytes,1,opt,name=short_name,json=shortName" json:"short_name,omitempty"`
	FullName     string `protobuf:"bytes,2,opt,name=full_name,json=fullName" json:"full_name,omitempty"`
	PasswordHash []byte `protobuf:"bytes,3,opt,name=password_hash,json=passwordHash,proto3" json:"password_hash,omitempty"`
	Hidden       bool   `protobuf:"varint,4,opt,name=hidden" json:"hidden,omitempty"`
}

func (m *AutoSetupAdminAccounts) Reset()                    { *m = AutoSetupAdminAccounts{} }
func (m *AutoSetupAdminAccounts) String() string            { return proto.CompactTextString(m) }
func (*AutoSetupAdminAccounts) ProtoMessage()               {}
func (*AutoSetupAdminAccounts) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *AutoSetupAdminAccounts) GetShortName() string {
	if m != nil {
		return m.ShortName
	}
	return ""
}

func (m *AutoSetupAdminAccounts) GetFullName() string {
	if m != nil {
		return m.FullName
	}
	return ""
}

func (m *AutoSetupAdminAccounts) GetPasswordHash() []byte {
	if m != nil {
		return m.PasswordHash
	}
	return nil
}

func (m *AutoSetupAdminAccounts) GetHidden() bool {
	if m != nil {
		return m.Hidden
	}
	return false
}

type DeviceInformation struct {
	Queries []string `protobuf:"bytes,1,rep,name=queries" json:"queries,omitempty"`
}

func (m *DeviceInformation) Reset()                    { *m = DeviceInformation{} }
func (m *DeviceInformation) String() string            { return proto.CompactTextString(m) }
func (*DeviceInformation) ProtoMessage()               {}
func (*DeviceInformation) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *DeviceInformation) GetQueries() []string {
	if m != nil {
		return m.Queries
	}
	return nil
}

type InstallProfile struct {
	Payload []byte `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (m *InstallProfile) Reset()                    { *m = InstallProfile{} }
func (m *InstallProfile) String() string            { return proto.CompactTextString(m) }
func (*InstallProfile) ProtoMessage()               {}
func (*InstallProfile) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *InstallProfile) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

type RemoveProfile struct {
	Identifier string `protobuf:"bytes,1,opt,name=identifier" json:"identifier,omitempty"`
}

func (m *RemoveProfile) Reset()                    { *m = RemoveProfile{} }
func (m *RemoveProfile) String() string            { return proto.CompactTextString(m) }
func (*RemoveProfile) ProtoMessage()               {}
func (*RemoveProfile) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

func (m *RemoveProfile) GetIdentifier() string {
	if m != nil {
		return m.Identifier
	}
	return ""
}

type DeleteUser struct {
	Username      string `protobuf:"bytes,1,opt,name=username" json:"username,omitempty"`
	ForceDeletion bool   `protobuf:"varint,2,opt,name=force_deletion,json=forceDeletion" json:"force_deletion,omitempty"`
}

func (m *DeleteUser) Reset()                    { *m = DeleteUser{} }
func (m *DeleteUser) String() string            { return proto.CompactTextString(m) }
func (*DeleteUser) ProtoMessage()               {}
func (*DeleteUser) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

func (m *DeleteUser) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *DeleteUser) GetForceDeletion() bool {
	if m != nil {
		return m.ForceDeletion
	}
	return false
}

type InstallApplication struct {
	ItunesStoreId         int64  `protobuf:"varint,1,opt,name=itunes_store_id,json=itunesStoreId" json:"itunes_store_id,omitempty"`
	Identifier            string `protobuf:"bytes,2,opt,name=identifier" json:"identifier,omitempty"`
	ManifestUrl           string `protobuf:"bytes,3,opt,name=manifest_url,json=manifestUrl" json:"manifest_url,omitempty"`
	ManagementFlags       int64  `protobuf:"varint,4,opt,name=management_flags,json=managementFlags" json:"management_flags,omitempty"`
	NotManaged            bool   `protobuf:"varint,5,opt,name=not_managed,json=notManaged" json:"not_managed,omitempty"`
	ChangeManagementState string `protobuf:"bytes,6,opt,name=change_management_state,json=changeManagementState" json:"change_management_state,omitempty"`
}

func (m *InstallApplication) Reset()                    { *m = InstallApplication{} }
func (m *InstallApplication) String() string            { return proto.CompactTextString(m) }
func (*InstallApplication) ProtoMessage()               {}
func (*InstallApplication) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{12} }

func (m *InstallApplication) GetItunesStoreId() int64 {
	if m != nil {
		return m.ItunesStoreId
	}
	return 0
}

func (m *InstallApplication) GetIdentifier() string {
	if m != nil {
		return m.Identifier
	}
	return ""
}

func (m *InstallApplication) GetManifestUrl() string {
	if m != nil {
		return m.ManifestUrl
	}
	return ""
}

func (m *InstallApplication) GetManagementFlags() int64 {
	if m != nil {
		return m.ManagementFlags
	}
	return 0
}

func (m *InstallApplication) GetNotManaged() bool {
	if m != nil {
		return m.NotManaged
	}
	return false
}

func (m *InstallApplication) GetChangeManagementState() string {
	if m != nil {
		return m.ChangeManagementState
	}
	return ""
}

type Settings struct {
	Settings []*Setting `protobuf:"bytes,1,rep,name=settings" json:"settings,omitempty"`
}

func (m *Settings) Reset()                    { *m = Settings{} }
func (m *Settings) String() string            { return proto.CompactTextString(m) }
func (*Settings) ProtoMessage()               {}
func (*Settings) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{13} }

func (m *Settings) GetSettings() []*Setting {
	if m != nil {
		return m.Settings
	}
	return nil
}

type Setting struct {
	Item       string             `protobuf:"bytes,1,opt,name=item" json:"item,omitempty"`
	DeviceName *DeviceNameSetting `protobuf:"bytes,2,opt,name=device_name,json=deviceName" json:"device_name,omitempty"`
	Hostname   *HostnameSetting   `protobuf:"bytes,3,opt,name=hostname" json:"hostname,omitempty"`
}

func (m *Setting) Reset()                    { *m = Setting{} }
func (m *Setting) String() string            { return proto.CompactTextString(m) }
func (*Setting) ProtoMessage()               {}
func (*Setting) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{14} }

func (m *Setting) GetItem() string {
	if m != nil {
		return m.Item
	}
	return ""
}

func (m *Setting) GetDeviceName() *DeviceNameSetting {
	if m != nil {
		return m.DeviceName
	}
	return nil
}

func (m *Setting) GetHostname() *HostnameSetting {
	if m != nil {
		return m.Hostname
	}
	return nil
}

type DeviceNameSetting struct {
	DeviceName string `protobuf:"bytes,1,opt,name=device_name,json=deviceName" json:"device_name,omitempty"`
}

func (m *DeviceNameSetting) Reset()                    { *m = DeviceNameSetting{} }
func (m *DeviceNameSetting) String() string            { return proto.CompactTextString(m) }
func (*DeviceNameSetting) ProtoMessage()               {}
func (*DeviceNameSetting) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{15} }

func (m *DeviceNameSetting) GetDeviceName() string {
	if m != nil {
		return m.DeviceName
	}
	return ""
}

type HostnameSetting struct {
	Hostname string `protobuf:"bytes,1,opt,name=hostname" json:"hostname,omitempty"`
}

func (m *HostnameSetting) Reset()                    { *m = HostnameSetting{} }
func (m *HostnameSetting) String() string            { return proto.CompactTextString(m) }
func (*HostnameSetting) ProtoMessage()               {}
func (*HostnameSetting) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{16} }

func (m *HostnameSetting) GetHostname() string {
	if m != nil {
		return m.Hostname
	}
	return ""
}

func init() {
	proto.RegisterType((*Event)(nil), "commandproto.Event")
	proto.RegisterType((*Payload)(nil), "commandproto.Payload")
	proto.RegisterType((*Command)(nil), "commandproto.Command")
	proto.RegisterType((*ScheduleOSUpdate)(nil), "commandproto.ScheduleOSUpdate")
	proto.RegisterType((*OSUpdate)(nil), "commandproto.OSUpdate")
	proto.RegisterType((*ScheduleOSUpdateScan)(nil), "commandproto.ScheduleOSUpdateScan")
	proto.RegisterType((*AccountConfiguration)(nil), "commandproto.AccountConfiguration")
	proto.RegisterType((*AutoSetupAdminAccounts)(nil), "commandproto.AutoSetupAdminAccounts")
	proto.RegisterType((*DeviceInformation)(nil), "commandproto.DeviceInformation")
	proto.RegisterType((*InstallProfile)(nil), "commandproto.InstallProfile")
	proto.RegisterType((*RemoveProfile)(nil), "commandproto.RemoveProfile")
	proto.RegisterType((*DeleteUser)(nil), "commandproto.DeleteUser")
	proto.RegisterType((*InstallApplication)(nil), "commandproto.InstallApplication")
	proto.RegisterType((*Settings)(nil), "commandproto.Settings")
	proto.RegisterType((*Setting)(nil), "commandproto.Setting")
	proto.RegisterType((*DeviceNameSetting)(nil), "commandproto.DeviceNameSetting")
	proto.RegisterType((*HostnameSetting)(nil), "commandproto.HostnameSetting")
}

func init() { proto.RegisterFile("command.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 1019 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x56, 0x5d, 0x6f, 0x1c, 0x35,
	0x14, 0xd5, 0xee, 0x26, 0xd9, 0xd9, 0x9b, 0x6c, 0x92, 0x9a, 0x24, 0x1d, 0x28, 0x6d, 0xc2, 0x04,
	0x50, 0x8a, 0x68, 0x0a, 0x05, 0x21, 0xf5, 0x01, 0x89, 0xd0, 0x14, 0x35, 0x82, 0x36, 0xc1, 0x4b,
	0x40, 0x3c, 0x20, 0xcb, 0xcc, 0x78, 0x77, 0xad, 0xce, 0x8c, 0xa7, 0xb6, 0x27, 0x68, 0x1f, 0x78,
	0xe2, 0x17, 0xf0, 0x8a, 0xc4, 0x2f, 0xe2, 0x4f, 0x21, 0x7f, 0x4d, 0x66, 0x3f, 0xfa, 0x36, 0xf7,
	0xf8, 0xf8, 0xf8, 0x8e, 0xef, 0xf1, 0xb5, 0x61, 0x98, 0x8a, 0xa2, 0xa0, 0x65, 0x76, 0x5a, 0x49,
	0xa1, 0x05, 0xda, 0xf2, 0xa1, 0x8d, 0x92, 0x3f, 0x61, 0xfd, 0xf9, 0x0d, 0x2b, 0x35, 0xda, 0x86,
	0x2e, 0xcf, 0xe2, 0xce, 0x51, 0xe7, 0x64, 0x80, 0xbb, 0x3c, 0x43, 0x08, 0xd6, 0x34, 0x2f, 0x58,
	0xdc, 0x3d, 0xea, 0x9c, 0xf4, 0xb0, 0xfd, 0x46, 0x8f, 0xa1, 0x5f, 0xd1, 0x59, 0x2e, 0x68, 0x16,
	0xf7, 0x8e, 0x3a, 0x27, 0x9b, 0x4f, 0xf6, 0x4f, 0xdb, 0x62, 0xa7, 0x57, 0x6e, 0x10, 0x07, 0x16,
	0x3a, 0x84, 0xcd, 0x8c, 0xdd, 0xf0, 0x94, 0x91, 0x3a, 0xe3, 0x59, 0xbc, 0x66, 0xd5, 0xc1, 0x41,
	0xd7, 0x19, 0xcf, 0x92, 0xdf, 0xa0, 0xef, 0x27, 0xa1, 0x0f, 0x20, 0x64, 0x46, 0xea, 0xba, 0x49,
	0x65, 0xd3, 0x63, 0xd7, 0x35, 0xcf, 0xcc, 0xfa, 0x3e, 0xb4, 0x69, 0x2d, 0xad, 0xff, 0xcc, 0x05,
	0x38, 0xb0, 0x92, 0xff, 0xd6, 0xa1, 0xef, 0x41, 0xa3, 0x2f, 0xd9, 0x9b, 0x9a, 0x29, 0x4d, 0xf4,
	0xac, 0x62, 0x41, 0xdf, 0x63, 0x3f, 0xcd, 0x2a, 0x86, 0x5e, 0x01, 0xf2, 0xe9, 0xf2, 0x72, 0x2c,
	0x64, 0x41, 0x35, 0x17, 0xa5, 0x5f, 0xea, 0x70, 0x7e, 0xa9, 0x73, 0xcb, 0xbb, 0xb8, 0xa5, 0xe1,
	0x3b, 0xd9, 0x22, 0x84, 0x9e, 0xc3, 0x0e, 0x2f, 0x95, 0xa6, 0x79, 0x4e, 0x2a, 0x29, 0xc6, 0x3c,
	0x67, 0x7e, 0xdf, 0xde, 0x9f, 0x17, 0xbb, 0x70, 0xa4, 0x2b, 0xc7, 0xc1, 0xdb, 0x7c, 0x2e, 0x46,
	0x3f, 0xc2, 0x3b, 0x41, 0x86, 0x56, 0x55, 0xce, 0x53, 0x97, 0xd7, 0x9a, 0x95, 0x3a, 0x5a, 0x29,
	0x75, 0x76, 0xcb, 0xc3, 0x88, 0x2f, 0x61, 0xe8, 0x17, 0xd8, 0xa7, 0x69, 0x2a, 0xea, 0x52, 0x93,
	0x54, 0x94, 0x63, 0x3e, 0xa9, 0xa5, 0x13, 0x5d, 0xb7, 0xa2, 0xc9, 0xbc, 0xe8, 0x99, 0xa3, 0x3e,
	0x6b, 0x33, 0xf1, 0x1e, 0x5d, 0x81, 0xa2, 0x1f, 0x00, 0xa9, 0x74, 0xca, 0xb2, 0x3a, 0x67, 0x44,
	0x28, 0x52, 0x57, 0x19, 0xd5, 0x2c, 0xde, 0xb0, 0xaa, 0x0f, 0xe6, 0x55, 0x47, 0x9e, 0x77, 0x39,
	0xba, 0xb6, 0x2c, 0xbc, 0x1b, 0x66, 0x5e, 0x2a, 0x87, 0xa0, 0x5f, 0xe1, 0xee, 0xb2, 0x1a, 0x51,
	0x29, 0x2d, 0xe3, 0xfe, 0xaa, 0x44, 0x17, 0x25, 0x47, 0x29, 0x2d, 0xf1, 0xde, 0xa2, 0xac, 0x41,
	0xd1, 0xb7, 0xb0, 0x2d, 0x59, 0x21, 0x6e, 0x58, 0x53, 0x9a, 0xc8, 0x2a, 0xde, 0x9b, 0x57, 0xc4,
	0x96, 0x13, 0x2a, 0x33, 0x94, 0xed, 0x10, 0x3d, 0x35, 0xf6, 0xce, 0x99, 0x66, 0xa4, 0x56, 0x4c,
	0xc6, 0x03, 0x2b, 0x10, 0x2f, 0x1a, 0xc5, 0x10, 0xae, 0x15, 0x93, 0xc6, 0xf8, 0xe1, 0x1b, 0x3d,
	0x81, 0x48, 0x31, 0xad, 0x79, 0x39, 0x51, 0x31, 0xd8, 0x79, 0x07, 0x0b, 0xbf, 0xe2, 0x47, 0x71,
	0xc3, 0x4b, 0xce, 0x61, 0x77, 0xf1, 0x07, 0xd1, 0x67, 0xd0, 0x77, 0xbb, 0xa2, 0xe2, 0xce, 0x51,
	0x6f, 0x59, 0xa6, 0xd9, 0xdc, 0x40, 0x4b, 0x30, 0x44, 0xcd, 0xec, 0x43, 0xd8, 0xac, 0xa4, 0xc8,
	0xea, 0x54, 0x93, 0xd7, 0x6c, 0xe6, 0x8f, 0x04, 0x78, 0xe8, 0x7b, 0x36, 0x43, 0x1f, 0xc1, 0x76,
	0x63, 0xbd, 0xb4, 0x39, 0x0d, 0x03, 0x3c, 0x0c, 0x9e, 0xb2, 0x60, 0xf2, 0x29, 0xec, 0xad, 0xda,
	0x7a, 0xb4, 0x07, 0xeb, 0x63, 0x21, 0x53, 0x77, 0xd8, 0x22, 0xec, 0x82, 0xe4, 0xdf, 0x2e, 0xec,
	0x9d, 0xad, 0x36, 0xcf, 0xb1, 0x7a, 0xcd, 0x2b, 0x52, 0x49, 0x5e, 0x50, 0x39, 0x23, 0x8a, 0xe9,
	0xba, 0x22, 0x8d, 0x51, 0x25, 0x73, 0x1e, 0x75, 0x62, 0x87, 0x86, 0x7a, 0xe5, 0x98, 0x23, 0x43,
	0x0c, 0x92, 0x9e, 0x86, 0x7e, 0x86, 0x87, 0x8a, 0xe9, 0xb7, 0x88, 0x51, 0x45, 0x24, 0x9b, 0xd4,
	0x39, 0x95, 0xae, 0x76, 0x5d, 0xab, 0x79, 0xac, 0x98, 0x5e, 0x21, 0x79, 0xa6, 0xb0, 0xe3, 0xda,
	0xd2, 0x11, 0x78, 0x97, 0xd6, 0x5a, 0x04, 0xc1, 0xac, 0xe0, 0x65, 0x90, 0x55, 0x71, 0xcf, 0x16,
	0xe1, 0xc3, 0x85, 0xf3, 0x53, 0x6b, 0xe1, 0xf4, 0x0c, 0xd9, 0x8b, 0x2a, 0x7c, 0x40, 0x57, 0xe2,
	0xc9, 0xdf, 0x1d, 0x38, 0x58, 0x3d, 0x05, 0xdd, 0x07, 0x50, 0x53, 0x21, 0x35, 0x29, 0x69, 0x11,
	0x5a, 0xd8, 0xc0, 0x22, 0xaf, 0x68, 0xc1, 0xd0, 0x3d, 0x18, 0x8c, 0xeb, 0x3c, 0x77, 0xa3, 0xae,
	0x52, 0x91, 0x01, 0xec, 0xe0, 0x31, 0x0c, 0x2b, 0xaa, 0xd4, 0x1f, 0x42, 0x66, 0x64, 0x4a, 0xd5,
	0xd4, 0xf6, 0xa2, 0x2d, 0xbc, 0x15, 0xc0, 0x17, 0x54, 0x4d, 0xd1, 0x01, 0x6c, 0x4c, 0x79, 0x96,
	0x31, 0xd7, 0x5e, 0x22, 0xec, 0xa3, 0xe4, 0x11, 0xdc, 0x59, 0x6a, 0x79, 0x28, 0x86, 0xfe, 0x9b,
	0x9a, 0x49, 0xee, 0xcd, 0x37, 0xc0, 0x21, 0x4c, 0x3e, 0x81, 0xed, 0xf9, 0xa6, 0x66, 0xb8, 0xe1,
	0xee, 0xe8, 0xd8, 0x75, 0x43, 0x98, 0x3c, 0x86, 0xe1, 0xdc, 0x29, 0x43, 0x0f, 0x00, 0x78, 0xc6,
	0x4a, 0xcd, 0xc7, 0x9c, 0xc9, 0x60, 0xca, 0x5b, 0x24, 0xb9, 0x04, 0xb8, 0x3d, 0x55, 0xe8, 0x3d,
	0x88, 0x4c, 0x05, 0x5b, 0x1b, 0xd2, 0xc4, 0xc6, 0xbe, 0xd6, 0x72, 0xc4, 0x9e, 0xbc, 0x60, 0xdf,
	0x08, 0x0f, 0x2d, 0x7a, 0xee, 0xc1, 0xe4, 0xaf, 0x2e, 0xa0, 0xe5, 0xc6, 0x89, 0x3e, 0x86, 0x1d,
	0xae, 0xeb, 0x92, 0x29, 0xa2, 0xb4, 0x90, 0x8c, 0xf8, 0x4b, 0xa9, 0x87, 0x87, 0x0e, 0x1e, 0x19,
	0xf4, 0x22, 0x5b, 0xc8, 0xb7, 0xbb, 0x98, 0xaf, 0xb9, 0x79, 0x0a, 0x5a, 0xf2, 0xb1, 0xb9, 0x7a,
	0x6a, 0x99, 0xdb, 0x7d, 0x1f, 0xe0, 0xcd, 0x80, 0x5d, 0xcb, 0x1c, 0x3d, 0x84, 0xdd, 0x82, 0x96,
	0x74, 0xc2, 0x0a, 0x56, 0x6a, 0x32, 0xce, 0xe9, 0x44, 0xd9, 0x02, 0xf4, 0xf0, 0xce, 0x2d, 0xfe,
	0x9d, 0x81, 0xcd, 0x99, 0x2d, 0x85, 0x26, 0x0e, 0xce, 0x6c, 0xc3, 0x8e, 0x30, 0x94, 0x42, 0xbf,
	0x74, 0x08, 0xfa, 0x0a, 0xee, 0xa6, 0x53, 0x5a, 0x4e, 0x18, 0x69, 0x49, 0x2a, 0x1d, 0xfa, 0xf0,
	0x00, 0xef, 0xbb, 0xe1, 0x97, 0xcd, 0xe8, 0xc8, 0x0c, 0x26, 0x5f, 0x43, 0x14, 0x9a, 0x0e, 0xfa,
	0xbc, 0xd5, 0x9e, 0x5c, 0x5f, 0xd9, 0x5f, 0xd9, 0x9e, 0x5a, 0xdd, 0xe9, 0x9f, 0x0e, 0xf4, 0x3d,
	0x6a, 0x1e, 0x0f, 0x5c, 0xb3, 0xc2, 0xd7, 0xc3, 0x7e, 0xa3, 0x6f, 0x9a, 0xb7, 0x40, 0xe3, 0xce,
	0xb7, 0xdc, 0xaa, 0xc6, 0xad, 0x41, 0xdf, 0x3f, 0x16, 0xac, 0x81, 0x9f, 0x42, 0x34, 0x15, 0x4a,
	0xdb, 0xe9, 0xee, 0x1e, 0xbd, 0x3f, 0x3f, 0xfd, 0x85, 0x1f, 0x6d, 0x92, 0x0b, 0xf4, 0xe4, 0xcb,
	0x60, 0xdf, 0x96, 0x76, 0xeb, 0x75, 0xd2, 0x32, 0x4f, 0x6b, 0xc1, 0xe4, 0x11, 0xec, 0x2c, 0x48,
	0x1a, 0xb7, 0x35, 0x39, 0x78, 0xb7, 0x85, 0xf8, 0xf7, 0x0d, 0x9b, 0xc5, 0x17, 0xff, 0x07, 0x00,
	0x00, 0xff, 0xff, 0x84, 0x9f, 0x25, 0x5c, 0x71, 0x09, 0x00, 0x00,
}

package mdm

import (
	"fmt"

	"github.com/gogo/protobuf/proto"
	"github.com/micromdm/micromdm/mdm/mdm/internal/mdmproto"
)

func MarshalCommandPayload(cmd *CommandPayload) ([]byte, error) {
	cmdToProto, err := commandToProto(cmd.Command)
	if err != nil {
		return nil, err
	}
	cmdproto := mdmproto.CommandPayload{
		CommandUuid: cmd.CommandUUID,
		Command:     cmdToProto,
	}
	pb, err := proto.Marshal(&cmdproto)
	return pb, err
}

func commandToProto(cmd *Command) (*mdmproto.Command, error) {
	cmdproto := mdmproto.Command{
		RequestType: cmd.RequestType,
	}
	switch cmd.RequestType {
	case "ProfileList",
		"ProvisioningProfileList",
		"CertificateList",
		"SecurityInfo",
		"RestartDevice",
		"ShutDownDevice",
		"StopMirroring",
		"ClearRestrictionsPassword",
		"UserList",
		"LogOutUser",
		"PlayLostModeSound",
		"DisableLostMode",
		"DeviceLocation",
		"ManagedMediaList",
		"DeviceConfigured",
		"AvailableOSUpdates",
		"NSExtensionMappings",
		"OSUpdateStatus",
		"EnableRemoteDesktop",
		"DisableRemoteDesktop",
		"ActivationLockBypassCode":

	case "InstallProfile":
		if cmd.InstallProfile == nil {
			break
		}
		cmdproto.Request = &mdmproto.Command_InstallProfile{
			InstallProfile: &mdmproto.InstallProfile{
				Payload: cmd.InstallProfile.Payload,
			},
		}
	case "RemoveProfile":
		if cmd.RemoveProfile == nil {
			break
		}
		cmdproto.Request = &mdmproto.Command_RemoveProfile{
			RemoveProfile: &mdmproto.RemoveProfile{
				Identifier: cmd.RemoveProfile.Identifier,
			},
		}
	case "InstallProvisioningProfile":
		cmdproto.Request = &mdmproto.Command_InstallProvisioningProfile{
			InstallProvisioningProfile: &mdmproto.InstallProvisioningProfile{
				ProvisioningProfile: cmd.InstallProvisioningProfile.ProvisioningProfile,
			},
		}
	case "RemoveProvisioningProfile":
		cmdproto.Request = &mdmproto.Command_RemoveProfisioningProfile{
			RemoveProfisioningProfile: &mdmproto.RemoveProvisioningProfile{
				Uuid: cmd.RemoveProvisioningProfile.UUID,
			},
		}
	case "InstalledApplicationList":
		cmdproto.Request = &mdmproto.Command_InstalledApplicationList{
			InstalledApplicationList: &mdmproto.InstalledApplicationList{
				Identifiers:     cmd.InstalledApplicationList.Identifiers,
				ManagedAppsOnly: cmd.InstalledApplicationList.ManagedAppsOnly,
			},
		}
	case "DeviceInformation":
		if cmd.DeviceInformation == nil {
			break
		}
		cmdproto.Request = &mdmproto.Command_DeviceInformation{
			DeviceInformation: &mdmproto.DeviceInformation{
				Queries: cmd.DeviceInformation.Queries,
			},
		}
	case "DeviceLock":
		cmdproto.Request = &mdmproto.Command_DeviceLock{
			DeviceLock: &mdmproto.DeviceLock{
				Pin:         cmd.DeviceLock.PIN,
				Message:     cmd.DeviceLock.Message,
				PhoneNumber: cmd.DeviceLock.PhoneNumber,
			},
		}
	case "ClearPasscode":
		cmdproto.Request = &mdmproto.Command_ClearPasscode{
			ClearPasscode: &mdmproto.ClearPasscode{
				UnlockToken: cmd.ClearPasscode.UnlockToken,
			},
		}
	case "EraseDevice":
		cmdproto.Request = &mdmproto.Command_EraseDevice{
			EraseDevice: &mdmproto.EraseDevice{
				Pin:                    cmd.EraseDevice.PIN,
				PreserveDataPlan:       cmd.EraseDevice.PreserveDataPlan,
				DisallowProximitySetup: cmd.EraseDevice.DisallowProximitySetup,
			},
		}
	case "RequestMirroring":
		cmdproto.Request = &mdmproto.Command_RequestMirroring{
			RequestMirroring: &mdmproto.RequestMirroring{
				DestinationName:     cmd.RequestMirroring.DestinationName,
				DestinationDeviceId: cmd.RequestMirroring.DestinationDeviceID,
				ScanTime:            cmd.RequestMirroring.ScanTime,
				Password:            cmd.RequestMirroring.Password,
			},
		}
	case "Restrictions":
		cmdproto.Request = &mdmproto.Command_Restrictions{
			Restrictions: &mdmproto.Restrictions{
				ProfileRestrictions: cmd.Restrictions.ProfileRestrictions,
			},
		}
	case "UnlockUserAccount":
		cmdproto.Request = &mdmproto.Command_UnlockUserAccount{
			UnlockUserAccount: &mdmproto.UnlockUserAccount{
				Username: cmd.UnlockUserAccount.UserName,
			},
		}
	case "DeleteUser":
		cmdproto.Request = &mdmproto.Command_DeleteUser{
			DeleteUser: &mdmproto.DeleteUser{
				Username:      cmd.DeleteUser.UserName,
				ForceDeletion: cmd.DeleteUser.ForceDeletion,
			},
		}
	case "EnableLostMode":
		cmdproto.Request = &mdmproto.Command_EnableLostMode{
			EnableLostMode: &mdmproto.EnableLostMode{
				Message:     cmd.EnableLostMode.Message,
				PhoneNumber: cmd.EnableLostMode.PhoneNumber,
				Footnote:    cmd.EnableLostMode.Footnote,
			},
		}
	case "InstallEnterpriseApplication":
		var pbManifest *mdmproto.Manifest
		if cmd.InstallEnterpriseApplication.Manifest != nil {
			pbManifest = &(mdmproto.Manifest{})
			for _, item := range cmd.InstallEnterpriseApplication.Manifest.ManifestItems {
				var pbManifestItem mdmproto.ManifestItem
				for _, asset := range item.Assets {
					pbAsset := mdmproto.Asset{
						Kind:       asset.Kind,
						Md5Size:    asset.MD5Size,
						Md5S:       asset.MD5s,
						Sha256Size: asset.SHA256Size,
						Sha256S:    asset.SHA256s,
						Url:        asset.URL,
					}
					pbManifestItem.Assets = append(pbManifestItem.Assets, &pbAsset)
				}
				if item.Metadata != nil {
					pbManifestItem.Metadata = &(mdmproto.Metadata{
						BundleIdentifier: item.Metadata.BundleInfo.BundleIdentifier,
						BundleVersion:    item.Metadata.BundleInfo.BundleVersion,
						Kind:             item.Metadata.Kind,
						SizeInBytes:      item.Metadata.SizeInBytes,
						Title:            item.Metadata.Title,
						Subtitle:         item.Metadata.Subtitle,
					})
					for _, bundleInfo := range item.Metadata.Items {
						bpBundleInfo := &(mdmproto.BundleInfo{
							BundleIdentifier: bundleInfo.BundleIdentifier,
							BundleVersion:    bundleInfo.BundleVersion,
						})
						pbManifestItem.Metadata.Items = append(pbManifestItem.Metadata.Items, bpBundleInfo)
					}
				}
				if len(pbManifestItem.Assets) > 0 || pbManifestItem.Metadata != nil {
					pbManifest.ManifestItems = append(pbManifest.ManifestItems, &pbManifestItem)
				}
			}
		}
		cmdproto.Request = &mdmproto.Command_InstallEnterpriseApplication{
			InstallEnterpriseApplication: &mdmproto.InstallEnterpriseApplication{
				Manifest:                       pbManifest,
				ManifestUrl:                    emptyStringIfNil(cmd.InstallEnterpriseApplication.ManifestURL),
				ManifestUrlPinningCerts:        cmd.InstallEnterpriseApplication.ManifestURLPinningCerts,
				PinningRevocationCheckRequired: falseIfNil(cmd.InstallEnterpriseApplication.PinningRevocationCheckRequired),
			},
		}
	case "InstallApplication":
		var (
			options       *mdmproto.InstallApplicationOptions
			configuration *mdmproto.InstallApplicationConfiguration
			attributes    *mdmproto.InstallApplicationAttributes
		)
		if cmd.InstallApplication.Options != nil {
			options = &mdmproto.InstallApplicationOptions{
				PurchaseMethod: zeroInt64IfNil(cmd.InstallApplication.Options.PurchaseMethod),
			}
		}
		if cmd.InstallApplication.Configuration != nil {
			configuration = &mdmproto.InstallApplicationConfiguration{}
		}
		if cmd.InstallApplication.Attributes != nil {
			attributes = &mdmproto.InstallApplicationAttributes{}
		}
		cmdproto.Request = &mdmproto.Command_InstallApplication{
			InstallApplication: &mdmproto.InstallApplication{
				ItunesStoreId:         zeroInt64IfNil(cmd.InstallApplication.ITunesStoreID),
				Identifier:            emptyStringIfNil(cmd.InstallApplication.Identifier),
				ManagementFlags:       int64(zeroIntIfNil(cmd.InstallApplication.ManagementFlags)),
				ChangeManagementState: emptyStringIfNil(cmd.InstallApplication.ChangeManagementState),
				ManifestUrl:           emptyStringIfNil(cmd.InstallApplication.ManifestURL),
				Options:               options,
				Configuration:         configuration,
				Attributes:            attributes,
			},
		}
	case "AccountConfiguration":
		var autosetupadminaccounts []*mdmproto.AutoSetupAdminAccounts
		for _, acct := range cmd.AccountConfiguration.AutoSetupAdminAccounts {
			autosetupadminaccounts = append(autosetupadminaccounts, &mdmproto.AutoSetupAdminAccounts{
				ShortName:    acct.ShortName,
				FullName:     acct.FullName,
				PasswordHash: acct.PasswordHash,
				Hidden:       acct.Hidden,
			})
		}
		cmdproto.Request = &mdmproto.Command_AccountConfiguration{
			AccountConfiguration: &mdmproto.AccountConfiguration{
				SkipPrimarySetupAccountCreation:     cmd.AccountConfiguration.SkipPrimarySetupAccountCreation,
				SetPrimarySetupAccountAsRegularUser: cmd.AccountConfiguration.SetPrimarySetupAccountAsRegularUser,
				DontAutoPopulatePrimaryAccountInfo:  cmd.AccountConfiguration.DontAutoPopulatePrimaryAccountInfo,
				LockPrimaryAccountInfo:              cmd.AccountConfiguration.LockPrimaryAccountInfo,
				PrimaryAccountFullName:              cmd.AccountConfiguration.PrimaryAccountFullName,
				PrimaryAccountUserName:              cmd.AccountConfiguration.PrimaryAccountUserName,
				ManagedLocalUserShortName:           cmd.AccountConfiguration.ManagedLocalUserShortName,
				AutoSetupAdminAccounts:              autosetupadminaccounts,
			},
		}
	case "ApplyRedemptionCode":
		cmdproto.Request = &mdmproto.Command_ApplyRedemptionCode{
			ApplyRedemptionCode: &mdmproto.ApplyRedemptionCode{
				Identifier:     cmd.ApplyRedemptionCode.Identifier,
				RedemptionCode: cmd.ApplyRedemptionCode.RedemptionCode,
			},
		}
	case "ManagedApplicationList":
		cmdproto.Request = &mdmproto.Command_ManagedApplicationList{
			ManagedApplicationList: &mdmproto.ManagedApplicationList{
				Identifiers: cmd.ManagedApplicationList.Identifiers,
			},
		}
	case "RemoveApplication":
		cmdproto.Request = &mdmproto.Command_RemoveApplication{
			RemoveApplication: &mdmproto.RemoveApplication{
				Identifier: cmd.RemoveApplication.Identifier,
			},
		}
	case "InviteToProgram":
		cmdproto.Request = &mdmproto.Command_InviteToProgram{
			InviteToProgram: &mdmproto.InviteToProgram{
				ProgramId:     cmd.InviteToProgram.ProgramID,
				InvitationUrl: cmd.InviteToProgram.InvitationURL,
			},
		}
	case "ValidateApplications":
		cmdproto.Request = &mdmproto.Command_ValidataApplications{
			ValidataApplications: &mdmproto.ValidateApplications{
				Identifiers: cmd.ValidateApplications.Identifiers,
			},
		}
	case "InstallMedia":
		cmdproto.Request = &mdmproto.Command_InstallMedia{
			InstallMedia: &mdmproto.InstallMedia{
				ItunesStoreId: zeroInt64IfNil(cmd.InstallMedia.ITunesStoreID),
				MediaType:     cmd.InstallMedia.MediaType,
				MediaUrl:      cmd.InstallMedia.MediaURL,
			},
		}
	case "RemoveMedia":
		cmdproto.Request = &mdmproto.Command_RemoveMedia{
			RemoveMedia: &mdmproto.RemoveMedia{
				ItunesStoreId: zeroInt64IfNil(cmd.RemoveMedia.ITunesStoreID),
				MediaType:     cmd.RemoveMedia.MediaType,
				PersistentId:  cmd.RemoveMedia.PersistentID,
			},
		}
	case "Settings":
		var settings []*mdmproto.Setting
		for _, s := range cmd.Settings.Settings {
			settings = append(settings, settingToProto(s))
		}
		cmdproto.Request = &mdmproto.Command_Settings{
			Settings: &mdmproto.Settings{Settings: settings},
		}
	case "ManagedApplicationConfiguration":
		cmdproto.Request = &mdmproto.Command_ManagedApplicationConfiguration{
			ManagedApplicationConfiguration: &mdmproto.ManagedApplicationConfiguration{
				Identifiers: cmd.ManagedApplicationConfiguration.Identifiers,
			},
		}
	case "ManagedApplicationAttributes":
		cmdproto.Request = &mdmproto.Command_ManagedApplicationAttributes{
			ManagedApplicationAttributes: &mdmproto.ManagedApplicationAttributes{
				Identifiers: cmd.ManagedApplicationAttributes.Identifiers,
			},
		}
	case "ManagedApplicationFeedback":
		cmdproto.Request = &mdmproto.Command_ManagedApplicationFeedback{
			ManagedApplicationFeedback: &mdmproto.ManagedApplicationFeedback{
				Identifiers:    cmd.ManagedApplicationFeedback.Identifiers,
				DeleteFeedback: cmd.ManagedApplicationFeedback.DeleteFeedback,
			},
		}
	case "SetFirmwarePassword":
		cmdproto.Request = &mdmproto.Command_SetFirmwarePassword{
			SetFirmwarePassword: &mdmproto.SetFirmwarePassword{
				CurrentPassword: cmd.SetFirmwarePassword.CurrentPassword,
				NewPassword:     cmd.SetFirmwarePassword.NewPassword,
				AllowOroms:      cmd.SetFirmwarePassword.AllowOroms,
			},
		}
	case "SetBootstrapToken":
		cmdproto.Request = &mdmproto.Command_SetBootstrapToken{
			SetBootstrapToken: &mdmproto.SetBootstrapToken{
				BootstrapToken: cmd.SetBootstrapToken.BootstrapToken,
			},
		}
	case "VerifyFirmwarePassword":
		cmdproto.Request = &mdmproto.Command_VerifyFirmwarePassword{
			VerifyFirmwarePassword: &mdmproto.VerifyFirmwarePassword{
				Password: cmd.VerifyFirmwarePassword.Password,
			},
		}
	case "SetAutoAdminPassword":
		cmdproto.Request = &mdmproto.Command_SetAutoAdminPassword{
			SetAutoAdminPassword: &mdmproto.SetAutoAdminPassword{
				Guid:         cmd.SetAutoAdminPassword.GUID,
				PasswordHash: cmd.SetAutoAdminPassword.PasswordHash,
			},
		}
	case "ScheduleOSUpdate":
		var updates []*mdmproto.Update
		for _, u := range cmd.ScheduleOSUpdate.Updates {
			updates = append(updates, &mdmproto.Update{
				ProductKey:    u.ProductKey,
				InstallAction: u.InstallAction,
			})
		}
		cmdproto.Request = &mdmproto.Command_ScheduleOsUpdate{
			ScheduleOsUpdate: &mdmproto.ScheduleOSUpdate{
				Updates: updates,
			},
		}
	case "ScheduleOSUpdateScan":
		cmdproto.Request = &mdmproto.Command_ScheduleOsUpdateScan{
			ScheduleOsUpdateScan: &mdmproto.ScheduleOSUpdateScan{
				Force: cmd.ScheduleOSUpdateScan.Force,
			},
		}
	case "ActiveNSExtensions":
		cmdproto.Request = &mdmproto.Command_ActiveNsExtensions{
			ActiveNsExtensions: &mdmproto.ActiveNSExtensions{
				FilterExtensionPoints: cmd.ActiveNSExtensions.FilterExtensionPoints,
			},
		}
	case "RotateFileVaultKey":
		fvunlock := cmd.RotateFileVaultKey.FileVaultUnlock
		cmdproto.Request = &mdmproto.Command_RotateFilevaultKey{
			RotateFilevaultKey: &mdmproto.RotateFileVaultKey{
				KeyType:                    cmd.RotateFileVaultKey.KeyType,
				NewCertificate:             cmd.RotateFileVaultKey.NewCertificate,
				ReplyEncryptionCertificate: cmd.RotateFileVaultKey.ReplyEncryptionCertificate,
				FilevaultUnlock: &mdmproto.FileVaultUnlock{
					Password:                 fvunlock.Password,
					PrivateKeyExport:         fvunlock.PrivateKeyExport,
					PrivateKeyExportPassword: fvunlock.PrivateKeyExportPassword,
				},
			},
		}
	default:
		return nil, fmt.Errorf("unknown command type %s", cmd.RequestType)
	}
	return &cmdproto, nil
}

func settingToProto(s Setting) *mdmproto.Setting {
	pbs := mdmproto.Setting{Item: s.Item}
	switch s.Item {
	case "ApplicationConfiguration":
		pbs.ApplicationConfiguration = &mdmproto.ApplicationConfigurationSetting{
			Identifier:                  emptyStringIfNil(s.Identifier),
			ConfigurationDictionaryData: s.ConfigurationData,
		}
	case "VoiceRoaming":
		pbs.VoiceRoaming = &mdmproto.VoiceRoamingSetting{
			Enabled: falseIfNil(s.Enabled),
		}
	case "PersonalHotspot":
		pbs.PersonalHotspot = &mdmproto.PersonalHotspotSetting{
			Enabled: falseIfNil(s.Enabled),
		}
	case "Wallpaper":
		pbs.Wallpaper = &mdmproto.WallpaperSetting{
			Image: s.Image,
			Where: int64(zeroIntIfNil(s.Where)),
		}
	case "DataRoaming":
		pbs.DataRoaming = &mdmproto.DataRoamingSetting{
			Enabled: falseIfNil(s.Enabled),
		}
	case "Bluetooth":
		pbs.Bluetooth = &mdmproto.BluetoothSetting{
			Enabled: falseIfNil(s.Enabled),
		}
	case "ApplicationAttributes":
		var attributes *mdmproto.ApplicationAttributes
		if v, ok := s.Attributes["VPNUUID"]; ok {
			attributes.VpnUuid = v
		}
		pbs.ApplicationAttributes = &mdmproto.ApplicationAttributesSetting{
			Identifier:            emptyStringIfNil(s.Identifier),
			ApplicationAttributes: attributes,
		}
	case "DeviceName":
		pbs.DeviceName = &mdmproto.DeviceNameSetting{
			DeviceName: emptyStringIfNil(s.DeviceName),
		}
	case "HostName":
		pbs.Hostname = &mdmproto.HostnameSetting{
			Hostname: emptyStringIfNil(s.HostName),
		}
	case "MDMOptions":
		options := &mdmproto.MDMOptions{}
		if v, ok := s.MDMOptions["ActivationLockAllowedWhileSupervised"]; ok {
			options.ActivationLockAllowedWhileSupervised = v.(bool)
		}
		pbs.MdmOptions = &mdmproto.MDMOptionsSetting{
			MdmOptions: options,
		}
	case "PasscodeLockGracePeriod":
		pbs.PasscodeLockGracePeriod = &mdmproto.PasscodeLockGracePeriodSetting{
			PasscodeLockGracePeriod: int64(zeroIntIfNil(s.PasscodeLockGracePeriod)),
		}
	case "MaximumResidentUsers":
		pbs.MaximumResidentUsers = &mdmproto.MaximumResidentUsersSetting{
			MaximumResidentUsers: int64(zeroIntIfNil(s.MaximumResidentUsers)),
		}
	case "DiagnosticSubmission":
		pbs.DiagnosticSubmission = &mdmproto.DiagnosticSubmissionSetting{
			Enabled: falseIfNil(s.Enabled),
		}
	case "AppAnalytics":
		pbs.AppAnalytics = &mdmproto.AppAnalyticsSetting{
			Enabled: falseIfNil(s.Enabled),
		}
	}
	return &pbs
}

func falseIfNil(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

func zeroInt64IfNil(n *int64) int64 {
	if n == nil {
		return 0
	}
	return *n
}

func zeroIntIfNil(n *int) int {
	if n == nil {
		return 0
	}
	return *n
}

func emptyStringIfNil(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

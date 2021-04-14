package mdm

import (
	"github.com/gogo/protobuf/proto"
	"github.com/micromdm/micromdm/mdm/appmanifest"
	"github.com/micromdm/micromdm/mdm/mdm/internal/mdmproto"
)

func protoToCommand(pb *mdmproto.Command) *Command {
	cmd := Command{
		RequestType: pb.RequestType,
	}
	switch pb.RequestType {
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
		cmd.InstallProfile = &InstallProfile{
			Payload: pb.GetInstallProfile().GetPayload(),
		}
	case "RemoveProfile":
		cmd.RemoveProfile = &RemoveProfile{
			Identifier: pb.GetRemoveProfile().GetIdentifier(),
		}
	case "InstallProvisioningProfile":
		cmd.InstallProvisioningProfile = &InstallProvisioningProfile{
			ProvisioningProfile: pb.GetInstallProvisioningProfile().GetProvisioningProfile(),
		}
	case "RemoveProvisioningProfile":
		cmd.RemoveProvisioningProfile = &RemoveProvisioningProfile{
			UUID: pb.GetRemoveProfisioningProfile().GetUuid(),
		}
	case "InstalledApplicationList":
		pbcmd := pb.GetInstalledApplicationList()
		cmd.InstalledApplicationList = &InstalledApplicationList{
			Identifiers:     pbcmd.GetIdentifiers(),
			ManagedAppsOnly: pbcmd.GetManagedAppsOnly(),
		}
	case "DeviceInformation":
		cmd.DeviceInformation = &DeviceInformation{
			Queries: pb.GetDeviceInformation().GetQueries(),
		}
	case "DeviceLock":
		pbc := pb.GetDeviceLock()
		cmd.DeviceLock = &DeviceLock{
			PIN:         pbc.GetPin(),
			Message:     pbc.GetMessage(),
			PhoneNumber: pbc.GetPhoneNumber(),
		}
	case "ClearPasscode":
		pbc := pb.GetClearPasscode()
		cmd.ClearPasscode = &ClearPasscode{
			UnlockToken: pbc.GetUnlockToken(),
		}
	case "EraseDevice":
		pbc := pb.GetEraseDevice()
		cmd.EraseDevice = &EraseDevice{
			PIN:                    pbc.GetPin(),
			PreserveDataPlan:       pbc.GetPreserveDataPlan(),
			DisallowProximitySetup: pbc.GetDisallowProximitySetup(),
		}
	case "RequestMirroring":
		pbc := pb.GetRequestMirroring()
		cmd.RequestMirroring = &RequestMirroring{
			DestinationName:     pbc.GetDestinationName(),
			DestinationDeviceID: pbc.GetDestinationDeviceId(),
			ScanTime:            pbc.GetScanTime(),
			Password:            pbc.GetPassword(),
		}
	case "Restrictions":
		pbc := pb.GetRestrictions()
		cmd.Restrictions = &Restrictions{
			ProfileRestrictions: pbc.GetProfileRestrictions(),
		}
	case "UnlockUserAccount":
		pbc := pb.GetUnlockUserAccount()
		cmd.UnlockUserAccount = &UnlockUserAccount{
			UserName: pbc.GetUsername(),
		}
	case "DeleteUser":
		pbc := pb.GetDeleteUser()
		cmd.DeleteUser = &DeleteUser{
			UserName:      pbc.GetUsername(),
			ForceDeletion: pbc.GetForceDeletion(),
		}
	case "EnableLostMode":
		pbc := pb.GetEnableLostMode()
		cmd.EnableLostMode = &EnableLostMode{
			Message:     pbc.GetMessage(),
			PhoneNumber: pbc.GetPhoneNumber(),
			Footnote:    pbc.GetFootnote(),
		}
	case "InstallEnterpriseApplication":
		pbc := pb.GetInstallEnterpriseApplication()
		var manifest *appmanifest.Manifest
		if pbManifest := pbc.GetManifest(); pbManifest != nil {
			manifest = &(appmanifest.Manifest{})
			if pbManifestItems := pbManifest.GetManifestItems(); pbManifestItems != nil {
				for _, pbManifestItem := range pbManifestItems {
					var manifestItem appmanifest.Item
					if pbAssets := pbManifestItem.GetAssets(); pbAssets != nil {
						for _, pbAsset := range pbAssets {
							asset := appmanifest.Asset{
								Kind:       pbAsset.Kind,
								MD5s:       pbAsset.Md5S,
								MD5Size:    pbAsset.Md5Size,
								URL:        pbAsset.Url,
								SHA256s:    pbAsset.Sha256S,
								SHA256Size: pbAsset.Sha256Size,
							}
							manifestItem.Assets = append(manifestItem.Assets, asset)
						}
					}
					if pbMetadata := pbManifestItem.GetMetadata(); pbMetadata != nil {
						manifestItem.Metadata = &(appmanifest.Metadata{
							Kind:        pbMetadata.Kind,
							SizeInBytes: pbMetadata.SizeInBytes,
							Title:       pbMetadata.Title,
							Subtitle:    pbMetadata.Subtitle,
							BundleInfo: appmanifest.BundleInfo{
								BundleIdentifier: pbMetadata.BundleIdentifier,
								BundleVersion:    pbMetadata.BundleVersion,
							},
						})
						if pbMetadataItems := pbMetadata.GetItems(); pbMetadataItems != nil {
							for _, pbMetadataItem := range pbMetadataItems {
								bundleInfo := appmanifest.BundleInfo{
									BundleIdentifier: pbMetadataItem.BundleIdentifier,
									BundleVersion:    pbMetadataItem.BundleVersion,
								}
								manifestItem.Metadata.Items = append(manifestItem.Metadata.Items, bundleInfo)
							}
						}
					}
					if len(manifestItem.Assets) > 0 || manifestItem.Metadata != nil {
						manifest.ManifestItems = append(manifest.ManifestItems, manifestItem)
					}
				}
			}
		}

		cmd.InstallEnterpriseApplication = &InstallEnterpriseApplication{
			Manifest:                       manifest,
			ManifestURL:                    nilIfEmptyString(pbc.GetManifestUrl()),
			ManifestURLPinningCerts:        pbc.GetManifestUrlPinningCerts(),
			PinningRevocationCheckRequired: nilIfFalse(pbc.GetPinningRevocationCheckRequired()),
		}
	case "InstallApplication":
		pbc := pb.GetInstallApplication()
		var mgmt *int
		mgmtFlags := nilIfZeroInt64(pbc.GetManagementFlags())
		if mgmtFlags != nil {
			mgmti := int(*mgmtFlags)
			mgmt = &mgmti
		}

		var (
			options       *InstallApplicationOptions
			configuration *InstallApplicationConfiguration
			attributes    *InstallApplicationAttributes
		)

		pboptions := pbc.GetOptions()
		if pboptions != nil {
			options = &InstallApplicationOptions{
				PurchaseMethod: new(int64),
			}
			*options.PurchaseMethod = pboptions.GetPurchaseMethod();
		}

		pbconfig := pbc.GetConfiguration()
		if pbconfig != nil {
			configuration = &InstallApplicationConfiguration{}
		}

		pbattributes := pbc.GetAttributes()
		if pbattributes != nil {
			attributes = &InstallApplicationAttributes{}
		}

		cmd.InstallApplication = &InstallApplication{
			ITunesStoreID:         nilIfZeroInt64(pbc.GetItunesStoreId()),
			Identifier:            nilIfEmptyString(pbc.GetIdentifier()),
			ManagementFlags:       mgmt,
			ChangeManagementState: nilIfEmptyString(pbc.GetChangeManagementState()),
			ManifestURL:           nilIfEmptyString(pbc.GetManifestUrl()),
			Options:               options,
			Configuration:         configuration,
			Attributes:            attributes,
		}
	case "AccountConfiguration":
		pbc := pb.GetAccountConfiguration()
		var autosetupadminaccounts []AdminAccount
		for _, acct := range pbc.GetAutoSetupAdminAccounts() {
			autosetupadminaccounts = append(autosetupadminaccounts, AdminAccount{
				ShortName:    acct.GetShortName(),
				FullName:     acct.GetFullName(),
				PasswordHash: acct.GetPasswordHash(),
				Hidden:       acct.GetHidden(),
			})
		}
		cmd.AccountConfiguration = &AccountConfiguration{
			SkipPrimarySetupAccountCreation:     pbc.GetSkipPrimarySetupAccountCreation(),
			SetPrimarySetupAccountAsRegularUser: pbc.GetSetPrimarySetupAccountAsRegularUser(),
			DontAutoPopulatePrimaryAccountInfo:  pbc.GetDontAutoPopulatePrimaryAccountInfo(),
			LockPrimaryAccountInfo:              pbc.GetLockPrimaryAccountInfo(),
			PrimaryAccountFullName:              pbc.GetPrimaryAccountFullName(),
			PrimaryAccountUserName:              pbc.GetPrimaryAccountUserName(),
			AutoSetupAdminAccounts:              autosetupadminaccounts,
		}
	case "ApplyRedemptionCode":
		pbc := pb.GetApplyRedemptionCode()
		cmd.ApplyRedemptionCode = &ApplyRedemptionCode{
			Identifier:     pbc.GetIdentifier(),
			RedemptionCode: pbc.GetRedemptionCode(),
		}
	case "ManagedApplicationList":
		pbc := pb.GetManagedApplicationList()
		cmd.ManagedApplicationList = &ManagedApplicationList{
			Identifiers: pbc.GetIdentifiers(),
		}
	case "RemoveApplication":
		pbc := pb.GetRemoveApplication()
		cmd.RemoveApplication = &RemoveApplication{
			Identifier: pbc.GetIdentifier(),
		}
	case "InviteToProgram":
		pbc := pb.GetInviteToProgram()
		cmd.InviteToProgram = &InviteToProgram{
			ProgramID:     pbc.GetProgramId(),
			InvitationURL: pbc.GetInvitationUrl(),
		}
	case "ValidateApplications":
		pbc := pb.GetValidataApplications()
		cmd.ValidateApplications = &ValidateApplications{
			Identifiers: pbc.GetIdentifiers(),
		}
	case "InstallMedia":
		pbc := pb.GetInstallMedia()
		cmd.InstallMedia = &InstallMedia{
			ITunesStoreID: nilIfZeroInt64(pbc.GetItunesStoreId()),
			MediaType:     pbc.GetMediaType(),
			MediaURL:      pbc.GetMediaUrl(),
		}
	case "RemoveMedia":
		pbc := pb.GetRemoveMedia()
		cmd.RemoveMedia = &RemoveMedia{
			ITunesStoreID: nilIfZeroInt64(pbc.GetItunesStoreId()),
			MediaType:     pbc.GetMediaType(),
			PersistentID:  pbc.GetPersistentId(),
		}
	case "Settings":
		pbc := pb.GetSettings()
		var settings []Setting
		for _, s := range pbc.GetSettings() {
			settings = append(settings, protoToSetting(s))
		}
		cmd.Settings = &Settings{
			Settings: settings,
		}
	case "ManagedApplicationConfiguration":
		pbc := pb.GetManagedApplicationConfiguration()
		cmd.ManagedApplicationConfiguration = &ManagedApplicationConfiguration{
			Identifiers: pbc.GetIdentifiers(),
		}
	case "ManagedApplicationAttributes":
		pbc := pb.GetManagedApplicationAttributes()
		cmd.ManagedApplicationAttributes = &ManagedApplicationAttributes{
			Identifiers: pbc.GetIdentifiers(),
		}
	case "ManagedApplicationFeedback":
		pbc := pb.GetManagedApplicationFeedback()
		cmd.ManagedApplicationFeedback = &ManagedApplicationFeedback{
			Identifiers:    pbc.GetIdentifiers(),
			DeleteFeedback: pbc.GetDeleteFeedback(),
		}
	case "SetFirmwarePassword":
		pbc := pb.GetSetFirmwarePassword()
		cmd.SetFirmwarePassword = &SetFirmwarePassword{
			CurrentPassword: pbc.GetCurrentPassword(),
			NewPassword:     pbc.GetNewPassword(),
			AllowOroms:      pbc.GetAllowOroms(),
		}
	case "VerifyFirmwarePassword":
		pbc := pb.GetVerifyFirmwarePassword()
		cmd.VerifyFirmwarePassword = &VerifyFirmwarePassword{
			Password: pbc.GetPassword(),
		}
	case "SetAutoAdminPassword":
		pbc := pb.GetSetAutoAdminPassword()
		cmd.SetAutoAdminPassword = &SetAutoAdminPassword{
			GUID:         pbc.GetGuid(),
			PasswordHash: pbc.GetPasswordHash(),
		}
	case "ScheduleOSUpdate":
		pbc := pb.GetScheduleOsUpdate()
		var updates []OSUpdate
		for _, up := range pbc.GetUpdates() {
			updates = append(updates, OSUpdate{
				ProductKey:    up.GetProductKey(),
				InstallAction: up.GetInstallAction(),
			})
		}
		cmd.ScheduleOSUpdate = &ScheduleOSUpdate{
			Updates: updates,
		}
	case "ScheduleOSUpdateScan":
		pbc := pb.GetScheduleOsUpdateScan()
		cmd.ScheduleOSUpdateScan = &ScheduleOSUpdateScan{
			Force: pbc.GetForce(),
		}
	case "ActiveNSExtensions":
		pbc := pb.GetActiveNsExtensions()
		cmd.ActiveNSExtensions = &ActiveNSExtensions{
			FilterExtensionPoints: pbc.GetFilterExtensionPoints(),
		}
	case "RotateFileVaultKey":
		pbc := pb.GetRotateFilevaultKey()
		fvunlock := pbc.GetFilevaultUnlock()
		cmd.RotateFileVaultKey = &RotateFileVaultKey{
			KeyType:                    pbc.GetKeyType(),
			NewCertificate:             pbc.GetNewCertificate(),
			ReplyEncryptionCertificate: pbc.GetReplyEncryptionCertificate(),
			FileVaultUnlock: FileVaultUnlock{
				Password:                 fvunlock.GetPassword(),
				PrivateKeyExport:         fvunlock.GetPrivateKeyExport(),
				PrivateKeyExportPassword: fvunlock.GetPrivateKeyExportPassword(),
			},
		}

	}
	return &cmd
}

func protoToSetting(s *mdmproto.Setting) Setting {
	setting := Setting{
		Item: s.GetItem(),
	}
	switch s.Item {
	case "ApplicationConfiguration":
		pbs := s.GetApplicationConfiguration()
		setting.Identifier = nilIfEmptyString(pbs.GetIdentifier())
		setting.ConfigurationData = pbs.GetConfigurationDictionaryData()
	case "VoiceRoaming":
		pbs := s.GetVoiceRoaming()
		setting.Enabled = nilIfFalse(pbs.GetEnabled())
	case "PersonalHotspot":
		pbs := s.GetPersonalHotspot()
		setting.Enabled = nilIfFalse(pbs.GetEnabled())
	case "Wallpaper":
		pbs := s.GetWallpaper()
		setting.Image = pbs.GetImage()
		setting.Where = nilIfZeroInt(int(pbs.GetWhere()))
	case "DataRoaming":
		pbs := s.GetDataRoaming()
		setting.Enabled = nilIfFalse(pbs.GetEnabled())
	case "Bluetooth":
		pbs := s.GetBluetooth()
		setting.Enabled = nilIfFalse(pbs.GetEnabled())
	case "ApplicationAttributes":
		pbs := s.GetApplicationAttributes()
		attr := pbs.GetApplicationAttributes()
		vpnUUID := attr.GetVpnUuid()
		if vpnUUID != "" {
			setting.Attributes = map[string]string{"VPNUUID": vpnUUID}
		}
		setting.Identifier = nilIfEmptyString(pbs.GetIdentifier())
	case "DeviceName":
		pbs := s.GetDeviceName()
		setting.DeviceName = nilIfEmptyString(pbs.GetDeviceName())
	case "TimeZone":
		pbs := s.GetTimeZone()
		setting.TimeZone = nilIfEmptyString(pbs.GetTimeZone())
	case "HostName":
		pbs := s.GetHostname()
		setting.HostName = nilIfEmptyString(pbs.GetHostname())
	case "MDMOptions":
		pbs := s.GetMdmOptions()
		activationLockAllowed := pbs.GetMdmOptions().GetActivationLockAllowedWhileSupervised()
		bootstrapTokenAllowed := pbs.GetMdmOptions().GetBootstrapTokenAllowed()
		promptUserToAllowBootstrapToken := pbs.GetMdmOptions().GetPromptUserToAllowBootstrapTokenForAuthentication()
		setting.MDMOptions = map[string]interface {
		}{
			"ActivationLockAllowedWhileSupervised": activationLockAllowed,
			"BootstrapTokenAllowed": bootstrapTokenAllowed,
			"PromptUserToAllowBootstrapTokenForAuthentication": promptUserToAllowBootstrapToken,
		}
	case "PasscodeLockGracePeriod":
		pbs := s.GetPasscodeLockGracePeriod()
		setting.PasscodeLockGracePeriod = nilIfZeroInt(int(pbs.GetPasscodeLockGracePeriod()))
	case "MaximumResidentUsers":
		pbs := s.GetMaximumResidentUsers()
		setting.MaximumResidentUsers = nilIfZeroInt(int(pbs.GetMaximumResidentUsers()))
	case "DiagnosticSubmission":
		pbs := s.GetDiagnosticSubmission()
		setting.Enabled = nilIfFalse(pbs.GetEnabled())
	case "AppAnalytics":
		pbs := s.GetAppAnalytics()
		setting.Enabled = nilIfFalse(pbs.GetEnabled())
	}
	return setting
}

func UnmarshalCommandPayload(data []byte, payload *CommandPayload) error {
	var pb mdmproto.CommandPayload
	if err := proto.Unmarshal(data, &pb); err != nil {
		return err
	}
	payload.CommandUUID = pb.CommandUuid
	payload.Command = protoToCommand(pb.Command)
	return nil
}

func nilIfZeroInt64(n int64) *int64 {
	if n == 0 {
		return nil
	}
	return &n
}

func nilIfZeroInt(n int) *int {
	if n == 0 {
		return nil
	}
	return &n
}

func nilIfEmptyString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func nilIfFalse(b bool) *bool {
	if !b {
		return nil

	}
	return &b
}

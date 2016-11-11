package mdm

type MCProfileErrorDomain int

const (
	prMalformedProfile MCProfileErrorDomain = iota + 1000
	prUnsupportedProfileVersion
	prMissingRequiredField
	prBadDataTypeInField
	prBadSignature
	prEmptyProfile
	prCannotDecrypt
	prNonUniqueUUIDs
	prNonUniquePayloadIdentifiers
	prProfileInstallationFailure
	prUnsupportedFieldValue
)

type MCPayloadErrorDomain int

const (
	plMalformedPayload MCPayloadErrorDomain = iota + 2000
	plUnsupportedPayloadVersion
	plMissingRequiredField
	plBadDataTypeInField
	plUnsupportedFieldValue
	plInternalError
)

type MCRestrictionsErrorDomain int

const (
	rsInconsistentRestrictionSense MCRestrictionsErrorDomain = iota + 3000
	rsInconsistentValueComparisonSense
)

type MCInstallationErrorDomain int

const (
	inCannotParseProfile MCInstallationErrorDomain = iota + 4000
	inInstallationFailure
	inDuplicateUUID
	inProfileNotQueued
	inUserCancelled
	inPasscodeNotCompliant
	inProfileRemovalDateInPast
	inUnrecognisedFileFormat
	inMismatchedCertificates
	inDeviceLocked
	inUpdatedProfileWrongIdentifier
	inFinalProfileNotConfiguration
	inProfileNotUpdatable
	inUpdateFailed
	inNoDeviceIdentity
	inReplacementNoMDMPayload
	inInternalError
	inMultipleHTTPProxyPayloads
	inMultipleCellularPayloads
	inMultipleApplockPayloads
	inUIInstallProhibited
	inProfileMustBeNonInteractive
	inProfileMustBeInstalledByMDM
	inUnnacceptablePayload
	inProfileNotFound
	inInvalidSupervision
	inRemovalDateInPast
	inProfileRequiresPasscodeChange
	inMultipleHomeScreenPayloads
	inMultipleNotificationPayloads
	inUnacceptablePayloadMultiuser
	inPayloadContainsSensitiveInfo
)

type MCPasscodeErrorDomain int

const (
	pcPasscodeTooShort MCPasscodeErrorDomain = iota + 5000
	pcTooFewUniqueChars
	pcTooFewComplexChars
	pcRepeatingChars
	pcAscendingDescendingChars
	pcRequiresNumber
	pcRequiresAlpha
	pcPasscodeExpired
	pcPasscodeTooRecent
	_
	pcDeviceLocked
	pcWrongPasscode
	_
	pcCannotClearPasscode
	pcCannotSetPasscode
	pcCannotSetGracePeriod
	pcCannotSetFingerprintUnlock
	pcCannotSetFingerprintPurchase
	pcCannotSetMaxAttempts
)

type MCKeychainErrorDomain int

const (
	kcKeychainSystemError MCKeychainErrorDomain = iota + 6000
	kcEmptyString
	kcCannotCreateQuery
)

type MCEmailErrorDomain int

const (
	emHostUnreachable MCEmailErrorDomain = iota + 7000
	emInvalidCredentials
	emUnknownValidationError
	emSMIMECertificateNotFound
	emSMIMECertificateBad
	emIMAPMisconfigured
	emPOPMisconfigured
	emSMTPMisconfigured
)

type MCWebClipErrorDomain int

const (
	wcCannotInstallWebClip MCWebClipErrorDomain = iota + 8000
)

type MCCertificateErrorDomain int

const (
	ceInvalidPassword MCCertificateErrorDomain = iota + 9000
	ceTooManyCertificatesInPayload
	ceCannotStoreCertificate
	ceCannotStoreWAPIData
	ceCannotStoreRootCertificate
	ceCertificateMalformed
	ceCertificateNotIdentity
)

type MCDefaultsErrorDomain int

const (
	deCannotInstallDefaults MCDefaultsErrorDomain = iota + 10000
	deInvalidSigner
)

type MCAPNErrorDomain int

const (
	apnCannotInstallAPN MCAPNErrorDomain = iota + 11000
	apnCustomAPNAlreadyInstalled
)

type MCMDMErrorDomain int

const (
	mdmInvalidAccessRights MCMDMErrorDomain = iota + 12000
	mdmMultipleMDMInstances
	mdmCannotCheckIn
	mdmInvalidChallengeResponse
	mdmInvalidPushCertificate
	mdmCannotFindCertificate
	mdmRedirectRefused
	mdmNotAuthorized
	mdmMalformedRequest
	mdmInvalidReplacementProfile
	mdmInternalConsistencyError
	mdmInvalidMDMConfiguration
	mdmMDMReplacementMismatch
	mdmProfileNotManaged
	mdmProvisioningProfileNotManaged
	mdmCannotGetPushToken
	mdmMissingIdentity
	mdmCannotCreateEscrowKeybag
	mdmCannotCopyEscrowKeybagData
	mdmCannotCopyEscrowSecret
	mdmUnauthorizedByServer
	mdmInvalidRequestType
	mdmInvalidTopic
)

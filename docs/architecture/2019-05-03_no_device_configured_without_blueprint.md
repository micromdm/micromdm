# Do not send DeviceConfigured when there are no user configured blueprints.

## Context

DEP enrolled devices use an `AwaitDeviceConfigured` mode, which can be optionally enabled in a DEP profile. If enabled, the device will wait at the Remote Management screen, and not proceed until the MDM issues a `DeviceConfigured` command.
MicroMDM introduced the concept of a `Blueprint`, which an administrator can configure to apply `InstallProfile` and `InstallApplication` payloads. At the end of the list of pre-configured commands the MDM would default to sending the `DeviceConfigured` command. 

It would be useful to allow a separate service to make a decision about whether it is ok to proceed with device provisioning past enrollment.

## Decision

When there are no blueprints configured, the MicroMDM server will not send _any_ commands to the device. 
That means, if `AwaitDeviceConfigured` is enabled on the device, upon enrollment the device will remain at the Remote Management screen until an API request queues the DeviceConfigured command.

## Status

Accepted

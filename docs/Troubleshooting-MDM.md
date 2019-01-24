# Troubleshooting MDM

Note: The DEP & MDM Testing + Troubleshooting is the most relevant for MDM developers, but please review the Sanity Check section first unless 100% confident!

## Emergency Apple Push Notification Services (APNs) certificate issues

Contact Apple support: <https://support.apple.com/en-us/HT208643>

## Sanity Check

### Apple system status

Make sure that Apple isn't currently reporting issues with Apple Business Manager/Apple School Manager, Device Enrollment Program, iOS Device Activation (<https://support.apple.com/en-us/ht201407>), Volume Purchase Program, and APNs using the links below. Note that the statuses of some or all of these services are manually updated, so the status indicators are unlikely to be accurate.

* ABM/ASM, DEP, iOS Activation, VPP: <https://www.apple.com/support/systemstatus/>
* APNs, APNs Sandbox: <https://developer.apple.com/system-status/>

## Don't forget about the DEP/ASM/ABM account and terms of service

If you're seeing Device Enrollment Protocol (DEP) or other Mobile Device Management (MDM) failures across the board and they don't seem to be immediately explainable, make sure to try logging into your DEP/ABM/ASM portal account and check if there are new terms of services that need to be accepted (<https://support.apple.com/en-us/HT203063>).

The DEP API should throw this error if new terms of service need to be accepted: `T_C_NOT_SIGNED`

## Checking for or requesting a DEP profile on macOS

Request a DEP profile for a macOS device in 10.13 and above:

```shell
sudo /usr/bin/profiles renew -type enrollment
```

Display the DEP profile for a macOS device in 10.13 and above:

```shell
sudo /usr/bin/profiles show -type enrollment
```

Request/show DEP profile for a macOS device (deprecated and not functional beyond 10.13.4)

```shell
sudo /usr/libexec/mdmclient dep nag
```

## Known issues with MDM

### storedownloadd

* `storedownloadd` queue gets stuck after an unsuccessful download and the process must be restarted for future downloads to work.

```shell
/var/log/commerce.log message when  storedownloadd is hanged:
May 21 17:13:40 admins-Mac storedownloadd[532]: DownloadManifest: removePurgeablePath: /var/folders/qg/yc8wn1f13jb_mly2mj34pfkh0000gn/C/com.apple.appstore/0
May 21 17:13:40 admins-Mac storedownloadd[532]: DownloadQueue: Could not add download
May 21 17:13:40 admins-Mac storedownloadd[532]: DownloadQueue: Will start any ready downloads anyways.
```

## Only one local user can be MDM enabled at a time

(Requires further details & citations) Unless using directory accounts, Apple's MDM protocol says that only one local user can be enabled for MDM. This has historically caused the most issues in some commercial solutions that attempt to install device-level VPP apps via the user channel [needs better citations].

This has been an issue since device-level VPP was introduced:

<https://www.jamf.com/jamf-nation/articles/372/enabling-mdm-for-local-user-accounts>

## DEP & MDM Testing + Troubleshooting

### Force DEP check in

Force a machine to check up on its DEP configuration record, display configuration contents, and prompt user to enroll in MDM if the machine currently isn't.

10.13.4 and above:

```shell
sudo /usr/bin/profiles renew -type enrollment
```

10.11 and above:

```shell
sudo /usr/libexec/mdmclient dep nag
```

### Check if a machine was enrolled via DEP (10.13+)

Show whether a machine has a device enrollment profile present (10.13.0+), and whether the MDM enrollment is user approved (10.13.4+)

```shell
/usr/bin/profiles status -type enrollment
```

### DEP & MDM Debug Logging

Enable MDM debug logging

For macOS 10.12+ use this configuration profile:

<https://gist.github.com/opragel/2b9c518f9a27dce787ed45da832708e2>

For OS X 10.11 and below, use the following commands to enable MDM debug logging:

```shell
sudo defaults write /Library/Preferences/com.apple.MCXDebug debugOutput -2
sudo defaults write /Library/Preferences/com.apple.MCXDebug collateLogs 1
sudo touch /var/db/MDM_EnableDebug
```

Logs of interest:

```shell
/Library/Logs/ManagedClient/ManagedClient.log
/var/log/commerce.log
```

Stream relevant logs on macOS 10.12+:

```shell
log stream --info --debug --predicate 'subsystem contains "com.apple.ManagedClient.cloudconfigurationd"'
log stream --info --debug --predicate 'processImagePath contains "mdmclient" OR processImagePath contains "storedownloadd"'
```

## Re-enroll a DEP Mac without having to wipe the device

10.13

Owen's big warning: Apple does not support this and as of 10.13 they made completely clearing the MDM profile if marked non-removable outside of recovery impossible - as in, switching MDMs without booting to recovery was impossible unless you took abnormally special care to put yourself in a place to migrate MDMs (you'll know if you have). I strongly recommend ignoring the script below unless it very coincidentally solves a big problem for you. This was only useful to me for a select few early 10.13 machines that were set to the right MDM server, but had the wrong local user locally enabled for MDM. It became a non-starter with UAMDM in 10.13.4, which will continue to be true even if Apple doesn't finish sealing off access by 10.14.

Recommendation: boot to recovery, perform corrections, boot back into the OS, and re-initiate MDM enrollment as necessary.

<https://gist.github.com/opragel/12555098f5894267c3aba2a7c023a823>

10.12

    # Remove indicator that setup assistant has already run
    sudo rm /var/db/.AppleSetupDone
    # Clear all configuration profiles off machine (not entirely clean)
    sudo rm -rf /var/db/ConfigurationProfiles/
    # Remove Apple Push Notification service daemon keychain
    sudo rm /Library/Keychains/apsd.keychain
    # Reboot the machine. It should bring you back to setup assistant
    # where you can re-enroll using DEP.

## Setup Assistant

### Run Setup Assistant as an app in regular user space

```shell
sudo /System/Library/CoreServices/Setup\ Assistant.app/Contents/MacOS/Setup\ Assistant -MBDebug
```

### Setup Assistant Troubleshooting

* Press Command + Option + Control + T during Setup Assistant to open Terminal.
* Press Command + Option + Control + C during Setup Assistant to open Console.

## Check the Topic for APNS

The topic in the cert is readable via

```shell
cat mdm.cert.pem |openssl x509 -subject -noout
```

The output of that will look something like

```shell
subject= /UID=com.apple.mgmt.External.foobar/CN=APSP:foobar/C=US
```

And that should match the `Topic` key in your MDM Enrollment profile. If they do not match, delete the enrollment profile on the server and re-generate it. You will need to re-sign the profile, example here: [Sign the enrollment profile with Hancock](https://github.com/micromdm/micromdm/wiki/Sign-the-enrollment-profile-with-Hancock).

_Delete the profile:_

```shell
mdmctl delete profiles -id com.github.micromdm.micromdm.enroll
```

_Re-generate the profile:_

```shell
curl https://mmdm.example.org/mdm/enroll
```

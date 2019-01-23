# Introduction

This guide is a Quickstart tutorial for getting up and running with MicroMDM.
The steps in this guide are intended to get you up and running as quick as possible, and therefore will not necessarily match the steps required to set up a long running production instance.

# Getting Help

The best place to get help is the `#micromdm` channel on the MacAdmins Slack team. Join us there by getting an [invitation here](https://macadmins.herokuapp.com/).

# Getting an MDM Push Certificate

An MDM is pretty much useless without an MDM Push Certificate.
Depending on your setup, this can be a long process.
For an overview of getting an MDM push certificate (and all the other MDM-related certificates) see [this blog post](https://micromdm.io/blog/certificates/).

The full process for obtaining an APNs **Vendor** certificate is documented in a [blog post](http://enterprisemac.bruienne.com/2015/06/06/mdm-azing-setting-up-your-own-mdm-server/) by Pepijn Bruienne. Use the [mdmctl mdmcert](https://github.com/micromdm/micromdm/wiki/Generating-MicroMDM-MDM-Certificates) command for performing the same steps in MicroMDM.



Meanwhile, there are two other ways to obtain a certificate that might prove to be more expedient if you're looking to test and/or contribute to MicroMDM:

a) If you already have the ServerApp with Device Management enabled, you can export the APNS certificate from Keychain Access. The certificate looks like `APSP:uuid` and has a UserID of `com.apple.mgmt.XServer:uuid`. Do **NOT** use this certificate for anything other than testing or development. You will have trouble renewing it when it expires.

b) Jesse Peterson has set up an MDM certificate service [here](https://mdmcert.download/). You can register and then download the certificate with the [`certhelper`](https://github.com/micromdm/tools/releases/latest) binary.
If you're using the `certhelper` tool to request a certificate through `mdmcert.download`, follow the steps on the wiki page [here](https://github.com/micromdm/micromdm/wiki/mdmcert.download).

# Up and Running with MicroMDM

1) Download the latest [binary](https://github.com/micromdm/micromdm/releases/latest) from the Releases page and copy it somewhere in your PATH. (usr/local/bin) for example.
2) Run the binary.

```shell
sudo micromdm serve \
  -server-url=https://my-server-url \
  -api-key MySecretAPIKey \
  -filerepo /path/to/pkg/folder \
  -tls-cert /path/to/tls.crt \
  -tls-key  /path/to/tls.key
```

If the server URL is a internet routable IP address, micromdm will automatically bind to ports `443` and `80` and download a certificate from Let's Encrypt to use for HTTPS.

**Attention!** The above is an awesome feature, but it is currently broken. So you should manage the certificates yourself. You can do it either with [Lego](https://github.com/micromdm/micromdm/wiki/Generating-LetsEncrypt-Certs-with-the-Lego-Client) or [certbot](https://github.com/micromdm/micromdm/wiki/Using-Certbot-with-MicroMDM).

If you're running behind a firewall or don't want to bind with the default ports, you can pass a few extra flags, see `micromdm serve -h` for a full list of options for the `serve` command.

If everything went well, the binary should be running. You can go to `https://your-server-url/mdm/enroll` which will download an enrollment profile for your Mac.

# Connecting mdmctl
To setup mdmctl to connect to the micromdm server you need to set some configuration parameters (which save to `~/.micromdm/servers.json`):

```shell
mdmctl config set \
  -name production \
  -api-token MySecretAPIKey \
  -server-url https://my-server-url \
```

You'll also need to initially "select" this configuration by switching to it. THis is the mechanism by which you can multiple micromdm servers configured and you switch between them:

```shelp
mdmctl config switch -name production
```

You'll also need to provide the `-api-key` switch with the same password to the `micromdm serve` command:

```shell
sudo micromdm serve \
  -api-key MySecretAPIKey \
```

# Connecting to DEP
Note: Bootstrapping management tools like Munki are probably the reason most of us want an MDM right now, but DEP is not necessary for MicroMDM to work. You can skip to the "Profiles and Applications Section" if you don't have a DEP account.

Got a DEP account? You can set up a virtual server for MicroMDM and sync your devices:
https://github.com/micromdm/micromdm/wiki/Connect-MicroMDM-with-DEP
This must be done before establishing DEP profiles for devices.

Once you're connected to DEP and can sync devices, you can assign a profile to them:

Use `mdmctl` to generate a sample profile.json. Note the list of serial numbers to assign the profile to.

Generate your own GUID for the org_magic string.

```shell
mdmctl apply dep-profiles -template > /tmp/profile.json
{
   "profile_name": "Test Profile",
   "url":"https://dev.micromdm.io/mdm/enroll",
   "await_device_configured":false,
   "is_mdm_removable":true,
   "department": "IT Department",
   "org_magic": "913FABBB-0032-4E13-9966-D6BBAC900331",
   "support_phone_number": "1-555-555-5555",
   "support_email_address": "org-email@example.com",
   "skip_setup_items": [ "Registration", "AppleID", "TOS"],
   "devices": ["SERIAL1","SERIAL2"]
}
```

Assign the profile:

```shell
$ mdmctl apply dep-profiles -f /tmp/profile.json
Defined DEP Profile with UUID 4B05B09E8AC7E7FC12C8F3338E099310

$ mdmctl get dep-profiles -f - -uuid=4B05B09E8AC7E7FC12C8F3338E099310
{
  "profile_name": "Test Profile",
  "url": "https://mdm.acmeinc.com/getconfig",
  "is_mdm_removable": true,
  "support_phone_number": "1-555-555-5555",
  "support_email_address": "org-email@example.com",
  "org_magic": "913FABBB-0032-4E13-9966-D6BBAC900331",
  "skip_setup_items": [
    "Registration",
    "AppleID",
    "TOS"
  ]
}
```

# Profiles and Applications
DEP enrollment is cool and all, but we want to be able to enforce some kind of configuration when a new device enrolls. The easiest way to do that is with a Blueprint. We're working on making blueprints configurable for each device, but for now this blueprint is global - it applies to every device enrolled.

## Profiles
MicroMDM stores profiles in its database for use with Blueprints. To store a profile use this invocation:

```shell
$ mdmctl apply profiles -f /path/to/MyProfile.mobileconfig
```

```shell
$ mdmctl get profiles
Identifier             Length
com.example.MyProfile  874
```

## Applications
MDM is capable of installing [distribution style pkgs](https://github.com/micromdm/micromdm/wiki/distributing-packages-with-InstallApplication) if they are signed. MicroMDM attempts to make this option simple for the user. the `mdmctl` command can sign and upload a pkg file to the MicroMDM server's filerepo directory.

Assuming you have a a developer certificate to sign the pkg in the keychain(or if the package is already signed), here is how you'd import munkitools:
```
mdmctl apply app -pkg ~/Desktop/mdmvid/munkitools-3.0.0.3333.pkg -sign "Developer ID Installer: groob (myid)" -upload
```

`mdmctl` will automatically generate an app manifest file and upload it to the server.
You can use the `mdmctl get apps` command to query the list of imported apps and get the ManifestURL, which will be required by a blueprint.

## Blueprints
A blueprint is a JSON file which contains a list of Application Manifest URLs and Profile identifiers and an array of apply-at directives.

The blueprint JSON file contains a few keys. To generate a sample file use the `-template flag`:

```shell
$ mdmctl apply blueprints -template > /tmp/blueprint.json
{
  "uuid": "a480e3ee-a96b-4a4d-abcd-1af07229cc7e",
  "name": "exampleName",
  "install_application_manifest_urls": [
    "https://mdm.example.org/repo/exampleAppManifest.plist"
  ],
  "profile_ids": [
    "com.example.my.profile"
  ],
  "apply_at": [
    "Enroll"
  ]
}
```


Both the manifest and pkg file that the manifest points to must be downloadable from a public URL. You can put them on S3 or any other fileserver. MicroMDM also supports serving a repository of files with the `-filerepo /path/to/folder` flag of the `serve` command.
The profile paths can be anywhere on the computer where micromdm is running.

Modify a profile to include URLs paths to app manifests and profile IDs that have been uploaded:


```shell
$ mdmctl apply blueprints -f /path/to/blueprint

$ mdmctl get blueprints
Name         UUID                     Manifests  Profiles  Apply At
exampleName  a480e3ee-a96b-4a4d-abcd  1          1         Enroll
```

The "Apply At" (or "apply_at" in the JSON above) tells the when the MDM server will apply this Blueprint (that is - install the profiles and applications). Currently the only ApplyAt value is Enroll which happens at, well, Enrollment time (either DEP or standard enrollment).

# Commands and Push API.

MicroMDM supports a custom commands API that you can use to queue commands.
To do so, send a JSON command as a POST request to `/v1/commands`.
The API supports the vast majority of MDM commands specified in the MDM protocol reference, but the API still needs full documentation (See Issue [#78](https://github.com/micromdm/micromdm/issues/78)).

```shell
POST /v1/commands
{
    "request_type": "InstallApplication",
    "udid" : "184012D9-753A-5DFC-8149-5C9AF257629F",
    "manifest_url" : "https://mdm.example.com/repo/manifests/munkitools-2.5.1.2637.plist",
    "management_flags" : 1
}
```

For both scheduling a command and sending a push notification, you will need to know the device's UDID.
The simplest way to do that is to run

```shell
mdmctl get devices
```

which will print some details about a device, including the UDID.

Note you may need to authenticate to the API if you specified an API key with the `-api-key` parameter. If this is on the API will require HTTP Basic authentication. The username is "micromdm" and the password will be your API key.

## Sending a Push notification to the device.

The simplest way to issue a push notification and force the device to check in is a `GET` request to `/push/$DEVICE-UDID-HERE`.

Example: `https://dev.micromdm.io/push/FA01680E-98CA-5557-8F59-7716ECFEE964`

Like the command API you may need to specify HTTP Basic authentication credentials if you specified an API key with the `-api-key` parameter.  See command section above for more information.

## Further Actions

[Adding Apple Push Certificates](https://github.com/micromdm/micromdm/wiki/Generating-MicroMDM-MDM-Certificates)

[Webhooks](https://github.com/micromdm/micromdm/wiki/Webhooks)

## Additional Information/Videos

**Getting MicroMDM working and working with MicroMDM**

[![Getting MicroMDM working and working with MicroMDM](http://img.youtube.com/vi/WGKT-PyHz6I/0.jpg)](http://www.youtube.com/watch?v=WGKT-PyHz6I)

# mdmctl Command Reference

| sub-command       | description |
| -------------     | ------------- |
| [get](#mdmctl-get)                           | Display one or many resources  |
| [apply](#mdmctl-apply)                       | Apply a resource  |
| [config](#mdmctl-config)                     | Display or set configuration  |
| [remove](#mdmctl-remove)                     | Remove a resource  |
| [mdmcert](#mdmctl-mdmcert)                   | Create new MDM Push Certificate  |
| [mdmcert.download](#mdmctl-mdmcert.download) | Request new MDM Push Certificate from https://mdmcert.download  |
| [version](#mdmctl-version)                   | Display mdmctl version  |


***
<!--###### get ######################################################### -->

<a id="mdmctl-get"></a>
## get
_Display one or many resources_

<details>

<br/>
Valid resource types:

  * [devices](#mdmctl-get-devices)
  * [blueprints](#mdmctl-get-blueprints)
  * [dep-tokens](#mdmctl-get-dep-tokens)
  * [dep-devices](#mdmctl-get-dep-devices)
  * [dep-account](#mdmctl-get-dep-account)
  * [dep-profiles](#mdmctl-get-dep-profiles)
  * [dep-autoassigner](#mdmctl-get-dep-autoassigner)
  * [users](#mdmctl-get-users)
  * [profiles](#mdmctl-get-profiles)
  * [apps](#mdmctl-get-apps)

<a id="mdmctl-get-devices"></a>
##### devices
```
USAGE
  mdmctl get devices [flags]

FLAGS
  -serials   comma seperated list of serials to search
```

<a id="mdmctl-get-blueprints"></a>
##### blueprints
```
USAGE
  mdmctl get blueprints [flags]

FLAGS
  -f -    filename of JSON to save to
  -name   name of blueprint
```

<a id="mdmctl-get-dep-tokens"></a>
##### dep-tokens
```
USAGE
  mdmctl get dep-tokens [flags]

FLAGS
  -export-public-key mdm-files/DEPPublicKey   Filename of public key to write (to be uploaded to deploy.apple.com)
  -export-token mdm-files/DEPOAuthToken.json  Filename to save decrypted oauth token (JSON)
  -v false                                    Display full ConsumerKey in summary list
```

<a id="mdmctl-get-dep-devices"></a>
##### dep-devices
```
USAGE
  mdmctl get dep-devices [flags]

FLAGS
  -serials   comma separated list of device serials
```

<a id="mdmctl-get-dep-account"></a>
##### dep-account
```
USAGE
  mdmctl get dep-account [flags]

FLAGS
```

<a id="mdmctl-get-dep-profiles"></a>
##### dep-profiles
```
USAGE
  mdmctl get dep-profiles [flags]

FLAGS
  -f      filename of DEP profile to apply
  -uuid   DEP Profile UUID(required)
```

<a id="mdmctl-get-autoassigner"></a>
##### dep-autoassigner
```
USAGE
  mdmctl get dep-autoassigner [flags]

FLAGS

```

<a id="mdmctl-get-users"></a>
##### users
```
USAGE
  mdmctl get users [flags]

FLAGS

```

<a id="mdmctl-get-profiles"></a>
##### profiles
```
USAGE
  mdmctl get blueprints [flags]

FLAGS
  -f -  filename of profile to write
  -id   profile Identifier
```

<a id="mdmctl-get-apps"></a>
##### apps
```
USAGE
  mdmctl get apps [flags]

FLAGS
  -f -    path to save file to. defaults to stdout.
  -name   specify the name of the app to get full details
```

</details>

***
<!--###### apply ######################################################### -->

<a id="mdmctl-apply"></a>
## apply
_Apply a resource_

<details>

<br/>
Valid resource types:

  * [blueprints](#mdmctl-apply-blueprints)
  * [profiles](#mdmctl-apply-profiles)
  * [users](#mdmctl-apply-users)
  * [dep-tokens](#mdmctl-apply-dep-tokens)
  * [dep-profiles](#mdmctl-apply-dep-profiles)
  * [dep-autoassigner](#mdmctl-apply-dep-autoassigner)
  * [app](#mdmctl-apply-app)
  * [block](#mdmctl-apply-block)

<a id="mdmctl-apply-blueprints"></a>
##### blueprints
```
USAGE
  mdmctl apply blueprints [flags]

FLAGS
  -f               filename of blueprint JSON to apply
  -template false  print a new blueprint template
```

<a id="mdmctl-apply-profiles"></a>
##### profiles
```
USAGE
  mdmctl apply profiles [flags]

FLAGS
  -f   filename of profile to apply
```

<a id="mdmctl-apply-users"></a>
##### users
```
USAGE
  mdmctl apply users [flags]

FLAGS
  -f               Path to user manifest
  -password        Password of the user. Only required when creating a new user.
  -template false  Print a JSON example of a user manifest.
```

<a id="mdmctl-apply-dep-tokens"></a>
##### dep-tokens
```
USAGE
  mdmctl apply dep-tokens [flags]

FLAGS
  -import mdm-files/DEPOAuthToken.json  Filename of p7m encrypted token file (downloaded from DEP portal)
```

<a id="mdmctl-apply-dep-profiles"></a>
##### dep-profiles
```
USAGE
  mdmctl apply dep-profiles [flags]

FLAGS
  -anchor                 filename of PEM cert(s) to add to anchor certs in template
  -f                      filename of DEP profile to apply
  -filter                 set the auto-assign filter to for the defined profile
  -template false         print a JSON example of a DEP profile
  -use-server-cert false  use the server cert(s) to add to anchor certs in template
```

<a id="mdmctl-apply-dep-autoassigner"></a>
##### dep-autoassigner
```
USAGE
  mdmctl apply dep-autoassigner [flags]

FLAGS
  -filter *  filter string (only '*' supported right now)
  -uuid      DEP profile UUID to set
```

<a id="mdmctl-apply-app"></a>
##### app
```
USAGE
  mdmctl apply app [flags]

FLAGS
  -manifest -        path to an app manifest. optional,
                     will be created if file does not exist.
  -md5size 10485760  md5 hash size in bytes (optional)
  -pkg               path to a distribution pkg.
  -pkg-url           use custom pkg url
  -sign              sign package before importing, requires specifying a product ID (optional)
  -upload false      upload package and/or manifest to micromdm repository.
```

<a id="mdmctl-apply-block"></a>
##### block
```
USAGE
  mdmctl apply block [flags]

FLAGS
  -udid   UDID of a device to block.
```

</details>

***
<!--###### config ######################################################### -->

<a id="mdmctl-config"></a>
## config
_Display or set configuration_

<details>

<br/>
Valid switches:

  * [print](#mdmctl-config-print)
  * [set](#mdmctl-config-set)
  * [switch](#mdmctl-config-switch)

<a id="mdmctl-config-print"></a>
##### print
Display current configuration
```
USAGE
  mdmctl config print
```

<a id="mdmctl-config-set"></a>
##### set
```
USAGE
  mdmctl config set [flags]

FLAGS
  -api-token          api token to connect to micromdm server
  -name               name of the server
  -server-url         server url of micromdm server
  -skip-verify false  skip verification of server certificate (insecure)
```

<a id="mdmctl-config-switch"></a>
##### switch
```
USAGE
  mdmctl config switch [flags]

FLAGS
  -name   name of the server to switch to
```

</details>

***
<!--###### remove ######################################################### -->

<a id="mdmctl-remove"></a>
## remove
_Remove a resource_

<details>

<br/>
Valid resource types:

  * [blueprints](#mdmctl-remove-blueprints)
  * [devices](#mdmctl-remove-devices)
  * [profiles](#mdmctl-remove-profiles)
  * [block](#mdmctl-remove-block)
  * [dep-autoassigner](#mdmctl-remove-dep-autoassigner)

<a id="mdmctl-remove-blueprints"></a>
##### blueprints
```
USAGE
  mdmctl remove blueprints [flags]

FLAGS
  -name   name of blueprint, optionally comma separated
```

<a id="mdmctl-remove-devices"></a>
##### devices
```
USAGE
  mdmctl remove devices [flags]

FLAGS
  -udid   device UDID, optionally comma separated
```

<a id="mdmctl-remove-profiles"></a>
##### profiles
```
USAGE
  mdmctl remove profiles [flags]

FLAGS
  -id   profile Identifier, optionally comma separated
```

<a id="mdmctl-remove-block"></a>
##### block
```
USAGE
  mdmctl remove block [flags]

FLAGS
  -udid   UDID of device to unblock
```

<a id="mdmctl-remove-dep-autoassigner"></a>
##### dep-autoassigner
```
USAGE
  mdmctl remove dep-autoassigner [flags]

FLAGS
  -filter *  filter string (only '*' supported right now)
```

</details>

***
<!--###### mdmcert ######################################################### -->

<a id="mdmctl-mdmcert"></a>
## mdmcert
_Create new MDM Push Certificate_

<details>

<br/>This utility helps obtain a MDM Push Certificate using the Apple Developer MDM CSR option in the enterprise developer portal.

Valid switches:

  * [vendor](#mdmctl-mdmcert-vendor)
  * [push](#mdmctl-mdmcert-push)
  * [upload](#mdmctl-mdmcert-upload)

<a id="mdmctl-mdmcert-vendor"></a>
##### vendor
```
USAGE
    mdmctl mdmcert vendor [flags]

FLAGS
    -cert mdm-certificates/mdm.cer                         Path to the MDM Vendor certificate from dev portal.
    -cn micromdm-vendor                                    CommonName for the CSR Subject.
    -country US                                            Two letter country code for the CSR Subject(example: US).
    -email                                                 Email address to use in CSR Subject.
    -out mdm-certificates/VendorCertificateRequest.csr     Path to save the MDM Vendor CSR.
    -password                                              Password to encrypt/read the RSA key.
    -private-key mdm-certificates/VendorPrivateKey.key     Path to the vendor private key. A new RSA key will be created at this path.
    -push-csr mdm-certificates/PushCertificateRequest.csr  Path to the user CSR(required for the -sign step).
    -sign false                                            Signs a user CSR with the MDM vendor certificate.
```

<a id="mdmctl-mdmcert-push"></a>
##### push
```
USAGE
    mdmctl mdmcert push [flags]

FLAGS
    -cn micromdm-user                                            CommonName for the CSR Subject.
    -country US                                                  Two letter country code for the CSR Subject(Example: US).
    -email                                                       Email address to use in CSR Subject.
    -out mdm-certificates/PushCertificateRequest.csr             Path to save the MDM Push Certificate request.
    -password                                                    Password to encrypt/read the RSA key.
    -private-key mdm-certificates/PushCertificatePrivateKey.key  Path to the push certificate private key. A new RSA key will be created at this path.
```

<a id="mdmctl-mdmcert-upload"></a>
##### upload
```
USAGE
    mdmctl mdmcert upload [flags]

FLAGS
    -cert                                                        Path to the MDM Push Certificate.
    -password                                                    Password to encrypt/read the RSA key.
    -private-key mdm-certificates/PushCertificatePrivateKey.key  Path to the push certificate private key.
```

First you must create a vendor CSR which you will upload to the enterprise developer portal and get a signed MDM Vendor certificate. Use the MDM-CSR option in the dev portal when creating the certificate.
The MDM Vendor certificate is required in order to obtain the MDM push certificate. After you complete the MDM-CSR step, copy the downloaded file to the same folder as the private key. By default this will be
mdm-certificates

```mdmctl mdmcert vendor -password=secret -country=US -email=admin@acme.co```

Next, create a push CSR. This step generates a CSR required to get a signed a push certificate.

```mdmctl mdmcert push -password=secret -country=US -email=admin@acme.co```

Once you created the push CSR, you mush sign the push CSR with the MDM Vendor Certificate, and get a push certificate request file.

```mdmctl mdmcert vendor -sign -cert=./mdm-certificates/mdm.cer -password=secret```

Once generated, upload the PushCertificateRequest.plist file to https://identity.apple.com to obtain your MDM Push Certificate.
Use the push private key and the push cert you got from identity.apple.com in your MDM server.

</details>

***
<!--###### mdmcert.download ######################################################### -->
<a id="mdmctl-mdmcert.download"></a>
## mdmcert.download
_Request new MDM Push Certificate from https://mdmcert.download_

<details>

```
USAGE
  mdmctl mdmcert.download [flags]

FLAGS
  -cn mdm-push                                 CommonName for the CSR Subject.
  -country US                                  Two letter country code for the CSR Subject (example: US).
  -decrypt                                     Decrypts and mdmcert.download push certificate request
  -email                                       Email address to use in mdmcert request & CSR Subject
  -new false                                   Generates a new privkey and uploads new MDM request
  -pki-cert mdmcert.download.pki.crt           Path for generated MDMCert pki exchange certificate
  -pki-password                                Password to encrypt/read the RSA key.
  -pki-private-key mdmcert.download.pki.key    Path for generated MDMCert pki exchange private key
  -push-csr mdmcert.download.push.csr          Path for generated Push Certificate CSR
  -push-password                               Password to encrypt/read the push RSA key.
  -push-private-key mdmcert.download.push.key  Path to the generated Push Cert private key
  -push-req mdmcert.download.push.req          Path for generated Push Certificate Request
```

</details>

***
<!--###### version ######################################################### -->

<a id="mdmctl-version"></a>
## version
_Display mdmctl version_

<details>

```
USAGE
  mdmctl version
```

</details>
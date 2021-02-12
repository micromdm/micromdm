The quickstart page is intended to provide a quick guide for setting up a running MicroMDM service. It does not cover all the possible deployment environments. More complex configuration details will be covered in their own sections. 

Before proceeding, make sure to read the [Introduction](./introduction.md) page, which describes requirements you might need before using MicroMDM.

# Up and Running with MicroMDM

First, [download](https://github.com/micromdm/micromdm/releases/latest) the latest release of MicroMDM. Copy the `micromdm` binary somewhere in your `$PATH`. 

Next, run the server binary. 

```
sudo micromdm serve \
  -server-url=https://my-server-url \
  -api-key MySecretAPIKey \
  -filerepo /path/to/pkg/folder \
  -tls-cert /path/to/tls.crt \
  -tls-key  /path/to/tls.key
```

- `server-url` MUST be the `https://` URL you(and your devices) will use to connect to MicroMDM.
- `api-key` is a secret you MUST create to protect the API. It will be used to authenticate API requests both from your own integrations, as well as `mdmdctl`. 
- `filerepo` is an optional key which needs to point to a directory `micromdm` can read and write to. It is used for packages uploaded by `mdmctl apply app`. 
 It is not necessary if you do not intend to push custom packages via InstallApplication commands.
- `tls-*` flags are used to specify the path to your TLS certificate and key. You MUST configure them if `micromdm` will terminate connections to a device.  
If you prefer to terminate TLS with a loadbalancer, you can set `-tls=false`. 
**NOTE**: You must use the `-flag=false` form to turn off a boolean flag.

Command line flags are used for configuration, because they are easy to discover and document.  
For a full list of available configurations and their usage, run `micromdm -help`.  
In some cases you can choose to use environment variables to provide the same configuration to `micromdm`. Each flag has a corresponding `MICROMDM_CONFIG_FLAG` environment variable option.
For example, `-api-key` becomes `MICROMDM_API_KEY`. 

**NOTE**: In a production environment, secrets should _always_ be set via environment variables and not CLI flags, so they don't remain in your shell history or server monitoring logs.  

This section described how to start the `micromdm` process interactively in a shell, but that won't persist a server restart or exiting your current session. Having the process remain persistent depends on your environment. Since systemd is a common choice, there are [notes](https://github.com/micromdm/micromdm/wiki/Using-MicroMDM-with-systemd) from users on the wiki. 

# Configure mdmctl

At this point, you should have a running server process, but you do not yet have a working MDM service. Some important configuration is still required. MicroMDM does not come with a built in web interface for administrators. Instead, it comes with a CLI utility called `mdmctl` which uses the API to provide helpful commands for the MDM administrator.  
The `mdmctl` binary does not need to run on the same computer as `micromdm`. It is usually convenient to run it on your laptop, but you can of course run it anywhere you want, as long as it can reach the server URL.

To use `mdmctl`, we will first have to configure it. 

```
mdmctl config set \
  -name production \
  -api-token MySecretAPIKey \
  -server-url https://my-server-url
```

The value for `-api-key` you set on the server is the `-api-token` value you're setting for `mdmctl`. 

And immediately after, run: 

```
mdmctl config switch -name production
```

The configuration `mdmctl` uses lives in the `~/.micromdm/servers.json` file and you can view it in a text editor, or with the `mdmctl config print` command. 

You're likely to run more than one instance of micromdm(ex: production and staging). You can store several configurations with the `-name` flag, and then switch between them with `mdmctl config switch` command.

`mdmctl` has many helpul commands. You can discover them with the `-help` flag. Many actions are grouped by subcommands like `get`, `apply`, `remove`. For example, `mdmctl get devices` will list devices enrolled in the MDM.

# Configure an APNS certificate

To communicate with your device fleet, MDM needs an APNS certificate issued by Apple. As noted in the introduction, this process requires that you have an Enterprise Developer Account, and the `MDM CSR` option enabled under the *Certificates, IDs & Profiles* tab for *iOS*. 

Apple has a separate flow for the MDM vendor than the one for customers. For an in-house deployment without third parties, you must complete both the vendor and the customer process yourself. The `mdmctl mdmcert` command will help you with your APNS certificate needs.

Create a request for the MDM CSR, with a password used to encrypt the private key. 
After this step you will have a new `mdm-certificates` directory, with the necessary files. 
```
mdmctl mdmcert vendor -password=secret -country=US -email=admin@acme.co
```
### Generate MDM CSR
Log in to the Apple Developer Portal (https://developer.apple.com/account), and navigate to the Certificates, IDs & Profiles section (https://developer.apple.com/account/resources/certificates/list).
  1. Click the plus symbol (+) next to *Certificates*
  2. Select *MDM CSR* under the *Services* section, click *Continue*
  3. Upload the `VendorCertificateRequest.csr` file, click *Continue*
  4. Click *Download* to download the certificate.
  5. Move the downloaded certificate file (likely called mdm.cer) to the `mdm-certificates` folder. 

You now have the vendor side of the certificate flow complete, and you need to complete the customer side of this flow, with the help of the vendor cert. 

Sign a push certificate request with the vendor certificate. This step uses the private key you created above, so specify the same password to be able to decrypt it. 

```
mdmctl mdmcert vendor -sign -cert=./mdm-certificates/mdm.cer -password=secret
```

You should now have a `mdm-certificates/PushCertificateRequest.plist` file. 

Sign in to [identity.apple.com](https://identity.apple.com) and upload the `PushCertificateRequest.plist` file to get the APNS certificate. The site offers a notes field, it's a good idea to use it to mark which server you're obtaining a certificate for, as you will come back for renewals. 

If you're getting certificates for multiple environments (staging, production) or running multiple in house MDM instances, you MUST sign a separate push request for each one. Using the same vendor certificate is okay, but using the same push certificate is not. 

If you've uploaded the plist, you will be offered a certificate, which is signed for the `mdm-certificates/PushCertificatePrivateKey.key` key. Copy the certificate to the same directory. 

Finally, upload the certificate to MicroMDM.

```
mdmctl mdmcert upload \
    -cert mdm-certificates/MDM_\ Acme\,\ Inc._Certificate.pem \
    -private-key mdm-certificates/PushCertificatePrivateKey.key \
    -password=secret
```

Keep a backup of the `mdm-certificates` directory somewhere safe(like a 1Password vault) in case you need to use the vendor certificate again. But once the upload has completed, you don't need the local copy around on your computer, since MicroMDM stores the APNS key. 

See the renewals sections at the end of this document for renewals steps. 

# Configure Apple Business Manager (DEP)

Using DEP is not required for MicroMDM, but is supported. If you do not need it, skip this section. 

In this section you will use `mdmctl` to create a public key, upload it to the DEP portal and then upload the DEP token to MicroMDM. This exchange will link your MDM server with DEP, and allow you to sync devices and associate them with MDM. 

First, run:
```
mdmctl get dep-tokens -export-public-key /tmp/DEPPublicKey.pem
```

Next, log in to your deployment program website, and [follow the instructions](https://help.apple.com/businessmanager/#/asm1c1be359d) for adding an MDM server. Once you get to the step to upload the public key certificate, use the file from the previous step.  
You'll be asked to save a file containing the server token. 

Upload the token file to the server:

```
mdmctl apply dep-tokens \
    -import /path/to/Some_Server_Token_smime.p7
```

Now, your server should be capable of syncing devices. To verify, you can run 

```
mdmctl get dep-account
```

Which will return the information about your account. 

# Renewing certificates

At least once a year you will have to renew your DEP and MDM certificate connections with Apple. It's a good idea not to wait until expiration, but to set a calendar reminder somewhere 9-10 months from the day you configured your services. 

To renew the MDM certificate you will re-do the entire process as described in the APNS section. The vendor cert needs to be re-created as well. When you get to the step of requesting a push certificate from identity.apple.com, make sure you are using the *Renew* button in the portal. Doing so ensures you maintain the same APNS topic year to year, which is necessary to avoid re-enrollment.

To renew your DEP tokens, follow the same instructions as setting up for the first time. 

# Restart the server

Once you configured DEP and the MDM certificates you should be good to begin enrolling devices. 
To avoid any issues, restart the `micromdm` process once you configured everything. 

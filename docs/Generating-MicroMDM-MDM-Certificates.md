# Generating MicroMDM MDM Certificates

**mdmctl** is utility that helps obtain a MDM Push Certificate using the Apple Developer MDM CSR option in the enterprise developer portal.

If you dont have a copy of mdmctl, download one on the [Release Page](https://github.com/micromdm/micromdm/releases)

## Vendor CSR

First, you must create a vendor CSR which you will upload to the enterprise developer portal and get a signed MDM Vendor certificate. 

	mdmctl mdmcert vendor -password=secret -country=US -email=admin@acme.co

This will create a folder in your current path called `mdm-certificates` with the contents `VendorCertificateRequest.csr` and `VendorPrivateKey.key` which you will use to send to Apple.

>Note: It goes without saying, keep the VendorPrivateKey safe and secure.


## Apple Developer Portal (Request the MDM CSR)

Log into your [Apple Developer Portal](https://developer.apple.com/account).

On the **Certificates, IDs & Profiles** Tab, make sure you are looking at **iOS, tvOS, watchOS** tab, or the MDM CSR wont be visible.

>Note: If you don't see the MDM CSR item, you will need to reach out to your apple developer support representative and have that option enabled.

Create a new certficate by clicking the `+` in the top right corner, and select **MDM CSR.**

Click past the next screen, because you already generated the key and csr with the `mdmctl` tool in the first step.

Upload the **VendorCertificateRequest.csr**, and you will be presented with a Download button to get the signed **mdm.cer**

To stay organized, copy the **mdm.cer** to your `../mdm-certificates` directory

## MicroMDM Push Certificate Request

Next, create a push CSR. This step generates a CSR required to get a signed a push certificate.

>Note: This is using the same secret as the initial vendor CSR

    mdmctl mdmcert push -password=secret -country=US -email=admin@acme.co

This command creates two files **PushCertificateRequest.csr** and **PushCertificatePrivateKey.key**.

## Sign the Push CSR with the MDM Vendor Certificate

This command will sign the push csr with the mdm vendor certificate you downloaded from apple developer.

>Note: This is using the same secret as the initial vendor CSR

	mdmctl mdmcert vendor -sign -cert=./mdm-certificates/mdm.cer -password=secret

You now have **PushCertificateRequest.plist** in your mdm-certficates folder

## Obtain your MDM Push Certificate

Log into [identity.apple.com](https://identity.apple.com)

Create a new certificate.  In the Notes field make sure you declare what server this certificate is being used for (there is no other way to know).
    
Upload the **PushCertificateRequest.plist**

Download the **MDM_ Acme._Certificate.pem** file
    
Keep this file in your `mdm-certificates` folder to stay organized

## Upload the Push Certificate to MicroMDM

The following command will upload your push certificate to micromdm
```
mdmctl mdmcert upload -cert mdm-certificates/MDM_\ Acme\,\ Inc._Certificate.pem -private-key mdm-certificates/PushCertificatePrivateKey.key -password=secret
```
## Download the DEP Token from MicroMDM

Now you need to download the DEP Public Key from MicroMDM, and use that key to register with [deploy.apple.com](http://deploy.apple.com/)

```
mdmctl get dep-tokens -export-public-key mdm-certificates/DEPPublicKey.key
```

>Note: The name of the public key (in this case DEPPublicKey) is arbitrary.  You can name it what ever you like.

This downloaded the **DEPPublicKey.key** to `mdm-certificates` folder

## Register with Deploy.apple.com

Now that you have the public key of your mdm server, you can now talk to Apple Push Network services.

Log into [deploy.apple.com](http://deploy.apple.com/) or [school.apple.com](http://school.apple.com/)

Click on **Device Enrollment Program** and **Manage Servers** and **Add MDM Server.**

Give it a reasonable name.  Preferably do not auto assign new devices to it unless you are sure you want to.

Upload the **DepPublicKey**, and in turn, Download the **Server Token**

Move the server token into your `mdm-certificates` folder for organization.

Now upload the `.p7m` to micromdm

```
mdmctl apply dep-tokens -import mdm-certificates/Server\ Name_Token_2017-11-10T18-20-57Z_smime.p7m
```
and you should get a positive response `imported DEP token`


## Folder Structure

Here is a listing of all the files in the `mdm-certificates` folder:

```
└── mdm-certificates
    # Push Certificate downloaded from identity.apple.com (final step).
    ├── MDM_\ Acme\ Inc_Certificate.pem 
    # Push Private key from mdmctl mdmcert push
    ├── PushCertificatePrivateKey.key
    # File to upload to identity.apple.com. This is the output from mdmctl mdmcert vendor -sign
    ├── PushCertificateRequest
    # CSR for push certificate. Created by running mdmctl mdmcert push
    ├── PushCertificateRequest.csr
    # CSR and Private key created by the mdmctl mdmcert vendor command in the first step.
    ├── VendorCertificateRequest.csr
    ├── VendorPrivateKey.key
    # Certificate from the enterprise portal. Downloaded after uploading the VendorCertificateRequest.csr for the 
    # MDM-CSR option. 
    └── mdm.cer

```

`MDM_\ Acme\ Inc_Certificate.pem` and `PushCertificatePrivateKey.key` are your push certificate and private key.

## Good to Know

**MDM CSR**: you only get 2 per developer account.  You actually only need 1, and can spin up instances of MicroMDM off of that vendor cert.
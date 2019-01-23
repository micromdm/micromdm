**If you're using a recent MicroMDM build (as of May 2nd, 2018), there's now a newer, much easier process to follow. If you're using a MicroMDM 1.2.0 or earlier, or an earlier source version scroll down to the legacy `certhelper` section.**

# Using the `mdmctl` method

Follow the [mdmcert.download instructions](https://mdmcert.download/instructions) until you reach step #4 (using your open source MDM software, in this case MicroMDM). Basically sign up and verify your email.

Then, create a new request with `mdmctl`:

```
$ mdmctl mdmcert.download -new -email=cool.mdm.admin@example.org
Request successfully sent to mdmcert.download. Your CSR should now
be signed. Check your email for next steps. Then use the -decrypt option
to extract the CSR request which will then be uploaded to Apple.
```

As it says (and continuing on to step #6 in the instructions) check your email and download the attached signed and encrypted CSR request. Then decrypt this attached file:

```
$ mdmctl mdmcert.download -decrypt=/Users/admin/Downloads/mdm_signed_request.20180502_222719_332.plist.b64.p7
Successfully able to decrypt the MDM Push Certificate request! Please upload
the file 'mdmcert.download.push.req' to Apple by visiting https://identity.apple.com
Once your Push Certificate is signed by Apple you can download it
and import it into MicroMDM using the `mdmctl mdmcert upload` command
```

That's it! Your work with [mdmcert.download](https://mdmcert.download) is done. But now we need to upload this request to Apple.

Now, upload the file `mdmcert.download.push.req` to [https://identity.apple.com](https://identity.apple.com) by using the green 'Create a Certificate' button at top right. You will then be able to download an actual Push Certificate! The filename will be similar to `MDM_ McMurtrie Consulting LLC_Certificate.pem` (note the accompanying private key will be in `mdmcert.download.push.key`.

Now we need to upload this certificate to MicroMDM for use:

```
$ mdmctl mdmcert upload -cert="/Users/admin/Downloads/MDM_ McMurtrie Consulting LLC_Certificate.pem" -private-key=mdmcert.download.push.key
```

You should now be able to use MDM with push certificates!

_Note: you can provide passwords for any of the PKI exchange or Push certificate private keys if you like. Use the `-h` option to list all the options for `mdmctl mdmcert.download`_.

# Using the legacy `certhelper` method

First, sign up for https://mdmcert.download and verify your email. Once you're linked to the instructions page, you can follow the steps below.
Make sure you use the same email address in the steps below as you used to sign up.

Step 0. Make a folder to run the rest of the steps in.
```
mkdir -p mdmcert
cd mdmcert
```

Step 1. Generate a certificate that will be used to decrypt the certificate payload from mdmcert.download
```
openssl genrsa -out server.key 2048
openssl rsa -in server.key -out server.key
openssl req -sha256 -new -key server.key -out server.csr -subj "/CN=micromdm.mdmcert.download"
openssl x509 -req -sha256 -days 365 -in server.csr -signkey server.key -out server.crt

```

The above steps will create the following files:

```
.
├── server.crt
├── server.csr
└── server.key
```

Step 2. Create a Certificate Signing Request
```
certhelper provider -csr -cn=mdm-certtool -password=secret -country=US -email=your_email@acme.co
```

certhelper can be found [here](https://github.com/micromdm/tools/releases). Please use v1.2.0 or higher.

Now your folder will include the following files:

```
├── ProviderPrivateKey.key
├── ProviderUnsignedPushCertificateRequest.csr
├── server.crt
├── server.csr
└── server.key
```

Step 3. Create a request for mdmcert.download
```
certhelper mdmcert.download -cert ./server.crt -csr=ProviderUnsignedPushCertificateRequest.csr -email=your_email@acme.co
```

Wait for a file to be sent to your email. Might take a few minutes or up to a day. 

Step 4. Download and decrypt the file. 
```
certhelper mdmcert.download -cert ./server.crt -key ./server.key -decode ./mdm_signed_request.20171122_094910_220.plist.b64.p7
```

The decode step will create a file called `mdmcert.download_PushCertificateRequest`. Use that in step 5.

Step 5. Upload your `mdmcert.download_PushCertificateRequest` file to [https://identity.apple.com](https://identity.apple.com). You will be able to download a certificate that looks like "MDM_ McMurtrie Consulting LLC_Certificate.pem".

The `ProviderPrivateKey.key` file will be your private key for the push certificate. 

Keep all the files you've created today safely, to start micromdm you'll need:
- The certificate from identity.apple.com
- The private key, `ProviderPrivateKey.key`
- The password for the private key, the password you specified in step 2.

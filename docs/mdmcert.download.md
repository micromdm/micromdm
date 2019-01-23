# Using the `mdmctl` method

Follow the [mdmcert.download instructions](https://mdmcert.download/instructions) until you reach step #4 (using your open source MDM software, in this case MicroMDM). Basically sign up and verify your email.

Then, create a new request with `mdmctl`:

```shell
$ mdmctl mdmcert.download -new -email=cool.mdm.admin@example.org
Request successfully sent to mdmcert.download. Your CSR should now
be signed. Check your email for next steps. Then use the -decrypt option
to extract the CSR request which will then be uploaded to Apple.
```

As it says (and continuing on to step #6 in the instructions) check your email and download the attached signed and encrypted CSR request. Then decrypt this attached file:

```shell
$ mdmctl mdmcert.download -decrypt=/Users/admin/Downloads/mdm_signed_request.20180502_222719_332.plist.b64.p7
Successfully able to decrypt the MDM Push Certificate request! Please upload
the file 'mdmcert.download.push.req' to Apple by visiting https://identity.apple.com
Once your Push Certificate is signed by Apple you can download it
and import it into MicroMDM using the `mdmctl mdmcert upload` command
```

That's it! Your work with [mdmcert.download](https://mdmcert.download) is done. But now we need to upload this request to Apple.

Now, upload the file `mdmcert.download.push.req` to [https://identity.apple.com](https://identity.apple.com) by using the green 'Create a Certificate' button at top right. You will then be able to download an actual Push Certificate! The filename will be similar to `MDM_ McMurtrie Consulting LLC_Certificate.pem` (note the accompanying private key will be in `mdmcert.download.push.key`.

Now we need to upload this certificate to MicroMDM for use:

```shell
$ mdmctl mdmcert upload \
    -cert="/Users/admin/Downloads/MDM_ McMurtrie Consulting LLC_Certificate.pem" \
    -private-key=mdmcert.download.push.key
```

You should now be able to use MDM with push certificates!

_Note: you can provide passwords for any of the PKI exchange or Push certificate private keys if you like. Use the `-h` option to list all the options for `mdmctl mdmcert.download`_.

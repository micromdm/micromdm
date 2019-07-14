# Overview

The MicroMDM server does not composite or sign profiles for you. Instead, it gives you the option of providing a profile which you have signed through an external process, before uploading it to the server. 
To make signing easier, MicroMDM adds the option to sign the profile as part of the `mdmctl` tool. 
First, you will need a private key and certificate. Any certificate will do, but to be verified on the device, you should choose one that is trusted on the devices where you're sending the payloads. The `Developer ID` certificates from Apple are examples of trusted certificates. 

# Signing with mdmctl

The `mdmctl apply profiles` command is used to upload a new configuration profile to the server. It can also be used to sign a profile. 
First, run `mdmctl apply profiles -help` to see all the available command line flags. 

Example command to sign a profile before uploading it to the server:
```
    mdmctl apply profiles \
        -f /path/to/profile.mobileconfig \
        -private-key /path/to/key.pem \
        -cert /path/to/certificate.pem \
        -password=private_key_password \
        -sign 
```

You can also specify the `-out /path/to/signed_output.mobileconfig` to save the signed output locally, instead of uploading it to the server. 
Using `-out -` will print the signed contents to the standard output, allowing you to pipe the output to another operation. 

# Signing with other tools

You can use the `security` command on the Mac to sign configuration profiles with a certificate stored in the Keychain. 
Example:
```
usr/bin/security cms -S -N 'Your KeyChain Cert' -i profile.mobileconfig -o signed.mobileconfig
```

Another popular choice is Jeremy Agostino's [Hancock](https://github.com/JeremyAgost/Hancock), a GUI utility for signing profiles and packages. 

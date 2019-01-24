# Webhooks

To connect MicroMDM with DEP, you'll need to create a new virtual server in the DEP portal, and import the necessary credentials to micromdm:

## Get public key

First, extract the public key which you will need to upload when creating a new virtual server:

```shell
mdmctl get dep-tokens -export-public-key /tmp/DEP_Public_key.pem
ConsumerKey                               AccessTokenExpiry

Wrote DEP public key to: /tmp/DEP_Public_key.pem
```

This assumes you've already configured your mdmctl client to talk to the micromdm server.

## Create server in DEP portal and get P7M file

Next, Go to the DEP portal and create a "MDM Server" (as its called in the DEP poral) record. You will be prompted to upload the public key. Once you successfully create the server, you'll be offered a new download - a file ending `.p7m`. This files contains the encrypted DEP tokens.

## Upload P7M file in micromdm

Next, You'll have to import the downloaded `p7m` file somewhere that MicroMDM can access.

```shell
mdmctl apply dep-tokens -import /Users/groob/Downloads/my_dep_token.p7m
imported DEP token
```

If you see a `imported DEP token` message, your MDM server will now be able to sync devices from DEP.  You can check to see that the token imported correctly by just listing the DEP tokens:

```shell
mdmctl get dep-tokens
ConsumerKey                               AccessTokenExpiry
CK_4e84a1be9638d4bce37443ff2dad55abcacb…  2018-04-22 07:54:12 +0000 UTC
```

## Exporting JSON token

If you need to interact with DEP directly, there is an option to export the OAUTH token used by the server:

```shell
mdmctl get dep-tokens -export-token /path/to/stoken.json
ConsumerKey                               AccessTokenExpiry
CK_4e84a1be9638d4bce37443ff2dad55abcacb…  2018-04-22 07:54:12 +0000 UTC

Wrote DEP token JSON to: /path/to/stoken.json
```

## Firewall

For MDM communication you need to open the following TCP ports

- 80 Redirects to 443
- 443 APN communication backup
- 2195 messages to the APN
- 2196 APN message feedback
- 5223 APN communication

## Other resources

The DEP API is documented by Apple [here](https://developer.apple.com/library/content/documentation/Miscellaneous/Reference/MobileDeviceManagementProtocolRef/4-Profile_Management/ProfileManagement.html#//apple_ref/doc/uid/TP40017387-CH7-SW68) and there are several DEP libraries:

- Go: <https://github.com/micromdm/dep>
- Python: <https://github.com/bruienne/depy>

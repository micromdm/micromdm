MicroMDM uses [Let's Encrypt](https://letsencrypt.org/) to acquire and renew certificates for the MDM server. Doing so reduces a lot of the difficulties with getting started with MDM, like getting a TLS cert or adding anchor certificates for enrollment and DEP. However some users would like to run MicroMDM behind a firewall, or use self-signed certificates. This page will explain how to do so. 

## Creating and using a self-signed certificate

Don't use `localhost` or any IP address for the certificate. Rather you'll need to use the same domain that you specified in the `-server-url` parameter to the server.

To generate your self-signed cert execute this OpenSSL one-liner (replacing the `DNSNAME=`) with your domain name:

```
 DNSNAME=mdm.example.org;  (cat /etc/ssl/openssl.cnf ; printf "\n[SAN]\nsubjectAltName=DNS:$DNSNAME\n") | openssl req -new -newkey rsa:2048 -days 365 -nodes -x509 -sha256 -keyout server.key -out server.crt -subj "/CN=$DNSNAME" -reqexts SAN -extensions SAN -config /dev/stdin
```

The output should look like:

```
Generating a 2048 bit RSA private key
...............+++
.........................+++
writing new private key to 'server.key'
-----
```

In short this generates an X.509 key-pair using a 2048-bit RSA key, SHA-256 signing algorithm as well as generating the SAN extensions for the certificate. If your openssl distribution is older, or has changed the location of the default `openssl.cnf` this command may fail. We hope to integrate this functionality in the future. See the tracking issue [#104](https://github.com/micromdm/micromdm/issues/104).

After this command has run successfully you should have two new files: a `server.key` and a `server.crt` file. These are the PEM-encoded private and public keys.

## Self certs and non-removable DEP profiles
Should you attempt to test self signed certificates and DEP with a non-removable profile, you may end up with the following error after running `mdmctl apply block -udid=4825064D-DE18-4426-AD0C-E49993A1786C`

```
transport=http method=POST status=200 proto=HTTP/1.1 host=::1 user_agent=Go-http-client/1.1 path=/v1/devices/4825064D-DE18-4426-AD0C-E49993A1786C/block
ts="2018/08/06 16:00:06" caller="http: TLS handshake error from [::1" msg="]:57126: remote error: tls: bad certificate"
```

Note this issue (#479) was fixed with PR #480. While you may receive this error, if you use the API tools in this repo (or curl), you can still force the device to unenroll:
```
/tools/api/send_push_notification 4825064D-DE18-4426-AD0C-E49993A1786C

transport=http method=GET status=200 proto=HTTP/2.0 host=::1 user_agent=curl/7.54.0 path=/push/4825064D-DE18-4426-AD0C-E49993A1786C
transport=http method=PUT status=401 proto=HTTP/2.0 host=::1 user_agent="MDM-OSX/1.0 mdmclient/1090" path=/mdm/connect
transport=http method=PUT status=200 proto=HTTP/2.0 host=::1 user_agent="MDM-OSX/1.0 mdmclient/1090" path=/mdm/checkin
```

Be sure to remove the block if you are testing:
`mdmctl remove block -udid=4825064D-DE18-4426-AD0C-E49993A1786C`

## Local testing & DNS

_Note_, for local testing you probably want to add the hostname to your hosts file if you don't have DNS configured:

```
sudo sh -c 'echo "127.0.0.1 mdm.acme.co\n" >> /etc/hosts'
```

## Using the new certificate

Once your certificate and private key are created, you can pass them to `micromdm serve`:

```
# Note the -tls-cert and tls-key flags at the bottom
sudo ./micromdm serve \
	-api-key MySuperSecretKey \
	-apns-cert apns.pem \
	-apns-key apns.key \
	-tls \
	-server-url https://mdm.example.org/ \
	-filerepo ./repo \
	-tls-cert server.crt \
	-tls-key server.key
```

You may also want to enable the "skip-verify" mode in your `mdmctl` configuration, too:

```
mdmctl config set \
	-api-token 12345
	-name MySelfSignedMDMServer
	-server-url https://mdm.example.org/
	-skip-verify true
```

This will disable certificate chain verification on the client side. Alternatively you should be able to add your certificate to your OS's trust store as well.

## Not running as root

The default HTTP ports are 443 and 80 (privileged ports) and the default database location is in `/var/db/micromdm` (also likely privileged). MicroMDM can change the port numbers and database location to be able to run without being root. Combing this with the self-signed certificate invocation looks like this:

```
# Note the -http-addr, -redir-addr, and -config-path flags at the bottom
./micromdm serve \
	-api-key MySuperSecretKey \
	-apns-cert apns.pem \
	-apns-key apns.key \
	-tls \
	-server-url https://mdm.example.org/ \
	-filerepo ./repo \
	-tls-cert server.crt \
	-tls-key server.key \
	-http-addr :8443 \
	-redir-addr :8080 \
	-config-path /path/to/my/config
```

Also note on Linux the capabilities system can give specific binaries listen-only ability:

```
sudo setcap 'cap_net_bind_service=+ep' ./micromdm
```

## DEP anchor certificate

When using self-signed certificates in the MDM server DEP has a mechanism to tell the device to trust the MDM server initially (as straight out of the box/freshly booted the device will not trust your MDM self-signed certificate). The DEP profile has capacity to specify an "anchor" cert.

MicroMDM makes working with these easy when using the `mdmctl apply dep-profiles -f /path/to/dep-profile.json` syntax to apply a DEP profile. It as two additional switches here: `-anchor` and `-use-server-cert`. The anchor switche takes a path to a PEM certificate and properly encodes it into the DEP profile. The use-server-cert switch simply uses the certificate mdmctl uses to connect to MicroMDM to inject into the DEP profile for you. The `-use-server-cert` switch is probably all you need to make quick work of using a self-signed certificate. E.g.:

```
mdmctl apply dep-profiles -use-server-cert -f /path/to/dep-profile.json
```

Or, to explicitly specific an anchor certificate (i.e. in cases of more complicated setups like proxied HTTPS, etc.):

```
mdmctl apply dep-profiles -anchor /path/to/anchor.pem -f /path/to/dep-profile.json
```

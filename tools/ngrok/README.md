# Local Development with ngrok

When developing and testing MicroMDM you will likely want to enroll devices. Normally that means either generating or getting new TLS certificates and hosting micromdm on some publicly accessible server.
`ngrok` is a tool that makes it easy to expose the MicroMDM server as a public URL that your devices can enroll to.

The scripts in this will help you get MicroMDM running for local development.
First, install the listed requirements.

Next, make sure you've cloned the `micromdm` repo exactly as described in the [contributing guide](../../CONTRIBUTING.md). Follow the steps in the same document to downloads required dependencies and build `micromdm` and `mdmctl`.

Finally, follow the described workflow for setting up your development server.

## Requirements

- [go](https://golang.org/)
  `brew install go`

- [jq](https://stedolan.github.io/jq/)
  `brew install jq`

- [ngrok](https://ngrok.com)
  `brew cask install ngrok`

## Workflow

Follow the steps in order, and use separate terminal windows for running ngrok and micromdm. These are long running processes, and you will want to see the output.
Some of the steps are only required for first time setup.

_Note:_ When the ngrok tunnel expires, you have to re-enroll the test device. The `configure_mdmctl` updates the DEP profile, so all you need to do is either reset the VM, or run `sudo /usr/bin/profiles -e -N` if testing with a real Mac.

1. Run ngrok. A tunnel will be opened. If you're using the free version of `ngrok` this URL is only valid for 8 hours, but you can always start with a new URL.
```
ngrok http 8080
```
Make sure to open a browser tab at `http://localhost:4040`. The ngrok web interface shows all the requests/responses between `micromdm` and its clients and can be invaluable for debugging.

2. Create a env file and define environment variables you will need for micromdm. This is a first time setup only step.

Contents of `env` file:
```
export API_TOKEN=supersecret
export PUSH_CERTIFICATE_PRIVATE_KEY_PASSWORD=supersecret
export SERVER_URL="$(curl -s localhost:4040/api/tunnels | jq '.tunnels[]|.public_url' -r |grep https)"
```

3. Set the `MICROMDM_ENV_PATH` variable to be the path of your `env` file. This will be used to source the file in all the helper scripts. Do this once in every new terminal window you open while developing with micromdm.

```
export MICROMDM_ENV_PATH="$(pwd)/env"
```

4. Start `micromdm`. Do this every time you rebuild the binary and when your `ngrok` tunnel expires.

```
./tools/ngrok/start_server
```

The server will run on port `8080`. The default configuration is minimal, but you can pass any additional flags as arguments to the script. For example `./tools/ngrok/start_server -http-debug`.

5. Configure the Push certificate and DEP tokens. This is a first time setup only step.
If you already have the `MDM CSR` option in your enterprise portal, [follow](https://github.com/micromdm/micromdm/blob/main/docs/user-guide/quickstart.md#configure-an-apns-certificate) the user guide steps.  
If you don't have access to the enterprise portal, you can get an MDM certificate from [mdmcert.download](https://github.com/micromdm/micromdm/wiki/mdmcert.download) 
You can also export the device management certificate from [Server.app](https://github.com/micromdm/micromdm/wiki/Export-the-Profile-Manager-Certificate-for-MicroMDM-Testing-and-Development).

Configure `mdmdctl`:

```
source $MICROMDM_ENV_PATH
./build/darwin/mdmctl config set -name ngrok -api-token=$API_TOKEN -server-url=$SERVER_URL
./build/darwin/mdmctl config switch -name=ngrok
```

Upload push certificate:
```
./build/darwin/mdmctl mdmcert upload  \
  -cert "$path_to/mdm-certificates/push_certificate.pem" \
  -private-key "$path_to/mdm-certificates/PushCertificatePrivateKey.key" \
  -password $PUSH_CERTIFICATE_PRIVATE_KEY_PASSWORD
```

Get a DEP OAuth token at deploy.apple.com
```
./build/darwin/mdmctl get dep-tokens -export-public-key /path/to/save/public.pem

# after you upload the public key to your DEP portal, you'll get back a .p7m token file.
./build/darwin/mdmctl apply dep-tokens -import "$dev_path/dep/local_dep_smime.p7m"
```

6. Configure mdmctl and start using micromdm. Do this every time your ngrok tunnel expires.
```
# substitute A_SERIAL_NUMBER for your test device. This serial number will be used to apply a DEP profile.
./tools/ngrok/configure_mdmctl A_SERIAL_NUMBER
```

`micromdm` and `mdmctl` should now be configured and you can begin testing.

## Everyday workflow.

```
# In your first terminal window. Don't forget to open http://localhost:4040
ngrok http 8080

# In a second terminal window
export MICROMDM_ENV_PATH="$(pwd)/env"
./tools/ngrok/start_server

# In a third terminal window.
export MICROMDM_ENV_PATH="$(pwd)/env"
./tools/ngrok/configure_mdmctl
```

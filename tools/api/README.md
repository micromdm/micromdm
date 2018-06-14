A sample of the common API requests are available as `curl` scripts. To use them, install all the requirements and follow the steps as described.

## Requirements

- [jq](https://stedolan.github.io/jq/)
  `brew install jq`

## Setup
Create an `env` file and define environment variables you will need to talk to the server.
This env file will be sourced by the scripts.

Contents of `env` file:
```
# the value of the -api-key flag that MicroMDM was started with.
export API_TOKEN=supersecret
export SERVER_URL=https://mdm.acme.co
```

In your shell, set the environment variable `MICROMDM_ENV_PATH` to point to your env file.
Do this every time you open a new shell to work with the scripts in this folder.
```
export MICROMDM_ENV_PATH="$(pwd)/env"
```

## Usage examples

```
# get all devices from micromdm. returns JSON.
./tools/api/get_devices

# use jq to filter response. For example, to get the udid of the first device.
./tools/api/get_devices | jq .devices[0].udid -r

# send a push notification to a device UDID
./tools/api/send_push_notification <device-udid>

# combine sending a push notification with the get devices request.
$udid=(tools/api/get_devices |jq .devices[0].udid -r)
./tools/api/send_push_notification $udid
```

The `api` folder groups all MDM commands in the `commands` folder. Use these for scheduling commands against a remote device. All of these take the device udid as a first argument and might take additional arguments when necessary.

```
# install a configuration profile
./tools/api/commands/install_profile $udid /path/to/profile.mobileconfig

# remove a configuration profile
./tools/api/commands/remove_profile $udid Your-Profile-PayloadIdentifier
```

All the scripts in this folder are tiny. Consider reading them to get an understanding of the API.
Also consider learning `jq`.
MicroMDM returns and accepts JSON for all of the API endpoints, and `jq` is a powerful tool for reading and manipulating JSON data.

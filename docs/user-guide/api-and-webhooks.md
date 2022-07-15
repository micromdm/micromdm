MicroMDM is intended to be a dependency for implementing a higher level workflow. To do so, it exposes all its capabilities through an HTTP API and also forwards all the device responses to an HTTP endpoint of your choice. 

For examples of what you can do with the API, reference the [tools/api](https://github.com/micromdm/micromdm/tree/main/tools/api) directory of the micromdm repository. The directory includes executable `curl` examples of API requests which you can run, but also reference when implementing an API client in another language.

Throughout this document we will reference the provided scripts as shorthand for making API requests.

# Configure a webhook URL to process device events 

When you schedule a command on the device, it will respond. As a rule of thumb, MicroMDM does not process the responses from the devices outside of maintaining the command queue and storing tokens/certificates it needs to keep communication with the devices. Instead, it allows the configuration of an HTTP Webhook URL, which gives you the ability to process the responses how you wish. 

To configure the webhook URL, start `micromdm` with the `-command-webhook-url` flag. 

## Events

Each event sent to the webhook url contains a json object in the body of the request which represents the event.

| Property          | Description                                      |
|-------------------|--------------------------------------------------|
| topic             | The type of MicroMDM event. See values below.    |
| event_id          | A unique id representing the event.              |
| created_at        | The timestamp that MicroMDM generated the event. |
| checkin_event     | Optional payload based on the topic.             |
| acknowledge_event | Optional payload based on the topic.             |


The following MicroMDM Topics are exposed via the webhook functionality:

| Topic                             | Payload                                  |
|-----------------------------------|------------------------------------------|
| [mdm.Authenticate](#authenticate) | [checkin_event](#checkin-events)         |
| [mdm.TokenUpdate](#token-update)  | [checkin_event](#checkin-events)         |
| [mdm.CheckOut](#checkout)         | [checkin_event](#checkin-events)         |
| [mdm.Connect](#connect)           | [acknowledge_event](#acknowledge-events) |


The following is an example of the json payload in the body of the request.

```json
{
    "topic": "mdm.Connect",
    "event_id": "09fd3348-24c8-43cf-91d7-1a5d5d599012",
    "created_at": "2018-08-13T14:30:20.491166786Z",
    "checkin_event": {
        ...
        "raw_payload": ""
    },
    "acknowledge_event": {
        ...
        "raw_payload": ""
    }
}
```

Depending on which topic the event is for one of the optional `_event` payloads will be filled in. Each `_event` payload contains a `raw_payload` field which is a base64 encoded representation of the actual data the device sent back to the MDM server. For more detailed information on these raw payloads please refer to the official Apple MDM Protocol Reference.

https://developer.apple.com/enterprise/documentation/MDM-Protocol-Reference.pdf

### Checkin Events

The MDM check-in protocol is used during initialization to validate a deviceʼs eligibility for MDM enrollment and to inform the server that a deviceʼs push token has been updated.

| Property    | Description                                                                                               |
|-------------|-----------------------------------------------------------------------------------------------------------|
| udid        | UDID of the device.                                                                                       |
| url_params  | Any additional http params added in the enrollment profile `CheckInURL` field are passed through to here. |
| raw_payload | The raw data sent from the device to MicroMDM.                                                            |

#### Authenticate

While the user is installing an MDM payload, the device sends an authenticate message.

```json
{
    "topic": "mdm.Authenticate",
    "event_id": "5ee634b5-bd58-4b3e-8b77-3f764d087787",
    "created_at": "2018-08-13T14:30:16.13272395Z",
    "checkin_event": {
        "udid": "A5EF1BA1-586D-4F29-B4F3-759DADAC2DDD",
        "url_params": null,
        "raw_payload": "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4KPCFET0NUWVBFIHBsaXN0IFBVQkxJQyAiLS8vQXBwbGUvL0RURCBQTElTVCAxLjAvL0VOIiAiaHR0cDovL3d3dy5hcHBsZS5jb20vRFREcy9Qcm9wZXJ0eUxpc3QtMS4wLmR0ZCI+CjxwbGlzdCB2ZXJzaW9uPSIxLjAiPgo8ZGljdD4KCTxrZXk+QnVpbGRWZXJzaW9uPC9rZXk+Cgk8c3RyaW5nPjE3RzY1PC9zdHJpbmc+Cgk8a2V5PkNoYWxsZW5nZTwva2V5PgoJPGRhdGE+CglZWEJ3YkdVPQoJPC9kYXRhPgoJPGtleT5EZXZpY2VOYW1lPC9rZXk+Cgk8c3RyaW5nPkhvc3ROYW1lPC9zdHJpbmc+Cgk8a2V5Pk1lc3NhZ2VUeXBlPC9rZXk+Cgk8c3RyaW5nPkF1dGhlbnRpY2F0ZTwvc3RyaW5nPgoJPGtleT5Nb2RlbDwva2V5PgoJPHN0cmluZz5NYWNCb29rUHJvMTQsMzwvc3RyaW5nPgoJPGtleT5Nb2RlbE5hbWU8L2tleT4KCTxzdHJpbmc+TWFjQm9vayBQcm88L3N0cmluZz4KCTxrZXk+T1NWZXJzaW9uPC9rZXk+Cgk8c3RyaW5nPjEwLjEzLjY8L3N0cmluZz4KCTxrZXk+UHJvZHVjdE5hbWU8L2tleT4KCTxzdHJpbmc+TWFjQm9va1BybzE0LDM8L3N0cmluZz4KCTxrZXk+U2VyaWFsTnVtYmVyPC9rZXk+Cgk8c3RyaW5nPkMwNFcwMEQzSFREMDwvc3RyaW5nPgoJPGtleT5Ub3BpYzwva2V5PgoJPHN0cmluZz5jb20uYXBwbGUubWdtdC5FeHRlcm5hbC4zYTdjOTVmMi1jMjE0LTRiNTMtYTU4NS1hNTUzYzI2MjNlMDQ8L3N0cmluZz4KCTxrZXk+VURJRDwva2V5PgoJPHN0cmluZz5BNUVGMUJBMS01ODZELTRGMjktQjRGMy03NTlEQURBQzJEREQ8L3N0cmluZz4KPC9kaWN0Pgo8L3BsaXN0Pgo="
    }
}
```

#### Token Update

A device sends a token update message to the check-in server whenever its device push token, push magic, or
unlock token change. These fields are needed by the server to send the device push notifications or passcode
resets. The device sends an initial token update message to the server when it has installed the MDM payload. The server
should send push messages to the device **only** after receiving the first token update message. 

```json
{
    "topic": "mdm.TokenUpdate",
    "event_id": "99eb840c-3e1f-48a7-a4f9-1ef092445129",
    "created_at": "2018-08-13T14:30:16.529223669Z",
    "checkin_event": {
        "udid": "A5EF1BA1-586D-4F29-B4F3-759DADAC2DDD",
        "url_params": null,
        "raw_payload": "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4KPCFET0NUWVBFIHBsaXN0IFBVQkxJQyAiLS8vQXBwbGUvL0RURCBQTElTVCAxLjAvL0VOIiAiaHR0cDovL3d3dy5hcHBsZS5jb20vRFREcy9Qcm9wZXJ0eUxpc3QtMS4wLmR0ZCI+CjxwbGlzdCB2ZXJzaW9uPSIxLjAiPgo8ZGljdD4KCTxrZXk+QXdhaXRpbmdDb25maWd1cmF0aW9uPC9rZXk+Cgk8ZmFsc2UvPgoJPGtleT5NZXNzYWdlVHlwZTwva2V5PgoJPHN0cmluZz5Ub2tlblVwZGF0ZTwvc3RyaW5nPgoJPGtleT5QdXNoTWFnaWM8L2tleT4KCTxzdHJpbmc+NzI4Rjg2NzAtNDkzNC00QTMwLUJBREYtN0I0M0MyNjFGQzE0PC9zdHJpbmc+Cgk8a2V5PlRva2VuPC9rZXk+Cgk8ZGF0YT4KCWMzRnlIV2hqRzN0bVNBNTVMRkZJUWc9PQoJPC9kYXRhPgoJPGtleT5Ub3BpYzwva2V5PgoJPHN0cmluZz5jb20uYXBwbGUubWdtdC5FeHRlcm5hbC4zYTdjOTVmMi1jMjE0LTRiNTMtYTU4NS1hNTUzYzI2MjNlMDQ8L3N0cmluZz4KCTxrZXk+VURJRDwva2V5PgoJPHN0cmluZz5BNUVGMUJBMS01ODZELTRGMjktQjRGMy03NTlEQURBQzJEREQ8L3N0cmluZz4KPC9kaWN0Pgo8L3BsaXN0Pgo="
    }
}
```

#### CheckOut

In iOS 5.0 and later, and in macOS v10.9, if the CheckOutWhenRemoved key in the MDM enrollment profile is set to true, the device attempts to send a CheckOut message when the MDM profile is removed.

```json
{
    "topic": "mdm.CheckOut",
    "event_id": "3f3a0620-6c3a-4dc5-bdc6-523f071e5865",
    "created_at": "2018-08-13T18:36:12.727463835Z",
    "checkin_event": {
        "udid": "A5EF1BA1-586D-4F29-B4F3-759DADAC2DDD",
        "url_params": null,
        "raw_payload": "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4KPCFET0NUWVBFIHBsaXN0IFBVQkxJQyAiLS8vQXBwbGUvL0RURCBQTElTVCAxLjAvL0VOIiAiaHR0cDovL3d3dy5hcHBsZS5jb20vRFREcy9Qcm9wZXJ0eUxpc3QtMS4wLmR0ZCI+CjxwbGlzdCB2ZXJzaW9uPSIxLjAiPgo8ZGljdD4KCTxrZXk+TWVzc2FnZVR5cGU8L2tleT4KCTxzdHJpbmc+Q2hlY2tPdXQ8L3N0cmluZz4KCTxrZXk+VG9waWM8L2tleT4KCTxzdHJpbmc+Y29tLmFwcGxlLm1nbXQuRXh0ZXJuYWwuM2E3Yzk1ZjItYzIxNC00YjUzLWE1ODUtYTU1M2MyNjIzZTA0PC9zdHJpbmc+Cgk8a2V5PlVESUQ8L2tleT4KCTxzdHJpbmc+QTVFRjFCQTEtNTg2RC00RjI5LUI0RjMtNzU5REFEQUMyREREPC9zdHJpbmc+CjwvZGljdD4KPC9wbGlzdD4K"
    }
}
```

### Acknowledge Events

Acknowledge events are sent when the device responds to a command sent to it.

| Property     | Description                                                                                              |
|--------------|----------------------------------------------------------------------------------------------------------|
| udid         | UDID of the device.                                                                                      |
| status       | Status. See the [Apple MDM Protocol Reference](https://developer.apple.com/enterprise/documentation/MDM-Protocol-Reference.pdf) for values. |
| command_uuid | UUID of the command that this response is for (if any).                                                  |
| url_params   | Any additional http params added in the enrollment profile `ServerURL` field are passed through to here. |
| raw_payload  | The raw data sent from the device to MicroMDM.                                                           |

#### Connect

After the MDM server has sent a command to the device, the device replies to the server by sending a plist-encoded dictionary containing the response to the command.

```json
{
    "topic": "mdm.Connect",
    "event_id": "3444be1d-bbf0-45dc-9da3-f707beeecf1b",
    "created_at": "2018-08-13T14:30:16.789541405Z",
    "acknowledge_event": {
        "udid": "A5EF1BA1-586D-4F29-B4F3-759DADAC2DDD",
        "status": "Acknowledged",
        "command_uuid": "41d35de3-a343-4146-ba4b-0069bae2a54f",
        "url_params": null,
        "raw_payload": "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4KPCFET0NUWVBFIHBsaXN0IFBVQkxJQyAiLS8vQXBwbGUvL0RURCBQTElTVCAxLjAvL0VOIiAiaHR0cDovL3d3dy5hcHBsZS5jb20vRFREcy9Qcm9wZXJ0eUxpc3QtMS4wLmR0ZCI+CjxwbGlzdCB2ZXJzaW9uPSIxLjAiPgo8ZGljdD4KCTxrZXk+Q29tbWFuZFVVSUQ8L2tleT4KCTxzdHJpbmc+NDFkMzVkZTMtYTM0My00MTQ2LWJhNGItMDA2OWJhZTJhNTRmPC9zdHJpbmc+Cgk8a2V5PlN0YXR1czwva2V5PgoJPHN0cmluZz5BY2tub3dsZWRnZWQ8L3N0cmluZz4KCTxrZXk+VURJRDwva2V5PgoJPHN0cmluZz5BNUVGMUJBMS01ODZELTRGMjktQjRGMy03NTlEQURBQzJEREQ8L3N0cmluZz4KPC9kaWN0Pgo8L3BsaXN0Pgo="
    }
}
```

## Example Code

Creating a simple webhook listener is as simple as listening for the POST requests from MicroMDM. Below is an example of a python [Flask](http://flask.pocoo.org/) server that just prints out all the messages it receives.

```python
from flask import Flask, request, abort

app = Flask(__name__)

@app.route('/webhook', methods=['POST'])
def webhook():
    print(request.json)
    return ''

if __name__ == '__main__':
    app.run()
```

For examples in more languages please see the [micromdm-webhook-blueprints](https://github.com/knightsc/micromdm-webhook-blueprints) project.


# Schedule Commands with the API

MicroMDM uses Basic Auth for its rest API, requiring you to provide `micromdm` as the user, and the `-api-key` value as the password. 

Let's take a look at the `./tools/api/send_push_notification` script as an example of how to use the API.

```
#!/bin/bash

source $MICROMDM_ENV_PATH
endpoint="push"
curl $CURL_OPTS -u "micromdm:$API_TOKEN" "$SERVER_URL/$endpoint/$1"
```

This `curl` request will use the `-u` flag for authentication, and make a request to `https://your-mdm/push/your-device-udid`. Once received, MicroMDM will send an MDM push notification to that device(or user) UDID. 

Assuming the device is online and able to respond, it will contact the `/mdm/connect` endpoint with a request like so:

```
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>Status</key>
	<string>Idle</string>
	<key>UDID</key>
	<string>55693EB3-DF03-5FD1-9263-F7CDB8AD7FFD</string>
</dict>
</plist>
```

This request will be processed by MicroMDM and also delivered to the webhook URL under the `raw_payload` property. 


An important endpoint in MicroMDM's API is the `/v1/commands` endpoint, which you can use to schedule MDM commands. 
This commands endpoint accepts a POST request, with a JSON payload that closely matches the [documented](https://developer.apple.com/documentation/devicemanagement/mdm_commands_and_queries?changes=latest_minor) device management specification by Apple. Knowing the request type you want to use, you can compose a command to this endpoint. 

For example, to install a configuration profile to a device, you can schedule the following command:
```
POST /v1/commands HTTP/1.1
Authorization: Basic bWljcm9tZG06c3VwZXJzZWNyZXQ=
{  
    "udid": "55693EB3-DF03-5FD1-9263-F7CDB8AD7FFD", 
    "request_type": "InstallProfile",
    "payload": "PD94bWwgdmVyc2lvbj0iMS4...base64_encoded_mobileconfig",  
}
```

MicroMDM will convert this request into a complete command, and schedule it on the queue. It will then send a push notification to ask the device to check in, and respond with the InstallProfile command. 

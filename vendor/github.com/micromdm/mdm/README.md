[![GoDoc](https://godoc.org/github.com/micromdm/mdm?status.svg)](http://godoc.org/github.com/micromdm/mdm) [![Build Status](https://travis-ci.org/micromdm/mdm.svg?branch=master)](https://travis-ci.org/micromdm/mdm)

The mdm package holds structs and helper methods for payloads in Apple's Mobile Device Management protocol.  
This package embeds the various payloads and responses in two structs - `Payload` and `Response`.

# How an MDM server executes commands on a device.
To communicate with a device, an MDM server must create a Payload property list with a specific RequestType and additional data for each request type. Let's use the DeviceInformation request as an example:


```
    // create a request
	request := &CommandRequest{
		RequestType: "DeviceInformation",
		Queries:     []string{"IsCloudBackupEnabled", "BatteryLevel"},
	}

    // NewPayload will create a proper Payload based on the CommandRequest struct
	payload, err := NewPayload(request)
	if err != nil {
		log.Fatal(err)
	}

	// Encode in a plist and print to stdout
    // uses the github.com/groob/plist package
	encoder := plist.NewEncoder(os.Stdout)
	encoder.Indent("  ")
	if err := encoder.Encode(payload); err != nil {
		log.Fatal(err)
	}
```

Resulting command payload:
```
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
  <dict>
    <key>Command</key>
    <dict>
      <key>Queries</key>
      <array>
        <string>IsCloudBackupEnabled</string>
        <string>BatteryLevel</string>
      </array>
      <key>RequestType</key>
      <string>DeviceInformation</string>
    </dict>
    <key>CommandUUID</key>
    <string>fa34b4b7-0553-4b3a-9c4b-76b8b357a622</string>
  </dict>
</plist>
```

An MDM server will queue this request and send a push notification to a device. When device checks in, the server will
reply with the queued plist.

Once the device receives and processes the payload plist, it will reply back to the server. The response will be another plist, which can be unmarshalled into the `Response` struct. Below is the response to our DeviceInformation request.

```
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>CommandUUID</key>
    <string>fa34b4b7-0553-4b3a-9c4b-76b8b357a622</string>
	<key>QueryResponses</key>
	<dict>
		<key>BatteryLevel</key>
		<real>1</real>
		<key>IsCloudBackupEnabled</key>
		<false/>
	</dict>
	<key>Status</key>
	<string>Acknowledged</string>
	<key>UDID</key>
	<string>1111111111111111111111111111111111111111</string>
</dict>
</plist>
```

package command_test

import (
	"reflect"
	"testing"

	"github.com/micromdm/micromdm/platform/command"
)

const testRawCmd = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Command</key>
    <dict>
        <key>ManagedOnly</key>
        <true/>
        <key>RequestType</key>
        <string>ProfileList</string>
    </dict>
    <key>CommandUUID</key>
    <string>0001_ProfileList</string>
</dict>
</plist>`

func TestRawEvent(t *testing.T) {
	cmd := &command.RawCommand{
		UDID:        "1234",
		CommandUUID: "0001_ProfileList",
		Raw:         []byte(testRawCmd),
	}
	cmd.Command.RequestType = "ProfileList"

	ev := command.NewRawEvent(cmd)

	buf, err := command.MarshalRawEvent(ev)
	if err != nil {
		t.Fatalf("could not marshal event: %v", err)
	}

	ev2 := new(command.RawEvent)
	if err = command.UnmarshalRawEvent(buf, ev2); err != nil {
		t.Fatalf("could not unmarshal event: %v", err)
	}

	if !reflect.DeepEqual(ev, ev2) {
		t.Error("expected events to be equal")
	}
}

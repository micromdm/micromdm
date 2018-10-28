package webhook

import "text/template"

var tmplStr = template.Must(template.New("scep_profile").Parse(stub))

const stub = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
    <dict>
        <key>PayloadContent</key>
        <array>
            <dict>
                <key>PayloadContent</key>
                <dict>
                    <key>Key Type</key>
                    <string>RSA</string>
                    <key>Keysize</key>
                    <integer>1024</integer>
                    <key>Retries</key>
                    <integer>3</integer>
                    <key>RetryDelay</key>
                    <integer>10</integer>
                    <key>URL</key>
					<string>https://2d8517bb.ngrok.io/scep</string>
                    <key>Subject</key>
                    <array>
                        <array>
                            <array>
                                <string>CN</string>
                				<string>io.micromdm.workflow.profile.scep.{{.UDID}}</string>
                            </array>
                        </array>
                    </array>
                </dict>
                <key>PayloadDescription</key>
                <string>Configures SCEP settings</string>
                <key>PayloadDisplayName</key>
                <string>SCEP</string>
                <key>PayloadIdentifier</key>
                <string>io.micromdm.workflow.profile.{{.UDID}}</string>
                <key>PayloadType</key>
                <string>com.apple.security.scep</string>
                <key>PayloadUUID</key>
                <string>{{.SCEPPayloadUUID}}</string>
                <key>PayloadVersion</key>
                <real>1</real>
            </dict>
            <dict>
                <key>PayloadCertificateFileName</key>
                <string>ca.crt</string>
                <key>PayloadContent</key>
                <data>{{.CAPEM}}</data>
                <key>PayloadDescription</key>
                <string>Configures certificate settings.</string>
                <key>PayloadDisplayName</key>
                <string>ca.crt</string>
                <key>PayloadIdentifier</key>
        		<string>io.micromdm.workflow.profile.scep_ca_root.{{.UDID}}</string>
                <key>PayloadType</key>
                <string>com.apple.security.root</string>
                <key>PayloadUUID</key>
                <string>CACertPayloadUUID</string>
                <key>PayloadVersion</key>
                <integer>1</integer>
            </dict>
        </array>
        <key>PayloadDisplayName</key>
        <string>Profile Workflow Identity</string>
        <key>PayloadIdentifier</key>
        <string>io.micromdm.workflow.profile.setup_svc.{{.UDID}}</string>
        <key>PayloadRemovalDisallowed</key>
        <false/>
        <key>PayloadType</key>
        <string>Configuration</string>
        <key>PayloadUUID</key>
        <string>{{.ProfileUUID}}</string>
        <key>PayloadVersion</key>
        <integer>1</integer>
    </dict>
</plist>`

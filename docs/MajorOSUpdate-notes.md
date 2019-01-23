Scheduling an OSUpdateScan I got this:

```
    <key>AvailableOSUpdates</key>
    <array>
        <dict>
            <key>AllowsInstallLater</key>
            <true/>
            <key>AppIdentifiersToClose</key>
            <array/>
            <key>HumanReadableName</key>
            <string>macOS Installer Notification</string>
            <key>HumanReadableNameLocale</key>
            <string>en</string>
            <key>IsConfigDataUpdate</key>
            <true/>
            <key>IsCritical</key>
            <false/>
            <key>IsFirmwareUpdate</key>
            <false/>
            <key>MetadataURL</key>
            <string>http://swcdn.apple.com/content/downloads/50/07/041-14451/vanwrs4pcxauvs2mzxztqwihcfgmkqaule/macOSInstallerNotification_GM.smd</string>
            <key>ProductKey</key>
            <string>041-14451</string>
            <key>RestartRequired</key>
            <false/>
            <key>Version</key>
            <string>2.0</string>
        </dict>
```

Then, scheduling an OSUpdate for `041-14451` I got the following from an AvailableOSUpdates command:

```
        <dict>
            <key>DownloadSize</key>
            <real>12000000000</real>
            <key>HumanReadableName</key>
            <string>macOS Mojave</string>
            <key>HumanReadableNameLocale</key>
            <string>en</string>
            <key>IsConfigDataUpdate</key>
            <false/>
            <key>IsCritical</key>
            <false/>
            <key>IsFirmwareUpdate</key>
            <false/>
            <key>IsMajorOSUpdate</key>
            <true/>
            <key>ProductKey</key>
            <string>_OSX_18A391</string>
            <key>RestartRequired</key>
            <true/>
            <key>Version</key>
            <string>18A391</string>
        </dict>
    </array>
```


And sure enough, the sucatalog for 10.12 and 10.13 has the following added (as of today)

```
            <key>041-14451</key>
            <dict>
                <key>ServerMetadataURL</key>
                <string>http://swcdn.apple.com/content/downloads/50/07/041-14451/vanwrs4pcxauvs2mzxztqwihcfgmkqaule/macOSInstallerNotification_GM.smd</string>
                <key>Packages</key>
                <array>
                    <dict>
                        <key>Digest</key>
                        <string>f11446026976c362a451b89d354c2313af19a1a3</string>
                        <key>Size</key>
                        <integer>1821971</integer>
                        <key>MetadataURL</key>
                        <string>https://swdist.apple.com/content/downloads/50/07/041-14451/vanwrs4pcxauvs2mzxztqwihcfgmkqaule/macOSInstallerNotification_GM.pkm</string>
                        <key>URL</key>
                        <string>http://swcdn.apple.com/content/downloads/50/07/041-14451/vanwrs4pcxauvs2mzxztqwihcfgmkqaule/macOSInstallerNotification_GM.pkg</string>
                    </dict>
                </array>
                <key>PostDate</key>
                <date>2018-10-16T20:05:00Z</date>
```

Downloading the notification package shows a `OSXNotification.bundle` payload with the following info.plist:

```
    ```<key>CFBundleIdentifier</key>
    <string>com.apple.installer.notification.macOS1014GM</string>
    <key>CFBundleInfoDictionaryVersion</key>
    <string>6.0</string>
    <key>CFBundleName</key>
    <string>OSXNotification</string>
    <key>CFBundlePackageType</key>
    <string>BNDL</string>
    <key>CFBundleShortVersionString</key>
    <string>2</string>
    <key>CFBundleSupportedPlatforms</key>
    <array>
        <string>MacOSX</string>
    </array>
    <key>CFBundleVersion</key>
    <string>2</string>
    <key>FreeSpaceRequired</key>
    <integer>12000000000</integer>
    <key>HumanReadableName</key>
    <string>macOS Mojave</string>
    <key>ItemID</key>
    <integer>1398502828</integer>
    <key>NSHumanReadableCopyright</key>
    <string>Copyright © 2017 Apple. All rights reserved.</string>
    <key>ProductBuildVersion</key>
    <string>18A391</string>```
```

When scheduling the `_OSX_18A391` product key, the NotifyOnly and DownloadOnly actions trigger error responses. `InstallASAP` and both are supported, but the status remains as `Idle`. Plugging the computer into power and _then_ scheduling the update caused `softwareupdated` to actually start the download. 
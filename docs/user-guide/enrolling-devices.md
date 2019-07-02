To enroll into MDM, a device needs to install an MDM Configuration Profile. There are many ways of doing that, you could even email one to your users and have them install it manually. 

MicroMDM has a built-in option for enrollment profiles and some helpers to tweak defaults. If the built-in options don't help, a separate enrollment service can be used. 

# Default Enrollment Method

MicroMDM offers an enrollment handler at its `/mdm/enroll` URL. This endpoint is used both for user and DEP enrollments and is the same across all your devices. The default enrollment profile includes all permissions and is also unsigned. 

# DEP Device Assignment
If you're enrolled in Apple Business Manager, you can configure all your devices to use MicroMDM during provisioning.

Get a DEP profile template and fill it out.
All the keys you can use [are documented](https://developer.apple.com/documentation/devicemanagement/profile?changes=latest_minor) by Apple.
```
mdmctl apply dep-profiles -template > /tmp/profile.json
```

Use
```
mdmctl apply dep-profiles -f /path//to/profile.json
``` 
to define the profile. 

You can use the `devices` array to assign this profile to specific serial numbers. 
You can also make this profile the default for all devices in your linked DEP server. 

```
mdmctl apply dep-profiles -f /path/to/dep-profile.json -filter='*'
```

For more details on auto-assignment, [check](https://github.com/micromdm/micromdm/wiki/DEP-auto-assignment) the wiki page.

# Replacing the default Enrollment Profile

You might want to customize the enrollment profile offered to your devices. To do so, you can download the default enrollment profile, tweak it, and upload a new one. 

To download and save the enrollment profile, you can use `curl`:

```
curl -o enroll.mobileconfig https://mdm.acme.co/mdm/enroll
```

Next, use a text editor to modify the payload. 
You can also [sign](https://github.com/micromdm/micromdm/wiki/Sign-the-enrollment-profile-with-Hancock) the profile. 

Once modified and signed, you can replace the default. 

```
mdmctl apply profiles -f /path/to/enroll.mobileconfig
```

Now, the profile is still offered at `/mdm/enroll`, but is the customized one.

# OTA Enrollment

For Over-the-Air profile delivery, [check out notes](https://github.com/micromdm/micromdm/wiki/OTA-Enrollment) from the wiki. 

# A custom enrollment endpoint

For some environments, serving the same enrollment profile might not be enough. You can create your own HTTP service which generates the configuration profile, and not use `/mdm/enroll` endpoint. 

How you do that is up to you, you just need to ensure that the profile has the CheckIn and ServerURL values pointing back to the micromdm endpoints.

```
        <key>CheckInURL</key>
        <string>https://mdm.acme.co/mdm/checkin</string>
        <key>ServerURL</key>
        <string>https://mdm.acme.co/mdm/connect</string>
```




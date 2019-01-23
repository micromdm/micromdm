![](https://user-images.githubusercontent.com/1526945/33252691-acfa1952-d30d-11e7-800a-c9a271d697b5.png)

MicroMDM allows an administrator to force the device to trigger an HTTP Status 401 (Unauthorized) response by adding the device UDID to a `block`. The next time the device connects to the server to receive a command, it will be forced to unenroll. 

Add a block:
```
mdmctl apply block -udid=564DF4AE-96B5-72BF-02C4-6B8B94BA5BF4
```
The command will add the UDID to a list of device blocks and also trigger a push notification, causing the device to connect ASAP.

The enrollment status can be checked by running `mdmctl get devices`. 
```
mdmctl get devices

UDID                                      SerialNumber  EnrollmentStatus  LastSeen
564DF4AE-96B5-72BF-02C4-6B8B94BA5BF4      C02T5LL6GZX2  false             2017-11-27 05:31:10.724424502 +0000 UTC
```

The device will be allowed to re-enroll should the enrollment profile be installed again, but will be checked out unless the block is removed. 

Removing the device block is similar to adding it:
```
mdmctl remove block -udid=564DF4AE-96B5-72BF-02C4-6B8B94BA5BF4
```

Run the remove command before enrolling the device again. 

# Known Issues
- BUG: Running the server with `-http-debug` causes the device to always see an HTTP Status OK (200) response. 
https://github.com/micromdm/go4/issues/7
- TODO: Device blocks are only allowed using a UDID. This is because the current implementation uses the `/connect` endpoint which does not normally show the serial number in the request. Adding a check-in middleware would allow blocking based on a few other attributes like serial number. 
- TODO: Listing blocked devices is not implemented. It would be possible to have `mdmctl get devices` to have custom columns, one of which could display the blocked status. 
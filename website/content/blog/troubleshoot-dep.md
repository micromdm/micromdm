+++
frontpage = true
author = "Victor Vrantchan"
date = "2017-11-06T13:04:37+02:00"
title = "How to troubleshoot your DEP/MDM Enrollments"
tags = ["mdm", "dep"]

+++

The [Device Enrollment Program](https://www.apple.com/business/dep/)(DEP) allows enterprises to configure their Macs to auto-enroll into a Mobile Device Management(MDM) server of their choice. DEP makes it possible to ensure that a new Mac [becomes managed during the unboxing process](https://blog.kolide.com/macos-on-boarding-at-kolide-fab71345986e), reducing the need for Netboot and complex imaging workflows. Of course, as any workflow that depends on the network, this enrollment process can fail, and it's important for an administrator to know how to troubleshoot it. This article has a few concrete steps that will help a macadmin using any MDM to debug why their enrollment process isn't working. 

# DEP, in brief

To make the DEP process work, Apple maintains a list of your serial numbers in a server side database. During first boot, the Mac will contact a remote server(`iprofiles.apple.com`) to get an "Activation Record". This record contains the MDM enrollment URL and a few fields that specify Setup Assistant configuration. 

{{<figure src="/troubleshoot-mdm/dep_enroll.gif" title="DEP Enrollment Screen" class="screenshot" >}}

In the WWDC session where DEP was introduced, Apple called it an _enrollment optimization_, and to this day, it lives to that characterization. If the setup assistant proceeds past the above screen, the DEP process is done, and the MDM is managing the Mac. At the very least, the enrollment profile should be installed.

Anyone can view what the current activation record for any mac looks like by running `sudo /usr/libexec/mdmclient dep nag`.

How does the activation record hold the enrollment URL for _your_ MDM? That part is configured by your MDM, by [talking to the API](https://developer.apple.com/library/content/documentation/Miscellaneous/Reference/MobileDeviceManagementProtocolRef/4-Profile_Management/ProfileManagement.html#//apple_ref/doc/uid/TP40017387-CH7-SW6) server at `mdmenrollment.apple.com`. The record that the MDM configures for your device looks roughly like the JSON structure below, which is almost identical in contents to the record your device receives. 

```
{
  "profile_name": "(Required) Human readable name.",
  "url": "https://mdm.acme.co/mdm/enroll",
  "allow_pairing": true,
  "is_supervised": false,
  "is_mandatory": false,
  "await_device_configured": false,
  "is_mdm_removable": true,
  "anchor_certs": [],
  "supervising_host_certs": [],
  "skip_setup_items": ["AppleID", "Android", "TOS"],
  "devices": ["SERIAL1","SERIAL2"]
}
```

As you can see, the Device Enrollment Program is nothing more than a matchmaking API to connect your device to its designated MDM server. I've attempted to summarize the DEP flow in the diagram below.

{{<figure src="/troubleshoot-mdm/dep_process.png" title="DEP Process" class="screenshot" >}}

# Troubleshooting Steps

While developing MicroMDM, we ran the enrollment steps many times, and often needed a way to find out why things weren't happening as we expected them to. Over time a common list of troubleshooting steps emerged. [Owen Pragel](https://twitter.com/opragel) wrote up the most common techniques to debug DEP and MDM on the [micromdm wiki](https://github.com/micromdm/micromdm/wiki/Troubleshooting).

### Debug Logging

To enable debug logging for the MDM processes on your Mac, [install this configuration profile](https://gist.github.com/opragel/2b9c518f9a27dce787ed45da832708e2). If you have a VM you're using to test a DEP workflow, it might be a good idea to add this profile to the VM image.

The `log` utility on macOS 10.12+ allows streaming very detailed logs with the `stream --info --debug` subcommand.
To filter out MDM specific log message, start by adding the following predicates.

```
log stream --info --debug --predicate 'subsystem contains "com.apple.ManagedClient.cloudconfigurationd"'
log stream --info --debug --predicate 'processImagePath contains "mdmclient" OR processImagePath contains "storedownloadd"'
```

As you're watching the log stream, you might want to adjust the predicates to add or remove new conditions. [This blog](https://eclecticlight.co/2016/10/01/using-the-logs-in-sierra-some-practical-tips/) is a good primer on the log command. `man log` has some additional usage examples as well.

### Build a VM Image

The way I recommend you test your MDM workflows is by creating a VMWware Fusion VM using [AutoDMG](https://github.com/MagerValp/AutoDMG) and Joseph Chilcote's [vfuse](https://github.com/chilcote/vfuse). If you snapshot the VM before first boot, you'll be able to reset to the snapshot and re-enroll. 


My `vfuse` template, which sets the serial number to DEP enabled Mac and uses an AutoDMG image:
```
{
    "fusion_path": "",
    "source_dmg": "/Users/groob/Desktop/osx-10.13.1-17B48.apfs.dmg",
    "output_dir": "/Users/groob/Desktop",
    "output_name": "dep-hs-groob",
    "cache": false,
    "mem_size": 4096,
    "disk_type": 0,
    "bridged": false,
    "mac_address": "",
    "enable3d": false,
    "vnc_port": 5901,
    "vnc_passwd": "",
    "hw_model": "MacBookAir7,2",
    "serial_number": "C02T9ZI1CXC4"
}
```

Build it with `sudo /usr/local/vfuse/vfuse -t dep-hs-groob.json --snapshot`.

_Note_ I really want to stress out that you must snapshot your VM before first boot. Otherwise the activation record is cached, and any changes you make to it on the MDM side, will not be reflected in your VM.

### Reset an enrollment

Thanks to Owen, we have one more useful [script to share](https://gist.github.com/opragel/12555098f5894267c3aba2a7c023a823) which can be used to trigger your test mac to an un-enrolled state. This is good for testing, although it's becoming clearer that over time VM or APFS snapshots are going to be the recommended way to reset your image.

If you're still using 10.12 or earlier, you can run these steps instead:

```
# Remove indicator that setup assistant has already run
sudo rm /var/db/.AppleSetupDone

# Clear all configuration profiles off machine (not entirely clean)
sudo rm -rf /var/db/ConfigurationProfiles/

# Remove Apple Push Notification service daemon keychain
sudo rm /Library/Keychains/apsd.keychain

# Reboot the machine. It should bring you back to setup assistant
# where you can re-enroll using DEP.
```

# Share your tips

You can find the above and more on the MicroMDM [wiki](https://github.com/micromdm/micromdm/wiki/Troubleshooting).
If you have a new tip to share, don't hesitate to edit the page and add it. 

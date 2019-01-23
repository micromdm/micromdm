# BoltDB maintenance

## Make a backup!

Before performing any operations on the underlying database please make a backup of your micromdm database!

## boltcopy

[boltcopy](https://github.com/jessepeterson/boltcopy) is a simple tool for copying BoltDB databases while omitting/including specific buckets. To get it:

```bash
go get github.com/jessepeterson/boltcopy
```

## Remove ARCHIVE buckets

Once upon a time MicroMDM archived every command response and checkin message from the device. These are not practically useful and have been removed from MicroMDM. However if you're still using the same database as back then you can clear them out:

```
boltcopy -b mdm.Checkin.ARCHIVE -b mdm.Command.ARCHIVE micromdm.db micromdm-diet.db
```

## "Reset" device enrollments

Do a lot of testing? Need to very quickly reset device enrollments but keep the existing APNS, DEP, and server config settings? Try this:

```
boltcopy -b mdm.Devices -b mdm.DeviceIdx -b mdm.PushInfo micromdm.db micromdm.nodevices.db
```
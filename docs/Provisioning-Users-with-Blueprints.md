# Provisioning users with blueprints

## Introduction

With blueprints, it's possible to provision administrative users automatically upon device enrolment. This could be useful if you want a standard service account to be included for example.

## Creating Administrative users

To create an administrative user, run `mdmctl apply users -template > exampleUser.json`

Fill in the appropriate user name values

```shell
{
  "user_shortname": "svcadmin",
  "user_longname": "Administrator",
  "hidden": false
}
```

Now upload the user to the MicroMDM server along with the desired password for the user using `mdmctl apply users -f ./exampleUser.json -password SuperSecret` . You can check to see that the user was created by running `mdmctl get users`

```shell
UUID                                  UDID  UserID  UserShortName  UserLongName
60000000-faaa-4000-c111-000aaa00000                svcadmin       Administrator
```

You'll also notice that if you open the original JSON file (`exampleUser.json` in our example), it will have been updated to include a `password_hash` and `uuid`.

## Blueprint Templates

If you generate a new blueprint template (`mdmctl apply blueprints -template`) you'll notice that there is a space for User UUIDs. Put the UUID that you got from `mdmctl get users` into the `user_uuids` section in order to automatically provision users on enrolment.

```shell
{
  "uuid": "uuid-here",
  "name": "exampleName",
  "install_application_manifest_urls": [
    "https://mdm.example.org/repo/exampleAppManifest.plist"
  ],
  "profile_ids": [
    "com.example.my.profile"
  ],
  "user_uuids": [
    "60000000-faaa-4000-c111-000aaa00000"      <--- Admin Account User IDs
  ],
  "skip_primary_setup_account_creation": true,
  "set_primary_setup_account_as_regular_user": true,
  "apply_at": [
    "Enroll"
  ]
}
```

Create your blueprint file with the UUID and any other settings you want, then upload it to micromdm using `mdmctl apply blueprints -f /path/to/your/blueprint.json`

The next step is that your DEP profile has to have `await_device_configured` set to `true`.

It's worth noting that if you skip Primary user account creation, you must also set `set_primary_setup_account_as_regular_user` to `true`.

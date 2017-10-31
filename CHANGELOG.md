# v1.2.0 October 31 2017

* Added support for modifying the default enrollment profile.
* Added support for user level profiles. 
* Added support for AccountConfiguration during DEP Enrollment. Specified in blueprints
* Addes support for multiple server configs in `mdmctl`.
* Added `mdmctl mdmcert upload` command which uploads/replaces the servers push certificate.
* Incorporated certhelper into mdmctl. See `mdmctl mdmcert -h`
* Added ENV variables for sensitive flags: `MICROMDM_APNS_KEY_PASSWORD`,`MICROMDM_API_KEY`
* Removed the `-redir-addr` flag. Redirect to HTTPS is only enabled when the 443 port is used.

# v1.1.0 June 05 2017

* Import and sign pkgs, generate appmanifest on import.
* Support syncing devices from DEP when token is added.
* Option to include SSL certificates in DEP profile template (-anchor and -use-server-cert) #107
* /push and /v1/commands API endpoints require API authentication #157
* Add `mdmctl` binary for interacting with the server over API. #127
* Save DEP cursor for use after restart. #109
* Add `-examples` flag to micromdm serve. #119
* Add HTTP logger (with apache format) for all endpoints. #85
* Serve a basic homepage at `/`. #113
* Decrypt armored private keys if the `-apns-password` flag is specified by the user. #105
* Improved command queue handling of NotNow and other responses. #96
* Fixed bug that allowed for duplicate device records on re-enrollment. #125
* Fixed data race in pubsub package. #97
* Fixed bug that would cause PushInfo Token for the device to be replaced by one for the user. #90


## [v1.5.0](https://github.com/micromdm/micromdm/compare/v1.4.0...master) (Unreleased)

* Fix DEP token update issue (#513, #510)
* Refactor certificate verification and implement UDID-cert authentication (#358, #429)
* Cleanup DEP library and integrate into main project (#504, #505)
* Add API endpoint to retrieve APNS certificate (#503)
* Remove deprecated `-apns` flags from server startup (#528)
* Move API calls to list endpoints from HTTP GET to HTTP POST (#522, #523, #524, #525, #526)
* Add support for the ApplicationConfiguration Setting (#521)
* Allow SCEP client validity to be adjusted via server startup flag (#577)
* Fix bug in mdmctl server saving, switch config when saving automatically (#565, #566)

## [v1.4.0](https://github.com/micromdm/micromdm/compare/v1.3.1...v1.4.0) September 6 2018

### Stability Improvements

* Handle DEP INVALID_CURSOR response (#497)
* Use config for block push (#479, #480)
* No longer store SCEP CA on disk or include in enrollment profile (#490)
* Further SCEP fixes (#492, #493)
* Base64 fixes for API CLI tools (#477)
* `mdmctl apply block` now works with self-signed certs (#480)
* Add API CLI tool for dep sync (#481)
* DeviceInformation command API example support query strings (#469)
* Allow setting curl options in environment variable (#455)
* Fix URL params decoding. (#467)
* Reorganize/refactor server init (#458)
* Allow supplying additional `curl` options in API CLI tools (#455)

Thanks to our contributors for this release: @erikng, @gerardkok, @knightsc, @marpaia, and @ochimo!

## [v1.3.1](https://github.com/micromdm/micromdm/compare/v1.3.0...v1.3.1) July 10 2018

* Update base container to Alpine 3.7 (#437)
* Fix bugs in SCEP enrollment (#451)
* Fix issue with APNS timeouts -- Issue #215 (#446)
* Add device_information and security_info commands with curl API (#448)
* Add support for InstallEnterpriseApplication command (#452)

## [v1.3.0](https://github.com/micromdm/micromdm/compare/v1.2.0...v1.3.0) 

### Auto-assigner

* Reorganize/refactor MDM, device, webhook services. #423, #424, #425, #426, #427
* Do not allow `mdmctl config set` without args. #421
* Fix for multiple UDID records. #422
* Added/refactored logging. #405, #425
* Added `-homepage` switch. #420
* Warn about deprecated APNS switches. #412
* Disallow bad TLS configuration with `-tls=false`. #414
* Refactored MDM types. #341, #415
* Added DEP auto-assigner feature. #405
* Fixed bug with authentication error messages. #411
* Added support for querying devices by serial(s). #363
* Added support for triggering a DEP sync via API. #404
* Added support for mdmcert.download directly to `mdmctl` #401
* Reject network MDM user attempts until we add support. #379
* Warn when starting without an API key. #396
* Added tools and documentation for ngrok, curl, and APIs. #392
* Fix issue with MDM command `AvailableOSUpdates` parsing. #368
* Validate APNs Push Certificate Topic. #373
* `mdmctl` now outputs to stdout vs. stderr. #360
* Added common HTTP library `httputil`. #350
* Added project Code of Conduct. #334
* Refactored services (mostly for HA). #348, #349, #351, #352, #353, #354, #355, #359
* Reorganized project layout. #333, #335, #336, #338, #340, #347
* Added support for version API. #327
* Added command response webhook feature. #315
* Added support for supplied `depsim` URL. #318
* Added Dockerfile. #316

## [v1.2.0](https://github.com/micromdm/micromdm/compare/v1.1.0...v1.2.0) October 31 2017

### User Profiles

* Added support for modifying the default enrollment profile.
* Added support for user level profiles.
* Added support for AccountConfiguration during DEP Enrollment. Specified in blueprints
* Addes support for multiple server configs in `mdmctl`.
* Added `mdmctl mdmcert upload` command which uploads/replaces the servers push certificate.
* Incorporated certhelper into mdmctl. See `mdmctl mdmcert -h`
* Added ENV variables for sensitive flags: `MICROMDM_APNS_KEY_PASSWORD`,`MICROMDM_API_KEY`
* Removed the `-redir-addr` flag. Redirect to HTTPS is only enabled when the 443 port is used.

## [v1.1.0](https://github.com/micromdm/micromdm/compare/v1.0.0...v1.1.0) June 05 2017

### YVR!

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

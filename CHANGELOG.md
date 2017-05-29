# v1.0.1 TBD

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


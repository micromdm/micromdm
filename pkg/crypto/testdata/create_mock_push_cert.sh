#!/usr/bin/env bash

openssl req -newkey rsa:2048 -nodes -keyout key.pem -x509 -days 18262 -out certificate.pem -subj "/UID=com.apple.mgmt.External.18a16429-886b-41f1-9c30-2bd04ae4fc37/CN=APSP:17a16429-886b-41f1-8c90-3bd02ae9fc57/C=US"

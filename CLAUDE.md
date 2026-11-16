# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

MicroMDM is an API-driven Mobile Device Management server for Apple devices, written in Go. This repo is a fork (`origin` = `paragpatel/micromdm`) of upstream https://github.com/micromdm/micromdm. The fork's own changes so far extend the `InstallApplication` command (`VPNUUID`, `removable` options). Day-to-day work happens on the `working` branch.

Two binaries:
- `cmd/micromdm` — the server (`micromdm serve`). State lives in BoltDB under `-config-path` (default `/var/db/micromdm`).
- `cmd/mdmctl` — CLI client for the server's API (`get`, `apply`, `remove`, `config`, `mdmcert`, etc.).

## Commands

```bash
make deps            # go mod download
make                 # builds build/<platform>/micromdm and build/<platform>/mdmctl (version info injected via -ldflags)
make test            # go test -cover ./...
make test-race       # race tests are not the default (bbolt issue)
go test ./mdm/mdm/ -run TestName -v   # run a single test
gofmt -s -l .        # CI fails if this prints anything
```

CI builds/tests on Linux, macOS, and Windows with Go 1.18–1.20, so avoid newer-only Go features.

`tools/api/` holds curl scripts for hitting a running server (needs `jq` and an `env` file with `API_TOKEN`/`SERVER_URL`, pointed to by `MICROMDM_ENV_PATH`). `tools/ngrok/` is the recommended local dev/demo setup.

## Architecture

**Wiring:** `cmd/micromdm/serve.go` parses flags, then `server/server.go` (`Server` struct) builds all services: BoltDB, pubsub, command queue, APNS, SCEP, enrollment, DEP, webhooks. `serve.go` then mounts every service's HTTP handlers on a gorilla/mux router. Device-facing endpoints (`/mdm/...`, `/scep`, `/mdm/enroll`, `/ota/...`) are unauthenticated or cert-verified. The `/v1/...` API endpoints use basic auth with `-api-key`.

**Package convention (go-kit):** each domain under `platform/<name>/` (device, blueprint, profile, user, command, apns, config, dep, remove, appstore, ...) usually has:
- `service.go`: the service interface plus go-kit endpoints
- one file per operation (`apply_*.go`, `get_*.go`, `remove_*.go`), each holding the request/response types, endpoint, and HTTP decode/encode
- `server.go`: `RegisterHTTPHandlers`
- `client.go`: the HTTP client used by `mdmctl`
- `builtin/`: the BoltDB-backed store
- `worker.go`: a pubsub subscriber that reacts to events
- `internal/<name>proto/`: protobuf definitions for serializing events and stored records

**Event-driven core:** components talk through an in-process pubsub (`platform/pubsub/inmem`) rather than calling each other directly. Topics are string constants, e.g.:
- `mdm.Connect`, `mdm.Authenticate`, `mdm.TokenUpdate`, `mdm.CheckOut` in `mdm/service.go`
- `mdm.Command`, `mdm.RawCommand` in `platform/command`
- `mdm.CommandQueued` in `platform/queue`
- `mdm.DeviceEnrolled`

Workers (device, blueprint, user, queue, webhook) subscribe to these topics. Events are protobuf-serialized.

**Command flow:**
1. An API client POSTs to `/v1/commands`. `platform/command` builds an `mdm.CommandRequest`, turns it into a plist payload, and publishes `mdm.Command`.
2. `platform/queue` (BoltDB-backed) stores the command per UDID and sends an APNS push.
3. The device connects to `/mdm/connect`. `mdm/` checks the device identity cert and asks the queue for the next command.
4. When the device acknowledges, the queue moves to the next command and publishes the event. `workflow/webhook` forwards checkin and acknowledge events to `-command-webhook-url`.

**MDM command definitions (`mdm/mdm/`):** `command.go` defines `CommandRequest`/`Command` and the payload struct for each command. A command has several representations that must stay in sync:
- the JSON API (`marshal_json.go`/`unmarshal_json.go`)
- the plist sent to the device (`marshal_plist.go`/`unmarshal_plist.go`)
- protobuf for pubsub (`marshal_proto.go`/`unmarshal_proto.go`, defined in `internal/mdmproto/mdm.proto`)

Adding or changing a command field usually means editing `command.go`, `mdm.proto`, regenerating `mdm.pb.go`, and updating `marshal_proto.go`/`unmarshal_proto.go`. Commit `3bac510` is an example. Tests are in `mdm_command_test.go`, with fixtures in `testdata/`.

**Regenerating protobufs:** each `internal/*proto/*.go` file has a `//go:generate protoc --go_out=. --go_opt=paths=source_relative <file>.proto` directive. The generated files were made with `protoc` v3.21.2 and `protoc-gen-go` v1.28.1. A local `protoc-21.2-osx-x86_64/` directory (untracked) provides that `protoc`. Put its `bin/` on your PATH, install `protoc-gen-go@v1.28.1`, then run `go generate` in the proto package's directory.

## Production deployment

The live server runs on `parag@85.90.246.210` (Ubuntu 22.04, a shared box that also runs many unrelated Docker services). Its public hostname is `https://mdm.vartit.com`.

- **systemd:** `micromdm.service` runs `/home/parag/mdm/build/linux/run.sh`. Logs: `journalctl -u micromdm` (`MICROMDM_HTTP_DEBUG=true`, so full request and response bodies are logged).
- **Binary:** `/home/parag/mdm/build/linux/micromdm` is built from this fork; `micromdm version` reports `v1.11.0-19-g475ea58`. It's a prebuilt Linux binary copied over (`make xp-micromdm`). `/home/parag/mdm/micromdm/` on the server is an old upstream checkout and is *not* what's running.
- **Data:** BoltDB is at `/var/db/micromdm/micromdm.db` (default `-config-path`). The push certificate and DEP config are stored in this DB.
- **Request path:**
  - Devices go through Cloudflare, then nginx on `:443` (`/etc/nginx/conf.d/mdm.conf`, Cloudflare-only allowlist), then `https://localhost:8443` (micromdm, self-signed `server.crt`/`server.key`).
- **Webhook path:**
  - micromdm sends webhooks to `-command-webhook-url http://mdmwebhook.local/webhook_/webhook`. `mdmwebhook.local` is `127.0.0.1` in `/etc/hosts`.
  - nginx on `127.0.0.1:80` proxies that to a Flask app on `localhost:8444` (`build/linux/webhook/webhook.py`, started by `run.sh`).
  - The Flask app publishes the raw event JSON to Redis channel `mdm_webhook`, and the decoded acknowledge plist to `mdm_webhook_payload`. That Redis server is external (`172.104.178.91:6379`), and downstream consumers read from it.
- **VPN:** MicroMDM does not issue VPN certificates.
  - An IKEv2 VPN runs in Docker (`ipsec-vpn-server`, image hwdsl2) on the same host. Its CA, server, and `vpnclient` certificates were all created 2024-03-14 and are valid until 2034, stored in `/etc/ipsec.d` inside the container.
  - The client profile `~/vpnclient.mobileconfig` (a `com.apple.vpn.managed` payload plus PKCS12 and root CA payloads) is pushed to devices as a profile.
  - The fork's `InstallApplication` `vpn_uuid` attribute ties an installed app to that VPN's `VPNUUID` (per-app VPN).
  - `~/micromdm-1.11.0/` on the server is an older source snapshot containing only the first `VPNUUID` commit; this git repo is the source of truth.
- `run.sh` sets `GODEBUG=x509sha1=1` and passes the API key on the command line.

### Certificates and backups

| What | Lifetime | How it's renewed |
|---|---|---|
| **APNs push certificate** | Yearly; current one expires 2027-11-16 | The `mdmctl mdmcert.download` API is deprecated, so use the mdmcert.download website, then **Renew** at identity.apple.com |
| **Device identity certificates** | 365 days (`-scep-client-validity`) | `InstallProfile` with the backend's per-child enrollment profile (identifier `com.github.micromdm.micromdm.enroll`) before expiry. Expired devices can't be reached. |
| **SCEP CA** | 5 years; MicroMDM created it at install, so it expires **2028-02-03** | No built-in way to rotate it |

- Step-by-step runbooks and the backup copies of certificates, keys and server config are in the internal repo `~/projects/mdm` (`server/backup-YYYY-MM-DD/README.md`).
- `tools/dbtools/` has read-only BoltDB inspectors (device/certificate report, push certificate and SCEP CA export). Run them on a **copy** of the DB.
- `/boltbackup` times out on large DBs, so take backup copies with the service stopped.

## Main consumer: KidsNanny

The only API client is the KidsNanny Spring Boot backend (`~/projects/kidsnanny`, package `com.kidsnanny.mdm`), which runs on another host (172.104.178.91).

- **Enrollment:** the backend serves `/parent/api/mdm/profile/{childDeviceId}`. It fetches `/mdm/enroll` from here, rewrites `CheckInURL` to its own proxy (`dash.kidsnanny.com.au/.../api/mdm/checkin/{childDeviceId}`), signs the result, and returns it. Check-ins therefore reach this server through the backend, and the backend maps UDIDs to children from them.
- **Commands it sends** (`/v1/commands`):
  - InstallProfile: restrictions, blocked apps, web content filter, and a VPN profile for the iOS app's own packet tunnel
  - RemoveProfile
  - DeviceInformation
  - InstalledApplicationList
  - ProfileList
  - InstallApplication with `vpn_uuid`, which takes over management of installed browsers so they use the per-app VPN. The per-device VPNUUID is a base UUID plus a rotating index. This is why the fork adds `vpn_uuid` and `removable`.
- **Results:** they come back only through the webhook → Redis channel `mdm_webhook` (`acknowledge_event`).
- Changing a command's JSON field names or the webhook event shape here breaks that backend.

## Docs

User-facing docs are in `docs/user-guide/` (introduction, quickstart, enrolling devices, API and webhooks).

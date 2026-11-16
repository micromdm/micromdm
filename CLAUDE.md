# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

MicroMDM is an API-driven Mobile Device Management server for Apple devices, written in Go. This repo is a fork (`origin` = `paragpatel/micromdm`) of upstream https://github.com/micromdm/micromdm.

**Upstream MicroMDM v1 is in maintenance mode; its support ended at the end of 2025.** Migrating to NanoMDM is planned: micro2nano `llorne` moves enrollments and `cmdapi` translates the JSON command API. Gaps to cover in the migration: enrollment profile serving, SCEP (reuse the exported CA), `/v1/devices`, and the fork's `vpn_uuid`/`removable` attributes.

Fork changes:
- `InstallApplication` attributes `vpn_uuid` and `removable`.
- `tools/dbtools`.
- The temporary flag `-accept-expired-device-certs-until` (see Certificates).

**Git workflow:**
1. Work on a feature branch.
2. Open a PR into `main`.
3. Merge `main` into `working`, the deploy branch.
4. Build from `working`.

Upstream syncs arrive as PRs from `micromdm:main`. The fork has no CI configured, so run `go test ./...` and `gofmt -s -l .` locally.

Runbooks (slash commands in `.claude/commands/`):
- `/mdm-status`: read-only health check.
- `/push-cert-renewal`: yearly APNs renewal.
- `/micromdm-deploy`: build, upgrade, rollback.

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
- **Binary:** `/home/parag/mdm/build/linux/micromdm`, built from `working` with `make xp-micromdm`. Since 2026-11-16 it reports `v1.12.1-23-gfdca542`.
  - The previous binary is kept as `micromdm-475ea58`; `upgrade-fdca542.sh` / `rollback-475ea58.sh` live next to it (run with `sudo` by the user).
  - `/home/parag/mdm/micromdm/` on the server is an old upstream checkout and is *not* what's running.
- **Access:** `ssh parag@85.90.246.210` works with a key. `sudo` needs the user's password, so stopping or restarting the service and swapping binaries are the user's steps.
- **Data:** BoltDB is at `/var/db/micromdm/micromdm.db` (default `-config-path`). The push certificate and DEP config are stored in this DB.
- **Request path:**
  - Devices go through Cloudflare, then nginx on `:443` (`/etc/nginx/conf.d/mdm.conf`, Cloudflare-only allowlist), then `https://localhost:8443` (micromdm, self-signed `server.crt`/`server.key`).
- **Webhook path:**
  - micromdm sends webhooks to `-command-webhook-url http://mdmwebhook.local/webhook_/webhook`. `mdmwebhook.local` is `127.0.0.1` in `/etc/hosts`.
  - nginx on `127.0.0.1:80` proxies that to a Flask app on `localhost:8444` (`build/linux/webhook/webhook.py`, started by `run.sh`).
  - The Flask app publishes the raw event JSON to Redis channel `mdm_webhook`, and the decoded acknowledge plist to `mdm_webhook_payload`. That Redis server is external (`172.104.178.91:6379`), and downstream consumers read from it.
- **VPN:** MicroMDM does not issue VPN certificates.
  - The KidsNanny "VPN" is the iOS app's own local packet tunnel (`com.kidsnanny.ios.kidsPacketTunnel`). It is configured by a VPN profile that the backend pushes.
  - `vpn_uuid` makes browsers use it per-app.
  - The IKEv2 container on this host (`ipsec-vpn-server`, hwdsl2; certs valid to 2034) is unrelated to the KidsNanny flow.
  - `~/micromdm-1.11.0/` on the server is an old source snapshot; this git repo is the source of truth.
- `run.sh` sets `GODEBUG=x509sha1=1` and passes the API key on the command line.

### Certificates and backups

| What | Lifetime | How it's renewed |
|---|---|---|
| **APNs push certificate** | Yearly; current one expires 2027-11-16 | The `mdmctl mdmcert.download` API is deprecated, so use the mdmcert.download website, then **Renew** at identity.apple.com |
| **Device identity certificates** | 365 days (`-scep-client-validity`) | The KidsNanny backend job (`mdm.cert-renewal.*`, daily 03:40 UTC) sends `InstallProfile` with the per-child enrollment profile (identifier `com.github.micromdm.micromdm.enroll`) 30 days before expiry. The phone gets a new SCEP cert and re-Authenticates, and Authenticate saves the new cert hash. |
| **SCEP CA** | 5 years; MicroMDM created it at install, so it expires **2028-02-03** | No built-in way to rotate it |

- Step-by-step runbooks and the backup copies of certificates, keys and server config are in the internal repo `~/projects/mdm` (`server/backup-YYYY-MM-DD/README.md`).
- `tools/dbtools/` has read-only BoltDB inspectors (device/certificate report, push certificate and SCEP CA export). Run them on a **copy** of the DB.
- `/boltbackup` times out on large DBs, so take backup copies with the service stopped.

**Temporary expired-certificate recovery** (`pkg/crypto/expired.go`, flag `-accept-expired-device-certs-until 2027-02-15` in `run.sh`):
- **What it does:** devices whose identity cert already expired were rejected before any command could reach them. While the flag is set, a request whose only failure is the signer cert's validity window is accepted when:
  - the cert is really expired;
  - it was issued and signed by our SCEP CA;
  - the message digest and signature verify.

  The UDID↔cert binding still applies.
- **How the phones recover:** the KidsNanny overlay `mdm.cert-renewal.include-expired=true` makes the backend send those phones the renewal profile.
- **Cleanup:** after 2027-02-15, remove both the flag and `include-expired`.

**Gotchas:**
- **iOS sends BER, not DER.** The `Mdm-Signature` CMS uses indefinite lengths, which Go's `encoding/asn1` rejects. Parse with `pkcs7.Parse` (it converts) and work on the parsed object.
- **The skew flag can't fix expired certs.** The pkcs7 "signing time ... outside of certificate validity" check compares the device's own signing time with the cert, so `-device-signature-skew` cannot help.
- **The deprecated mdmcert API hides its reason.** It returns `"failure"`; the real reason is in the `reason` field, which upstream commit 90831d0 prints.
- **The device record's `Enrolled` field is unreliable.** In the BoltDB device record it reads false; use the API's `enrollment_status`.
- **Queues grow without limit.** They hold undeliverable commands plus full command history (`-no-command-history` is not set). Clear dead devices with `DELETE /v1/commands/{udid}`.

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

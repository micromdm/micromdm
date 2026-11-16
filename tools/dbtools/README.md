# dbtools

Read-only helpers for inspecting a MicroMDM BoltDB **backup copy** (BoltDB takes an exclusive lock,
so run these against a copy, never the live `/var/db/micromdm/micromdm.db`).

Build for the Linux server: `GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o <name> ./tools/dbtools/<name>`

| Tool | Usage | Output |
|---|---|---|
| `devicereport` | `devicereport <db> > report.csv` | CSV per device: udid, serial, model, last_seen, identity-cert NotAfter (via `mdm.UDIDCertAuth` hash → `scep_certificates`), queue bytes. Bucket sizes on stderr. Note: `enrolled` from the DB is unreliable; use the API `enrollment_status`. |
| `pushcert` | `pushcert <db> <out-prefix>` | `<prefix>.pem` / `<prefix>.key` — the active APNs push cert + key (for rollback/backup). |
| `scepca` | `scepca <db> <out-dir>` | `scep-ca.pem` / `scep-ca.key` — the SCEP CA that signs device identity certs. |

`/boltbackup` cannot export large DBs (go4 httputil `WriteTimeout` is 60s); take a consistent copy with the
service stopped instead.

---
description: Build, deploy, and roll back the MicroMDM fork on production (85.90.246.210). The user runs the sudo step.
---

# Deploy the fork to production

Git flow (see CLAUDE.md):
1. Feature branch → PR into `main`.
2. Merge `main` into `working`.
3. Build from `working`.

Production runs `v1.12.1-23-gfdca542` (deployed 2026-11-16; before that `v1.11.0-19-g475ea58`).

1. **Build.** On `working`, run `make xp-micromdm`, then `build/linux/micromdm` and `micromdm version` should show the commit.
   `go test ./...` must pass (CI isn't configured on the fork).
2. **Stage the binary next to the running one.**
   `scp build/linux/micromdm parag@85.90.246.210:mdm/build/linux/micromdm-<sha>`, `chmod 755`.
   Check `./micromdm-<sha> version` and `serve -h` for any new flags.
3. **Write an upgrade script** for the user, following `~/mdm/build/linux/upgrade-fdca542.sh`. It must:
   - stop the service;
   - take a **consistent DB backup** to `~/mdm-backups/` (the service must be stopped; `/boltbackup` times out on this DB);
   - keep the old binary as `micromdm-<oldsha>` and `run.sh` as `run.sh.pre-<sha>`;
   - `install -m 755 -o root` the new binary;
   - edit `run.sh` (dry-run the sed first);
   - start the service and print status, version and startup errors.

   Also write the matching rollback script (`rollback-475ea58.sh` is the template).
   Run `bash -n` on both scripts.
4. **The user runs** `sudo bash ~/mdm/build/linux/upgrade-<sha>.sh`.
5. **Verify:**
   - `/mdm/connect` statuses since `ActiveEnterTimestamp`;
   - no panics or UDID mismatches;
   - test push to an active device → `success` and a 200 connect;
   - `/mdm-status`.
6. Record the version in CLAUDE.md, and the backup in `~/projects/mdm`.

## run.sh flags in use

- `-api-key`, `-server-url https://mdm.vartit.com`, `-tls-cert`/`-tls-key` (self-signed, behind nginx), `-http-addr :8443`
- `-command-webhook-url http://mdmwebhook.local/webhook_/webhook`
- `-accept-expired-device-certs-until 2027-02-15`: **temporary** expired-cert recovery (see CLAUDE.md). Remove it after that date, together with the KidsNanny overlay line `mdm.cert-renewal.include-expired=true`.
- `GODEBUG=x509sha1=1`, `MICROMDM_HTTP_DEBUG=true` (full requests in the journal, including `Mdm-Signature`).

Candidate flag for the next maintenance window: `-no-command-history`. After clearing queues, also run an offline `bbolt compact` (with the service stopped).

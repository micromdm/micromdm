---
description: Read-only health check of the production MicroMDM server and the KidsNanny cert-renewal job (service, version, push cert, connect success rate, expired-cert recovery, renewals).
argument-hint: "[status | renewals | device <udid-or-serial>]"
---

# MDM status (read-only)

Production MicroMDM: `ssh parag@85.90.246.210` (https://mdm.vartit.com). KidsNanny backend: `ssh KidsNanny`
(172.104.178.91, prod JVM in `/home/ec2-user/tmp-prod`). **Never** restart, edit, or write anything from this
command; `sudo` on the MDM box is the user's. Never print the API key (read it into a variable from `run.sh`).

## status (default)

Run on the MDM box and summarise as a table:

```bash
ssh parag@85.90.246.210 'systemctl is-active micromdm; ~/mdm/build/linux/micromdm version | head -1
K=$(grep -o "api-key [^ ]*" ~/mdm/build/linux/run.sh | cut -d" " -f2)
curl -sk -u "micromdm:$K" https://localhost:8443/v1/config/certificate | python3 -c "import sys,json,base64;open(\"/tmp/.c\",\"wb\").write(base64.b64decode(json.load(sys.stdin)[\"cert\"]))"
openssl x509 -inform der -in /tmp/.c -noout -subject -enddate; rm /tmp/.c
grep -o "accept-expired-device-certs-until [^ ]*" ~/mdm/build/linux/run.sh
journalctl -u micromdm --since "24 hours ago" --no-pager | grep -o "status=[0-9]* proto=[^ ]* host=[^ ]* user_agent=MDM/1.0 path=/mdm/[a-z]*" | awk "{print \$1, \$5}" | sort | uniq -c
journalctl -u micromdm --since "24 hours ago" --no-pager | grep -c "accepted expired device certificate"
journalctl -u micromdm --since "24 hours ago" --no-pager | grep -c "outside of certificate validity"
journalctl -u micromdm --since "24 hours ago" --no-pager | grep "component=apns" | grep -vc "err=null"
df -h / | tail -1; ls -la /var/db/micromdm/micromdm.db'
```

Flag in the report:
- **APNs push certificate** expiring within 60 days (it is renewed yearly, see `/push-cert-renewal`)
- **SCEP CA** expires **2028-02-03**
- `outside of certificate validity` > 0 while the recovery flag is active
- APNs errors
- `/mdm/connect` 500 share

The journal keeps only about a day.

## renewals

KidsNanny job (`mdm.cert-renewal.*`, daily 03:40 UTC):

```bash
ssh KidsNanny 'grep -h "MDM cert renewal scheduler is" /home/ec2-user/tmp-prod/logs/kidsnanny.log | tail -1
grep -h "mdm cert renewal" /home/ec2-user/tmp-prod/logs/kidsnanny*.log | grep -v "scheduler is" | tail -40'
```

Count:
- `queued InstallProfile` lines (split `cert EXPIRED` vs `expires`)
- `now has a new certificate` lines (successful renewals)
- `failed`/`NOT sent` lines

DB state needs the user, because the auto-mode classifier blocks Claude from prod DB reads:
`SELECT COUNT(*) total, SUM(mdm_cert_not_after < NOW()) expired, SUM(mdm_cert_not_after BETWEEN NOW() AND NOW() + INTERVAL 30 DAY) due_30d, SUM(mdm_cert_renewal_count >= 5) gave_up FROM mdm_devices;`

## device <udid-or-serial>

1. Find the UDID from `/v1/devices` (match `serial_number`).
2. Trace the journal: split on `---BEGIN Request---` and keep the blocks containing the UDID. Report MessageType/Status, HTTP status and `err=`.
3. The signer cert of the latest request is in the `Mdm-Signature` header; decode with
   `openssl pkcs7 -inform der -print_certs | openssl x509 -noout -enddate`.

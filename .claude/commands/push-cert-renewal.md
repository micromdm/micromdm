---
description: Yearly APNs push certificate renewal for production MicroMDM (mdmcert.download website flow + identity.apple.com Renew + upload). Next due before 2027-11-16; start in October.
---

# APNs push certificate renewal (yearly)

Current cert:
- topic `com.apple.mgmt.External.78ac43f5-0e85-4cdb-b807-e2207156b24b`
- valid to **2027-11-16**
- vendor "MDM Vendor: Jesse Peterson" via mdmcert.download

If the topic changes, **every device must re-enroll**.

1. **CSR on the server.** In a new empty dir: `~/mdm/build/linux/mdmctl mdmcert.download -new -email=info@vartit.com`.
   - The v1 API is deprecated, so the command fails with `"failure"`. It still leaves `mdmcert.download.push.key` and `.csr`.
   - Delete the `mdmcert.download.pki.*` files; they're not needed.
   - Copy the CSR to the Mac.
2. **mdmcert.download.**
   1. https://mdmcert.download → Sign in with GitHub (Parag's personal account; registration "Vartit Technology Private Limited", India).
   2. Sign page → paste the CSR PEM. The file picker may not be usable from the browser tools.
   3. Sign → Download blob (`apns_push.plist.b64.req`).

   Notes:
   - The "This is a renewal" checkbox only appears for accounts with prior requests; it is cosmetic.
   - The email shown in the page's example command comes from the GitHub login; it is not used.
3. **Verify the blob** before Apple:
   1. base64-decode it to a plist.
   2. `PushCertRequestCSR` must equal our CSR.
   3. The SHA-256 signature must verify with the vendor cert.
   4. The chain must verify.
4. **Apple.**
   1. https://identity.apple.com/pushcert as **satishkpatel@gmail.com** (the user signs in).
   2. Find the row with topic `…78ac43f5…` and press **Renew** (never "Create a Certificate").
   3. Upload the blob and download `MDM_<vendor>_Certificate.pem`.

   If Apple answers "Request Entity Too Large / Invalid file type or size", rebuild the blob with the chain **vendor → WWDR → root** (mdmcert v2 emits root first) and retry.
5. **Pre-flight on the server** (all must pass):
   - same topic as the active cert;
   - new cert pubkey == new key pubkey;
   - issuer Apple Application Integration 2;
   - valid now.

   Also export the current pair for rollback (`tools/dbtools/pushcert` on a DB copy).
6. **Upload** on localhost (no restart):
   `PUT https://localhost:8443/v1/config/certificate` with `{"cert": base64(pem), "key": base64(key)}`.
   `mdmctl mdmcert upload` works too, but `~/.micromdm/servers.json` has a placeholder token.
7. **Verify:**
   1. `GET /v1/config/certificate` shows the new expiry.
   2. `GET /push/<active udid>` → `success`.
   3. The device connects with 200.
8. **Back up** into `~/projects/mdm/server/backup-YYYY-MM-DD/apns-push/` (cert, key, CSR, blob, previous pair) and push to the internal GitLab.

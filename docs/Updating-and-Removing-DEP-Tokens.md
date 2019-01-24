# Updating and removing DEP tokens

Due to a current issue (#546), MicroMDM is able to store as many DEP tokens as possible BUT only uses the first loaded, even if expired.

## Updating and Downloading the new DEP Token

1. Log in to either [Apple School Manager](https://school.apple.com) or [Apple Business Manager](https://business.apple.com)
2. Browse to your MicroMDM server under MDM Servers.
3. Download your updated token.

## Removing Previous and Uploading a new DEP Token

1. Log in to a terminal session of your MicroMDM server (SSH, Google Cloud Shell, DigitalOcean "Access Console", etc).
2. Download [BoltBrowser](https://github.com/br0xen/boltbrowser) for your OS.
3. Stop the MicroMDM service.
4. **Make a backup of your MicroMDM Database!** `cp /var/db/micromdm/micromdm.db /var/db/micromdm/micromdm_Backup.db`
5. Access the DB with BoltBrowser: `/path/to/boltbrowser /var/db/micromdm/micromdm.db`
6. Press down until you highlight the mdm.DEPToken table.
7. On the CK Keys, press SHIFT+D to delete all of the CK entries.
8. Press `q` or `ESC` to exit BoltBrowser.
9. Start the MicroMDM service.
10. Using `mdmctl`, upload the new DEP Token: `/path/to/mdmctl apply dep-token -f /path/to/new/token.p7m`

_In some cases you may need to restart the MicroMDM service after uploading a new DEP Token_
11. Stop the MicroMDM service, then Start the MicroMDM service.

## Package Repo

The `-filerepo` flag points to a folder that contains signed packages and [application manifests](http://help.apple.com/deployment/osx/#/ior5df10f73a).
How you choose to organize the directory is up to you, the server will make the files and directories available at `https://serverurl/repo`

### Creating and Signing Packages
Note that only distribution style packages are supported.
Use `productbuild` to turn a component flat pkg into a distribution one.
```bash
productbuild --package someFlatPkg.pkg myNewPkg.pkg
```

Example signing using munkitools pkg.
You will need an app signing certificate from apple.
```bash
/usr/bin/productsign --sign "3rd Party Mac Developer Installer: Your Name (ID)" \
    ~/Downloads/munkitools-2.7.0.2757.pkg \
    /path/to/pkgrepo/pkgs/munkitools-2.7.0.2757.pkg
```

### Application Manifest
Documented by Apple [here](http://help.apple.com/deployment/osx/#/ior5df10f73a).
Pepijn Bruienne first described how to create an manifest plist for OS X pkgs [on his blog](http://enterprisemac.bruienne.com/2015/11/17/installing-os-x-pkgs-using-an-mdm-service/).

The default `md5-size` value is 10MB, but can be set to anything. The array of md5 hashes must match the md5-size specified in the md5-size field.
In case of a connection error, the client will retry the download starting with the last succesful chunk, so it's a good idea to keep the sizes relatively small.

You can manually create a plist, or you can use the [`appmanifest`](https://github.com/micromdm/tools/releases/tag/v1) utility.
```bash
./appmanifest -url https://mdm.example.com/repo/pkgs/munkitools-2.7.0.2757.plist \
	/path/to/signed/munkitools-2.7.0.2757.pkg > /path/to/pkgrepo/manifests/munkitools-2.7.0.2757.plist
```

OR
```bash
./mdmctl apply app -pkg /path/to/signed/munkitools-2.7.0.2757.pkg > /path/to/pkgrepo/manifests/munkitools-2.7.0.2757.plist
```

{{< note title="Metadata" >}}
The metadata portion of the plist seems to only matter for launchpad. It might be safe to omit it, but you can edit the plist and add it yourself as needed.
The information supplied in the metadata section should be accurate, but if you don't include any items in the items array, or use a bad identifier it does not seem to matter.
{{< /note >}}

# Configuring LetsEncrypt with Lego Client

As of this writing, MicroMDM's usage of Lets Encrypt (a free SSL certificate service) is currently broken. In the meantime, it is recommended to use the [Lego](https://github.com/xenolf/lego) client to download your certs.

On the server you are hosting MicroMDM, you may need to temporarily shut down the server if you bind it on port `80`, although Lego does support using a different port via the `--http :####` format.

1. Download the latest Lego release for the os you are hosting from [here](https://github.com/xenolf/lego/releases/latest).

2. Un-tar the tarball that is downloaded `tar -xvf /path/to/lego_v1.0.1_darwin_amd64.tar.gz`

3. Run lego:```shell /usr/local/bin/lego --domains yourmmdm.domain.com --email webmasteremail@domain.com --path /path/to/certs run'``

This should download your certs into a `certs` folder where you specified. Lego will automatically download the certs and the keys needed into two folders inside where you specififed called `accounts` and `certificates`.

If you specify the same path at time of renewal for your LetsEncrypt certs, all you will need to do restart MicroMDM.

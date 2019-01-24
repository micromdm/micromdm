# MicroMDM & Certbot

As of this writing, MicroMDM's usage of Lets Encrypt (a free SSL certificate service) is broken (see [this issue](https://github.com/micromdm/micromdm/issues/408)), therefore, **if you want MicroMDM to use HTTPS** (and yes, you do!), you will need to generate a certificate yourself and point MicroMDM to it.

Luckily, you can get one, for free, (given that you own a domain) using [Certbot](https://certbot.eff.org/), a python client for Let's Encrypt.

## Scope

In this article we will:

1. Install Certbot
2. Generate your certificate
3. Run MicroMDM using the certificate.

_The steps in this guide use Ubuntu 18.04 LTS as a reference but should be fairly easy to adapt to other setups_

## Install Certbot

On Ubuntu, the cetbot team maintains a PPA so installation is pretty simple once you add it to your repository:

```shell
sudo apt-get update
sudo apt-get install software-properties-common
sudo add-apt-repository ppa:certbot/certbot
sudo apt-get update
```

Once done, this will install certbot on your box:

```shell
sudo apt-get install certbot
```

To verify your installation simply type the following (it will list the certificates installed by `certbot`):

```shell
$ sudo certbot certificates
```

You should see something like this:

```
Saving debug log to /var/log/letsencrypt/letsencrypt.log
-------------------------------------------------------------------------------
No certs found.
-------------------------------------------------------------------------------
```

## Generate your certificate

In **this specific setup**, we are going to use the `--standalone` option to generate the certificate.

**Important!** For this to work, your domain *must resolve to the IP address where you are generating the cert, otherwise you'll see an error along these lines:
> `http://<yourdomain>/.well-known/acme-challenge/-K1FkTKDIEvjapSc_2aD2lICYQXvr6vKigQnKbmwqxE:
  Connection refused`.

Once that is set up, simply run:

```shell
$ sudo certbot certonly \
  --standalone \
  -d $yourdomain \
  -d $www.yourdomain
```

Use the `-d` flag to add any additional domains.

If all goes well you'll see:
```
IMPORTANT NOTES:
 - Congratulations! Your certificate and chain have been saved at:
   /etc/letsencrypt/live/<yourdomain>/fullchain.pem
```

A brief look at the `/etc/letsencrypt/live/$yourdomain/` folder reveals a few files:

```shell
cert.pem
chain.pem
fullchain.pem
privkey.pem
```

We'll need the last 2 for MicroMDM.


## Run MicroMDM using the certificates

This is the easy part. Once you have the binary in your path just run this:

```shell
sudo micromdm serve \
  -server-url=https://<yourserverURL> \
  -tls-cert=/etc/letsencrypt/live/<yourdomain>/fullchain.pem \
  -tls-key=/etc/letsencrypt/live/<yourdomain>/privkey.pem \
  -filerepo=<your path to micromdm packages folder>
```

### References

- [Let's Encrypt](https://letsencrypt.org/)
- [Certbot](https://certbot.eff.org/)

MicroMDM is a Go binary/Web Service which has the following deployment requirements/recommendations.

# Operating System

MicroMDM can be cross-compiled to a number of operating systems, including windows. The officially supported binaries built for release are for macOS and linux. macOS is recommended for development/testing and linux is recommended for a production deployment.

Official builds are available on
https://github.com/micromdm/micromdm/releases

MicroMDM is compatible with any major cloud provider's instances: AWS, GCP, Azure, DigitalOcean.
DigitalOcean is a great cloud provider to start with as they offer linux instances for as low as $5/mo.

# Certificates/Ports

By default, MicroMDM will bind to ports 80 and 443. This configuration, combined with a valid DNS name allows MicroMDM to automatically request and renew a certificate for your domain from Let's Encrypt.
The best supported configuration is:

 - A real domain(must resolve to a public instance).
 - Process should be allowed to bind to ports 80 and 443.

MicroMDM configures recommended [HTTP timeouts](https://blog.cloudflare.com/exposing-go-on-the-internet/) and enforces a [strict set of ciphersuites](https://wiki.mozilla.org/Security/Server_Side_TLS#Modern_compatibility), making it OK to expose on the internet.

Port 80 is only used to enforce a redirect from 80 -> 443 on the web server.

Other configurations, [like running MicroMDM with self signed certificates](https://github.com/micromdm/micromdm/wiki/Using-MicroMDM-with-self-signed-certificates) are also possible.

# Outgoing connections

MicroMDM makes a number of outgoing connections to Apple services, including `mdmenrollment.apple.com` and the APNS(push notification) service. These outgoing connections require an outbound connections on port 443.

# Usage

MicroMDM, like other web services will require more resources with usage. The more devices you enroll, the bigger the requirements. However, something that makes MDM unique from other servers is that devices will not check in with the server unless they receive a push notification. This makes the server usage mostly idle, unless a new MDM command (ex Installing a profile) is queued for a specific device.

If you're using MicroMDM mainly for bootstrapping macOS devices, the server will have relatively low utilization outside of your enrollment period.


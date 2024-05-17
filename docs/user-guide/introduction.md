# Getting Help

The best place to get help with setting up and ongoing maintenance of MicroMDM is the MacAdmins Slack. Join by getting an [invitation here](https://www.macadmins.org/).

Once you join Slack, the following channels will be useful.
- `#micromdm` for MicroMDM specific questions.
- `#mdm` and `#dep` for generic questions about MDM and deployment programs.

For defects and feature requests, please [open an issue](https://github.com/micromdm/micromdm/issues/new).

# Not a product!

MicroMDM is not a full featured device management product. The mission of MicroMDM is to enable a secure and scalable MDM deployment for Apple Devices, and expose the full set of Apple MDM commands and responses through an API.  
But it is more correct to think of MicroMDM as a lower level dependency for one or more products, not a solution that lives on its own. 

For example, MicroMDM has no high level options for configuration profiles. It accepts an already composed `mobileconfig` file and queues it for a single device at a time. Device Management products often have built-in support for signing profiles or pushing them to a chosen device group. MicroMDM expects those features to exist in a higher level companion service. 

MicroMDM has no Web UI. 

MicroMDM can enable disk encryption and escrow a secret, but it has no option for storing that secret on the server. 

Dynamic enrollment / user authentication workflows belong in an external service. MicroMDM serves the exact same enrollment profile at the `/mdm/enroll` server endpoint for every device. It also does not care about how the enrollment is protected from unauthorized devices. The number of possible workflows are infinite, and the recommendation is to point the devices at a separate URL for serving the enrollment profile. 

As you see, MicroMDM itself lacks many features that are usually present in device management products. But it also exposes a low level API that would allow an organization to build a product that is highly custom to ones environment. Over time, the community will likely share solutions that depend on MicroMDM and expose higher level workflows. 

*Before using MicroMDM, consider that there are a number of alternative, commercial products like Airwatch, SimpleMDM and Fleetsmith that might already fit your needs. Importantly, the companies behind these products employ developers dedicated to the development of their products. 
It is probably not a good idea to deploy MicroMDM as a cost cutting option, as the money you save will likely go towards hiring emplyees with the required domain knowledge and development expertise.*

# Apple Requirements

If you've decided to run an instance of MicroMDM in your organization, there are a few Apple specific requirements you need to meet.  
First, you need to [enroll](https://developer.apple.com/programs/enterprise/enroll/) your organization in the Apple Developer Enterprise Portal. Enrolling costs $299/year and requires that your organization have a [DUNS](https://en.wikipedia.org/wiki/Data_Universal_Numbering_System) number.  
Once signed up, or during the verification process in the first step, you need to ask Apple to enable the `MDM CSR` option. This option enables the signing of the APNS Push Certificate. The MDM CSR is typically reserved for commercial vendors, but Apple should enable it for you once you specify that you intend to use it for managing your company owned devices.

Finally, familiarize yourself with the [education](https://www.apple.com/education/it/) or [business](https://www.apple.com/business/it/) programs and enroll in Apple School/Business Manager(ABM). While MicroMDM does not require that you use the [deployment programs](https://support.apple.com/en-ca/HT204142) to enroll your devices, this is an increasingly popular option for enterprise deployments.

# Requirements for running MicroMDM

MicroMDM is a web server which handles connections from Apple Devices and exposes an authenticated web API for operators and other services to interact with the enrolled devices. If configured to do so, the server also syncs device records and profiles with ABM.  

The web server component is built into the `micromdm` binary, and launched with `micromdm serve [flags]`.  
For convenience, a second binary `mdmctl` is also provided and can be used to interact with the API. It does not implement _all_ the API endpoints, but should support most actions an administrator needs.
Notably, `mdmctl` does not need to beinstalled on the same server as `micromdm`. It is intended to be used by the server operator, and will be most commonly run on the operator's laptop. 

## Hardware Requirements

`micromdm` is a native Go binary and can be run on any hardware/VM/container environment. The release versions are built for `darwin` (for test environments on macOS) and `linux` for running on a server. Of course, you can compile for any platform that Go [supports](https://github.com/golang/go/wiki/MinimumRequirements). 

Currently, `micromdm` does not use a distributed database such as PostgreSQL and requires a persistent disk to be available. Because of this, only a single `micromdm` process can run at once, and high-availability setups are not possible. However, the database can be backed up/replicated while the server is running, allowing for failover with minimal downtime.  
If you need a more horizontally scalable open source MDM solution, take a look at [NanoMDM](https://github.com/micromdm/nanomdm).

Unlike many other services, once enrolled, devices do not maintain communication until an [APNS](https://developer.apple.com/library/archive/documentation/NetworkingInternet/Conceptual/RemoteNotificationsPG/APNSOverview.html#//apple_ref/doc/uid/TP40008194-CH8-SW1) command is scheduled to ask the device to check in. This architecture, makes the default setup of MicroMDM have relatively low hardware/requirements.  
For installations with under 10,000 device enrollments the recommended VM size is the GCP `n1-standard-2` [instance type](https://cloud.google.com/compute/docs/machine-types) or similar. This is roughly equivalent to 2 vCPUs and 7.5GB of RAM. 

API usage of any kind(ex: actions like pushing a configuration profile to all enrolled devices) will require additional resources, and will need additional planning. The recommendation is to [follow monitoring SRE best practices](https://landing.google.com/sre/workbook/chapters/monitoring/) and adjust your resources according to usage. 

## Network, TLS, DNS, APNS etc 

The best supported deployment configuration is to run `micromdm` on a internet accessible cloud VM, with an internet routable domain and TLS certificates from a public certificate authority like Let's Encrypt.  
Other options like using self signed certificates, private networks or proxies are possible, but mostly out of scope for the documentation in this guide. It's assumed that if you're deploying in an enterprise environment, you know how to deploy a webservice already. It should be roughly equivalent to deploying common software like nginx, wordpress or a Rails application. Users are encouraged to share their experiences in various enterprise environments on the Wiki, personal blogs or by speaking at conferences. 

You will need to allow the server to bind to port `443` for incoming TLS connections from devices and the admin API. This port is configurable if you'd like to use a different one. When binding to port `443`, it's encouraged to also allow connections on port `80`, which allows MicroMDM to redirect `http://` to `https://`. 

Your firewall MUST allow incoming connections to the ports MicroMDM is bound to(`443` and `80` by default). Your firewall MUST also allow outgoing connections to port `443`. This is necessary for sending messages over APNS. 

You should use a dedicated domain or subdomain for MicroMDM. One of the MDM requirements is that the server URL devices use does not change for the lifecycle of an enrollement. 

You will need to have a trusted TLS certificate in order for your devices to connect to MicroMDM. The certificate must be valid for the domain you're using with MicroMDM and can be terminated either by `micromdm` directly or any proxy running on the edge. 
Let's Encrypt is a great option for a CA, but any other will also work. 
The certificate file should be encoded as PEM and include the full chain:

```
-----BEGIN CERTIFICATE----- 
(Your leaf TLS certificate: your_domain_name.crt) 
-----END CERTIFICATE----- 
-----BEGIN CERTIFICATE----- 
(Your Intermediate certificate: TrustedCA.crt) 
-----END CERTIFICATE----- 
-----BEGIN CERTIFICATE----- 
(Your Root certificate: TrustedRoot.crt) 
-----END CERTIFICATE-----
```

Finally, your devices need to have access to APNS, a service run by Apple. This is not negotiable if you plan on using MDM.   
For a full description of APNS requirements, see https://support.apple.com/en-us/HT203609.
Apple also documents network requirements for devices connecting to the DEP services: https://support.apple.com/en-us/HT207516

# Operating

This document provided a high level description of various requirements to deploy MicroMDM. The rest of the user guide will provide detailed descriptions for how to operate the service and use the API. 

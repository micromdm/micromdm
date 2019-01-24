# OTA Enrollment

MicroMDM supports a method of profile delivery to the device called [OTA or Over-The-Air configuration or enrollment](https://developer.apple.com/library/archive/documentation/NetworkingInternet/Conceptual/iPhoneOTAConfiguration). OTA enrollment is the web-based profile enrollment technology that existed before MDM. Apple describes it as:

> iOS over-the-air enrollment and configuration provides an automated way to configure devices securely within the enterprise. This process provides IT with assurance that only trusted users are accessing corporate services and that their devices are properly configured to comply with established policies. Because configuration profiles can be both encrypted and locked, the settings cannot be removed, altered, or shared with others.

Technically speaking the [OTA process is a 3-phase HTTP request process](https://developer.apple.com/library/archive/documentation/NetworkingInternet/Conceptual/iPhoneOTAConfiguration/OTASecurity/OTASecurity.html) with its own separate SCEP certificate/enrollment step (which is separate from the MDM enrollment SCEP device identity certificate). OTA enrollment works for both macOS and iOS, despite the above blurb.

For MicroMDM the traditional way to get to the MDM enrollment profile is to visit the URL of a MicroMDM server here at:

* <https://mdm.example.org/mdm/enroll>

This will download the `.mobileconfig` enrollment profile which can then be double-clicked on, or used with `profiles -I`. For an OTA request you visit the OTA phase 1 URL:

* <https://mdm.example.org/ota/enroll>

## Why not just use the MDM enrollment endpoint

The main reason to pursue OTA enrollment over the plain `.mobileconfig` profile delivery is additional security. However it's important to note that _some_ of the security features of the OTA process are addressed when proper TLS is configured for MDM servers (which is mandatory for MDM). The OTA protocol was intended to work over non-HTTPS connections and additional measures for that situation are suggested (but not required). MicroMDM insists on TLS as well so some of those features for OTA aren't fully used.

Here are some specific reasons why folks would want to use OTA enrollment:

### OTA Phase 2 & 3 are signed by Apple-device-only certificates

Just like DEP device enrollment requests OTA Phase 2 & 3 requests are PKCS7-signed by an Apple device certificate, signed by an Apple device CA. This is a certificate that the device has that is signed _only_ by Apple. Both the MDM and OTA specifications guarantee that we can compare the signers and signing chain against a known Apple CA. Practically this means we can be reasonably sure that _only_ Apple devices are talking to to our endpoints.

_Note:_ one of the CA certificates mentioned here is actually expired. The Apple docs explicitly say to trust this expired CA and check that the signers match despite this expired CA certificate.

### OTA Phase 2 & 3 provide device details _before_ the MDM enrollment profile is delivered

Details like serial number, product, UDID, IMEI, user name, MAC address, etc. are all provided _before_ the MDM profile is delivered. This gives an opportunity to reject any further communication with the device based on its attributes. Or provides the ability to customize the MDM enrollment profile based on the device. This is also opportunity to keep state or tracking on things like UDID for an inventory system. One of the primary things the Apple OTA docs tout is tying the OTA process to a directory service to be able to validate e.g. users against a database.

MicroMDM doesn't currently use this information (but does parse it).

### Encrypt configuration profile to OTA SCEP (device-specific) certificate

The additional SCEP process gives the device an OTA profile that you can use to encrypt the delivered configuration (i.e. MDM enrollment) profile to. In this way the profile would only be readable by that particular device. This would be particularly useful for an OTA process that worked over non-encrypted HTTPS. However because MicroMDM requires TLS this is less of a concern.

### Protect the enrollment profile content

The three-step OTA process (protected by add'l SCEP and Apple-device CA) allows you to avoid the MDM enrollment profile from ever really being on disk or distributed elsewhere. This could be useful if you have any sensitive information in your enrollment profile. I.e. SCEP challenge passwords, URLs, CA certificates, embedded certificates, or any other sensitive configuration payloads in your enrollment profile.

### Additional hooks for user or organization authentication

One of the primary use cases Apple talks about is providing end-user authentication _before_ profile delivery in the OTA process. While Apple's example uses a simple web based authentication page, other techniques including OAuth, SAML, or anything else could potentially be used here. Technically those same techniques could be used to gate the "plain" enrollment profile, too, additional details could be verified about the device (i.e. in the OTA Phase 2 step, for example).

### MicroMDM OTA support

As above MicroMDM implements the OTA protocol but currently doesn't hook into any additional systems or make any additional decisions besides verifying the Apple CA. So with MicroMDM alone you get all the _inherent_ benefits of OTA, but some of the additional benefits listed above would need to have integrations or support code written for them.

If you'd like to see hooks, callbacks, pubsub events or the like for OTA enrollment, file issues or discuss on Slack #micromdm about what you're looking to do!

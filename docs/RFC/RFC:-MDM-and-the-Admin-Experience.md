# MDM and the admin experience

_For now this is just a bunch of ideas, if you just looking what this project is, don't take the following as a bill of functionalities. It's just brainstorming on ideal admin UX_

## Global MDM Behavior

Currents MDM all works the same: it offers a direct access to supported Payload and offer a more or less smart way to target the user.

None of the current MDM handle the complexity of profile management for the admin.

Properly managing settings with MDM imply deep knowledge of the client side.

For example, since the beginning it's possible to deploy multiple IMAP accounts via the same configuration profile. With 10.12 this feature has been broken, the mdmclient explicitly complain about it and refuse to install the payload because "only one IMAP account is allowed."

This kind of thing needs advanced skills and time for troubleshooting.

A good MDM should be able to check the configuration profile for compliance. Not only with Apple recommendations about optional/required fields. But also with real world experience. With a good MDM we should be able to mutualize the troubleshooting knowledge and handle most configuration errors in a pre-check scenario.

However, this feature must be a helper, not a guardian. In no case the MDM should refuse to save and push a configuration profile if the admin request it. The MDM should warn about a configuration is supposed not to be supported, but after a double check by the admin, it should push it.

For example, Profile Manager denies pushing a VPN settings when no value has been set as username. However this is perfectly supported by macOS.

In resumé, the MDM must be smart enough to help beginners to set things up, and smart enough to recognize when someone is more skilled.

## Targets

### Target selection

In all the ways offered to the admin to create a configuration profile, the target selection must be the same.

At the start, the profile should be set to manage configuration at the user or system level.

Then, deployment filter must be added. Filters must be regex compatible with creation wizard. Here is some example of filters needed (from real world scenario):

* users, group of users, group of users with subgroups
* devices, group of devices, group of devices with subgroups
* pattern matching on users and devices information fields (custom tags included)
* LDAP directory branch

Using the LDAP filter syntax could be a really good idea.

Everything related to users info must support every LDAP feature with no exception. Every single setup is different. And every single setup will use LDAP for user backend. And every single LDAP in this world is managed by human with different cultures and points of view.

So LDAP filter must not be limited. Attributes available must be extracted from the LDAP schema (specifying all supported classes for a user or a user group can be set during the initial configuration).

When admin inspects an item in the MDM (user, device, group…), the system must reflect current applied settings directly. Value must be presented as read only and payload variables must be replaced by real values if possible at current scope.

### Target exclusion

In some scenario, admin might want to exclude some users and computers from existing profiles.

To handle this scenario, exclusion rules must be supported and handled as a post action to the target lookup.

MDM should look for targets, then filter excluded items.

Excluded items must be managed from the item to exclude and reflected to the initial configuration profile read only.

For example, we have device groups named "Workstation", with computer A and B inside. A configuration profile called "Workstation Restrictions" is set to target all computers inside the device group Workstation.

Looking for computer info, admin must be able to see applied restrictions in detail (not just a link to the original profile). And must be able to request a specific exclusion.

Looking at configuration profile, admin must be able to edit the payload, targets and see specific exclusions. Exclusion can be removed from the configuration profile and the item. But added only on the item.

Exclusion can act for users, devices and groups.

## Configuration Profiles

## Scope for Configuration Profiles

Configuration Profiles have targets like seen previously. Targets define for whom or what the configuration is applied.

Scope will specify at which level: user or system.

For example, admin would like to deploy VPN settings (system level) for all computers used by mobile users.

Another example, admin might want to deploy a specific restriction (at the system level) when IT team are logged in (user check-in extension protocol).

## Removal Condition

In addition to basic removal condition supported by Apple, MDM should be able to remove the configuration based on custom workflow.

With the same example as before, admin would like to keep VPN settings when VIP user logout (to avoid losing settings when users is offsite) and remove all custom IT settings when management task is done.

## Raw Configuration Profiles

For expert users, raw configuration profiles like we know today must be available. Simply display all settings possible according to the selected scope and let the admin do whatever needed.

As said before, alter when something uncommon is done (regarding the Apple doc) or unlikely to work (according to word of mouth) but don't deny anything in this mode.

## Assisted Configuration Profiles

Most MDM admin will have many things to handle. MDM is just a part of their job. Assisted Configuration Profile is needed to simplify their work.

For example, it could be e-mail related. Assisted scenario could help admin to manage settings for personal mailbox (pre-filed by MDM, authenticated by the user on the end device) and also shared mailboxes with password built-in the profile.

Difference with raw scenario is simple: MDM will automatically break settings according to best practices.

Now that we know 10.12 has issues with multiple IMAP account in the same payload, assisted setup must automatically each e-mail account in a dedicated profile, even if seen as a unique one by the admin UI.

Another example is when we want restriction applied at the same time as a e-mail account. For BYOD we want them in the same profile. Remove restriction, remove your mail. But for in-house managed devices with no admin rights for the end user, restrictions and e-mails must be automatically split in multiple profiles to avoid removing the e-mail account (and all local settings + cache) when admin will update the restrictions.

Assisted Configuration Profile can also help to manage restriction in terms of unitary feature access. For example, iCloud Storage and online store are forbidden. That's all the admin want to know. State for other settings in the same payload must be handled according to the default system behavior when it's not managed.

## Cross references

Configuration Profile must be able to reference items from other profile. This will include a depenedence that must be managed.

For example, AD computer certificate can be deployed for everyone and every computer (so two payloads, system and user scope).

Then, computer cert might be used for 802.1X at login window, user one after login. Also, computer cert and user cert can be both used for L2TP/IPSec VPN.

This imply pre-check to identify targeted items who aren't also targeted by the referenced profile.

## Override management

### Override alert

When configuration is saved, MDM must run a global pre-check on all items targeted by the new profile. Pre-check must also be triggered when MDM actually push new content to the device.

Conflicting items must be identified and reported to the admin. Basic conflict detection on native payload is simple to do: array will merge, single value will override. This is valid for almost all cases.

Even if not perfect this can be improved in the future in a smarter way.

### Manual override

Configuration Profile has multiple targets. And as seen in target management part, when admin look at targeted items, configuration profile must be shown with details.

Useful features needed here is to allow override on a subtarget directly (group or final item).

Simple examples. Restrictions are deployed for everyone in the company, it's a mass scenario. But for some VIP users we want the exact same restrictions except that we will allow access to the App Store.

### Automatic Merge for Assisted Configuration Profiles

When using Assisted Configuration Profiles, admin might deny App Store for everyone and enforce Firewall for some. Since this belongs to the same payload, the MDM must handle the merge by itself.

### Automatic override

All profiles (assisted and raw) must have a override settings with the following kinds of values:

* must be installed
* automatic merge with precedence

Automatic merge will handle conflicting state for profiles applied on the same scope by applying to follow rules (coming from the good old times): closer to the user get better priority

If we take the Dock position on the screen, preference will be in this order:

1. user settings
2. group settings (closer in terms of the group inheritence is better)
3. device
4. group devices (closer in terms of the group inheritence is better)
5. directory branch (closer is better)

If payload is conflicting (same original target, no the same group tree branch, etc.), merge precedence is used to remove one of it.

If merge precedence is the same, both profiles are pushed and a warning is issued.

When automatic override rules remove a profile from the push list, a notice is issued explaining why it has been ignored and in benefit of what.

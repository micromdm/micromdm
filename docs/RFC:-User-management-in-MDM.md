_For now this is just a bunch of ideas, if you just looking what this project is, don't take the following as a bill of functionalities. It's just brainstorming_

# Why handle users and groups?

Nowadays, iOS like macOS are multiuser devices. Maybe not for everyone, but they are.

macOS is multiuser since the first version, and support networks user authentication to allow multiple users from the same company to connect on a shared devices.

iOS is at this time multiuser for school only, with support for managed apple ID authentication across all shared iPad the school own.

This with the BYOD scenario where you want to link an enrolled device to a specific user justify the need of user management and draw the general requirement for it.

# Sources of users and groups (and classrooms)

Trying to make an abstract of this use case can be interesting to understand how a user and group management can be implemented.

If we take the most complete scenario, the school one, we will have:
* users coming from Directory Service, Apple School Manager, and School Information System;
* groups coming from Directory Service;
* classrooms coming from School Information System.

The most interesting part to notice if for users. DS, ASM and SIS will talk about the same real user, but in different ways.

The MDM must be able to uniquely identify a user across all sources.

This mean, user record in the MDM DB must be references to users records across all users sources. And one user will have records in multiple sources. So a unique ID in the user schema of each source must be selected.

For directory services it should be LDAP ObjectUUID if it's LDAP, but an open architecture allowing JSON/REST directory service in the future could be interesting. For Apple School Manager, we also have UUID (different one), and for SIS, we can't know what it can be (maybe a student number, maybe nothing).

Common work will be with ASM and DS. Those two needs real and persistent unique ID.

SIS work will be more like rules for searches across other sources and will be used when classrooms are imported.

This research part will also be useful to create the initial link between MDM DB, ASM and DS. Maybe the Managed Apple ID can be linked to the Universal Principal Name. Maybe e-mail will be the same. Maybe Managed Apple ID UUID will be the student number available in employee numbers DS field…

The rule for the link between all users can't be crafted without knowledge of the DS and ASM customer config. So it has to be open for customization even if sample settings for common scenarios can be provided.

Group coming from the DS need to be also managed and refreshed for each user each time the user or its associated device check in. Nested groups need to be supported too.

Again, group representation in the DS is impossible to predict even if we can provide sample settings for AD and OD. For examples, groups can represent members with uid (OD) or full DN (AD), and nested groups as common members (AD) or with specific fields matching group UUID (OD). And it will be impossible to predict for custom LDAP schema.

MDM should have its own group representation like for users, with support of nested groups, and unique link between MDM group and DS group must be done based on real and persistent unique ID (and no, a DN isn't valid here).

DS integration should be done with an open API allowing simple extension in the future (to link to a JSON/REST directory service for example, or any kind of cloud-based directory services). The plugin architecture must force DS developers to translate its world to the MDM one. For example, plugin must force developers to provide UUID, username, aliases and e-mails for the user, and UUID, name, members and nested groups for groups. The work to link this to the DS field is delegated to the DS plugin developer.

Apple Directory Service implementation is a really good source of inspiration for this.

## Custom attributes

Custom attributes for users and groups are mandatory. Even if we think username, e-mail and groups can be enough, some case will require custom info coming from the directory service. For example, school might want to load username, e-mail, first name, last name, but also student number in classroom. Company might want a field to indicate if user is employee, contractor, provider, etc.

So, user and group global settings must allow the user to extend basics field with custom mapping.

## Multiple directory service

It's not needed for most scenario. But an interesting feature could be to support multiple DS since work it's almost done (to support DS + ASM).

Infinite number of DS is useful to allow augmented records when original directory service don't have all infos.

It's also useful in case of DS migration. This will allow IT to map current user records from old DS to the new one without loosing anything.

This feature is clearly not mandatory and will be used by really few people. But nice to have and keep in mind during development even if not implemented.

# How do I know who is where?

In BYOD scenario or DEP with personalized devices, the username linked to the device is provided as enrollment, so the user record must have a list of permanent device associations.

For shared devices, it's different. MDM protocol has two specific extensions, one for shared iOS devices and one for network users on macOS. Both extension does the same kind of task with different user ID: it allows the device to do a specific check-in with new push token when the user log in. So the MDM can now push user-specific payloads on this new push target. 

The system takes care for the MDM to associate new info to the right user profile and not mix it with others.

# Smart groups

MDM workflow can of course be assigned to a user or a group. But it might also need to be assigned to a search. So smart groups are needed.

For example, IT might want to push VPN profile to all users with apple-keyword LDAP field set to VPN value.
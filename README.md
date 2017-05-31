# MicroMDM - a devops friendly MDM server

[![CircleCI](https://circleci.com/gh/micromdm/micromdm/tree/master.svg?style=svg)](https://circleci.com/gh/micromdm/micromdm/tree/master)

MicroMDM is a Mobile Device Management server for Apple Devices currently focused on managing macOS.

## Introduction 

MDM is a large problem domain, and we are looking at how to solve many problems with device management and obstacles with existing solutions. However, at this stage of development we're focusing on a very specific use case which is current to many Mac Administrators: 

> When a device enrolls (through DEP or otherwise) we want to bootstrap some tools to manage the Mac. These are agents running on a mac like Munki or Chef/Puppet which will manage the software and configuration through the lifecycle of the Mac. 

MicroMDM is able to solve this common use-case today, and we're focusing on improving the user experience for administrators and developing related features.

MicroMDM is being actively developed but is **ready for you to start testing**.

To get started, see the [Quickstart](https://github.com/micromdm/micromdm/wiki/Quickstart) tutorial on the wiki. 

## Installing

Unless you're a developer, you probably want a binary you can run. [Download](https://github.com/micromdm/micromdm/releases/latest) the latest release, either for macOS (darwin) or linux. 

## Getting Help

The best place to get help is the `#micromdm` channel on the MacAdmins Slack team. Join us there by following getting an [invitation here](https://macadmins.herokuapp.com/).

## Helping out

Over the last year this project has gained a lot of interest from the community. We want to make it clear that at this point MicroMDM is a passion project, and is being developed by a [few fellow Mac Admins](https://github.com/micromdm/micromdm/graphs/contributors) on weekends and during conference hackathons. 

It would be great to get a few more developers to contribute, but at this point, there are _more_ important tasks than knowing how to code. Here is how you can help: 

- Read the documentation, install the tool and test MicroMDM.
- File bugs: https://github.com/micromdm/micromdm/issues  
- Participate in discussions. The `#micromdm` slack channel is best, but `#mdm` and `#dep` are a few relevant ones.   
- Edit the [project Wiki](https://github.com/micromdm/micromdm/wiki). The wiki page is open to anyone and you can make a lot of impact on the project by submitting additional documentation or designing proposals.
- Were you able to set up MicroMDM and enroll a few devices? Blogging about your experience. It can help others get started, or can help us figure out what we need to do better.

See the [CONTRIBUTING](CONTRIBUTING.md) page for additional info. 

## Design Goals

As mentioned in the introduction, a primary use case is bootstraping Macs. That's the short term (from now until WWDC). 
But there's a larger design goal we have in mind — what will differentiate MicroMDM from other vendor projects. Here it is in brief:

> MicroMDM aims to provide a declarative approach to device management. Too often vendor tools expect you to manage devices by filling out various forms in a web interface. While MicroMDM might have a web interface of it's own one day, the tool itself is inspired by popular DevOps processes like Configuration Management (chef/puppet/ansible, terraform) and orchestration frameworks (Kubernetes, Docker). 

For example here is a process of applying a DEP profile (not to be confused with an Apple Configuration Profile, of course):

```
$ mdmctl apply dep-profiles -template > /tmp/profile.json

$ mdmctl apply dep-profiles -f /tmp/profile.json
Defined DEP Profile with UUID 4B05B09E8AC7E7FC12C8F3338E099310

$ mdmctl get dep-profiles -f - -uuid=4B05B09E8AC7E7FC12C8F3338E099310
{
  "profile_name": "Test Profile",
  "url": "https://mdm.acmeinc.com/getconfig",
  "is_mdm_removable": true,
  "support_phone_number": "1-555-555-5555",
  "support_email_address": "org-email@example.com",
  "org_magic": "913FABBB-0032-4E13-9966-D6BBAC900331",
  "skip_setup_items": [
    "Registration",
    "AppleID",
    "TOS"
  ]
}
``` 

- Expose an API for developers and administratos. Today you can send MDM commands to the server using a RESTful API. We intend to make more processes scriptable.
- Provide a way for administrators to subscribe to events generated from the MDM interactions between client & server. MicroMDM works through a pubsub system at its core. For example, when a new device enrolls with the MDM server it doesn't record this in the database immediately, but instead creates an event which is sent on the message bus to other services that are listening. Today, this message bus is built in-memory, but we plan to expose the same hooks over the network, allowing developers to consume events in any language — not just Go.

[Here](https://docs.google.com/drawings/d/1B4w5xOmU-7D5pcW0kdiY7ia5fl7UnBfRpWR8KxzD1YI/edit?usp=sharing) is a slightly more in depth design overview of the pubsub system within MicroMDM.

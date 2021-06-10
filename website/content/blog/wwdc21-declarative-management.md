+++
title = "Tinkering with Declarative Management"
date = "2021-06-10T12:00:00-07:00"
tags = ["mdm", "declarative-management", "wwdc"]
author = "Jesse Peterson"
frontpage = true
+++

Amongst the many announcements at Apple's WWDC21 was [declarative management](https://developer.apple.com/news/?id=y3h32xgt). I would highly recommend [watching the WWDC session "Meet declarative device management"](https://developer.apple.com/videos/play/wwdc2021/10131/) to get a a general idea of what it is and how it works.

I wanted to tinker with the new capabilities and created some initial support for it and wanted to share how you can try this out, too. My hope is that, if you're inclined, you can get these tools and start playing with it.

Note this is a mostly technical how-to to get up and running and testing Declarative Management with NanoMDM. I hope to write more about the conceptual/paragidm changes DM brings as well more broad thinking in how to implement this protocol in a server. See the bottom for links to further discussion. Also note this is all really early stuff, based on beta software and documentation priovided by Apple and of course is all subject to major breaking changes.

<!--more-->

## Requirements

First you'll need to have some things to get this going.

* A running [NanoMDM](https://github.com/micromdm/nanomdm) instance (on a specific code branch, details below). This is probably the most difficult part to get going if you've never run an Open Source MDM solution before.
* A device running the iOS or iPadOS 15.0 pre-release
* Only [User Enrollments](https://support.apple.com/guide/deployment-reference-ios/user-enrollment-apdeb00576b2/web) are supported at this time with Declarative Management.
    * Apple Business Manager or Apple School Manager is required to create Managed Apple IDs
    * If you have ABM/ASM and Managed Apple IDs you can create a User Enrollment profile by adding the `ManagedAppleID` key to an MDM enrollment profile.
    * Check out the [WWDC20 videos from last year](https://developer.apple.com/videos/play/wwdc2019/303/) for more info on User Enrollment.

Putting it all together this means that this iOS/iPadOS 15.0 device needs to be enrolled into the the NanoMDM instance using User Enrollment. NanoMDM will let you know the style of enrollment of the device when it connects (as confirmation).

## Overview

If you look at the publically available documentation you can see that the [DeclarativeManagement check-in](https://developer.apple.com/documentation/devicemanagement/declarativemanagementrequest?changes=latest_minor) is, essentially, a wrapper for an HTTP service. It contains a body (called `Data`) and a "URL" endpoint (called `Endpoint`).

To that end NanoMDM simply treats it as such and synchronously dispatches to an external HTTP DM "server" with that HTTP data. We return the HTTP result to the MDM client in the body of the HTTP `DeclarativeManagement` check-in from this DM "server" — not as a Plist (like other check-ins) but as whatever the DM "server" responds with: usually nothing or JSON results.

Note that we also include the NanoMDM enrollment ID in the request in the `X-Enrollment-ID` header.

## Declarative management "server"

Before we setup NanoMDM we'll setup the toy Declarative Management server. This is a Flask python [app that lives in a GitHub Gist](https://gist.github.com/jessepeterson/5a633f627bfc23f74153add89aee07f1).

To get it going you might do something like:

```bash
$ python3 -m venv dm-venv
$ source dm-venv/bin/activate
$ pip install --upgrade pip
$ pip install flask
Collecting flask
[..snip..]
$ git clone https://gist.github.com/jessepeterson/5a633f627bfc23f74153add89aee07f1 dm
Cloning into 'dm'...
[..snip..]
$ cd dm
$ export FLASK_APP=app
$ export FLASK_DEBUG=1
$ flask run
[..snip..]
 * Debug mode: on
 * Running on http://127.0.0.1:5000/ (Press CTRL+C to quit)
[..snip..]
```

Note the "running on" line. We'll point NanoMDM at this server in a sec.

## NanoMDM branch

As this feature is very new this code only exists in a Pull Request against NanoMDM at the moment. So you'll have to check that specific branch and build it. That might look like this (this assumes you have the Go compiler set installed):

```bash
$ mkdir -p ~/go/src/github.com/micromdm/nanomdm
$ cd ~/go/src/github.com/micromdm
$ git clone git@github.com:jessepeterson/nanomdm.git
Cloning into 'nanomdm'...
[..snip..]
$ cd nanomdm
$ git checkout declarative-management
Branch declarative-management set up to track remote branch declarative-management from origin.
Switched to a new branch 'declarative-management'
$ make
[..snip..]
```

You should now have a binary of NanoMDM from the `declarative-management` branch in the `jessepeterson` fork of NanoMDM. This branch has a new switch: `-dm` which just takes a URL to the DM "server". This means our NanoMDM invocation would be something like:

```bash
./nanomdm-darwin-amd64 [normal-nanomdm-switches] -dm http://127.0.0.1:5000/ -dump
```

The `-dump` switch just tells NanoMDM to print out everything for us so we can inspect how the DM command runs.

## Activating declarative management

Once you've done all this, then you can active Declarative Management by sending your device the `DeclarativeManagement` command. With NanoMDM that looks like (from a check-out of NanoMDM):

```bash
./tools/cmdr.py command DeclarativeManagement | curl -v -T - -u nanomdm:api-key '[::1]:9000/v1/enqueue/B367CA4B-A874-DDA6-8168-A6295B54C7DF'
```

At this point assuming everything is working (and it is a lot of moving parts) you should see the Check-in requests for DeclarativeManagement. 

```bash
2021/06/10 12:51:05 level=info service=nanomdm msg=DeclarativeManagement id=B367CA4B-A874-DDA6-8168-A6295B54C7DF type=User Enrollment (Device) endpoint=declaration-items
```

Here we see that NanoMDM printed the enrollment ID, the check-in msg (DeclarativeManagement), it's type (User Enrollment) and the DM endpoint (declaration-items).

Even cooler this forwared the request over to http://127.0.0.1:5000/declaration-items for processing by appending the endpoint to the DM URL prefix. Our flask should have received this and printed:

```bash
127.0.0.1 - - [10/Jun/2021 12:51:05] "GET /declaration-items HTTP/1.1" 200 -
```

What you should then start seeing is the device doing various things with the actual Declaration Management protocol:

* `declaration-items` endpoint fetches all of the declarations.
* After requesting all of the declaration-tiems and seeing if they've changed (the first all of them will be new) it will start requesting the individual declarations
* `declaration/<type>/<uuid>` endpoints are the actual declarations getting fetched.
* The example/initial declarations are all specified directly in the server. There's an example set with a couple configurations, an activation, and a management declaration.
* Importantly you'll also see **status updates**. The device will proactively tell you things about its management state (which is amazing, in general). In these cases the status updates are reactionary to the declarations being fetched, but the concept is the same.
* `status` endpoint will start dumping a lot of JSON for any of the declarations, subscriptions, or other data that the status endpoint supports.

Here's an example of what this a status update looks like for changing declartions:

```json
{
  "StatusItems" : {
    "management" : {
      "declarations" : {
        "activations" : [
          {
            "active" : true,
            "identifier" : "49F6F16A-70EB-4A89-B092-465FAEC5E550",
            "valid" : "valid",
            "server-token" : "c306a35e-438c-55ca-af4c-e847d7bda0f9"
          }
        ],
        "configurations" : [
          {
            "reasons" : [
              {
                "details" : {
                  "Identifier" : "0FCD2F56-D5BC-48EA-B98D-E0CCC0C6F9E0",
                  "ServerToken" : "539f34e8-74f5-53d9-b986-374d1c4dedbc"
                },
                "description" : "Configuration (0FCD2F56-D5BC-48EA-B98D-E0CCC0C6F9E0:539f34e8-74f5-53d9-b986-374d1c4dedbc) is missing state.",
                "code" : "Error.MissingState"
              }
            ],
            "active" : false,
            "identifier" : "0FCD2F56-D5BC-48EA-B98D-E0CCC0C6F9E0",
            "valid" : "valid",
            "server-token" : "539f34e8-74f5-53d9-b986-374d1c4dedbc"
          },
          {
            "active" : true,
            "identifier" : "0FCD2F56-D5BC-48EA-B98D-E0CCC0C6F9E0",
            "valid" : "valid",
            "server-token" : "93665bad-af46-5391-b800-5483a4b07ebd"
          },
          {
            "active" : true,
            "identifier" : "85B5130A-4D0D-462B-AA0D-0C3B6630E5AA",
            "valid" : "unknown",
            "server-token" : "936e16fd-725f-5534-be1d-4e822ed78264"
          }
        ],
        "assets" : [

        ],
        "management" : [
          {
            "active" : false,
            "identifier" : "AF0E633E-7ADB-4B2A-A48C-1B20AA271D08",
            "valid" : "valid",
            "server-token" : "456fc8c6-85e9-52e9-b93d-eed694a4f487"
          }
        ]
      }
    }
  },
  "Errors" : [

  ]
}
```

In this example the `Error.MissingState` is telling us (I think) that one of the declarations is out of date and was retrieved. We'll get follow-up status updates letting us know each configuration was eventually "valid" or "invalid", it seems, too.

To get updated/changed items make some trivial change to one of the declarations in the toy server then send another `DeclarativeManagement` command. The device will check-in, fetch the declaration items, compute any changes, request the changed declarations, try to apply them, and send back statuses.

Note the toy server writes out each status to disk with the enrollment ID and timestamp in the filename.

## Feedback & Future

I created a [GH Discussion topic under the main MicroMDM site](https://github.com/micromdm/micromdm/discussions/759) for general discussion about open source MDM and Declarative Management. Hopefully we can coordinate discussion there. For any feedback about the above post (and related tools specifically) I would direct you to [the NanoMDM pull-request](https://github.com/micromdm/nanomdm/pull/24) adding this capability.

I hope this is enough to get you started tinkering with Declarative Managament and NanoMDM. Because of Declarative Management's limited scope at this time (iOS and iPadOS User Enrollments only) it is almost useless to the community involved in these open source tools (that is: macOS SysAdmin/IT management). But my hope is to get folks familiar ith the protocol and data models so we all can start thinking about future designs to start supporting declarative management when it comes to other platforms/environments.

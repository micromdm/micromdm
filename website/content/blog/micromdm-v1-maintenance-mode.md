+++
title = "MicroMDM v1 project maintenance mode"
date = "2025-06-20T14:10:00-07:00"
tags = ["micromdm", "nanomdm"]
author = "Jesse Peterson"
frontpage = true
+++

The MicroMDM project in its current form—what I refer to as MicroMDM “v1”—has effectively been in maintenance mode for years now. Development effort has been almost entirely focused on [NanoMDM](https://github.com/micromdm/nanohub) and the [“Nano”-suite of projects](https://github.com/micromdm?q=nano). This post serves to merely confirm and formalize the current situation: MicroMDM “v1” is indeed in maintenance mode.

<!--more-->

# What does that mean?

In short: no new feature work will happen in the MicroMDM “v1” codebase. Security and bug fixes will be on an available effort basis until the end of 2025. Further, the project is approaching an age where larger code dependencies may require refactors or other significant changes to integrate—even in the face of security issues. Such changes will be determined on a case-by-case basis against project maintainers' time, regardless of security necessity, even before the end of 2025.

I do want to be clear that the MicroMDM *organization* is still going strong. The Nano-suite of MDM projects and other organization projects are still maintained in all the same ways. This post is just a formal announcement of the current state of the MicroMDM v1 project code—that is: just the code that currently lives under [https://github.com/micromdm/micromdm](https://github.com/micromdm/micromdm) in the `main` branch. Beware we may soon change around project repositories, branches, or other code organization to reflect this maintenance mode status.

# And the future?

The official answer to what MicroMDM users should do now is to **migrate from MicroMDM to NanoMDM** and/or the Nano-suite of MDM components. They’re currently supported and generally maintained. The original [NanoMDM introductory post](https://micromdm.io/blog/introducing-nanomdm) mentions the [micro2nano](https://github.com/micromdm/micro2nano) tools which aid in transparently migrating existing MDM enrollments. I’ll also mention [this talk from last year’s MacDevOps::YVR](https://www.youtube.com/watch?v=ZBKidZVebAs) on my organization’s migration from MicroMDM to the Nano-suite of tools. Apple announced improvements to DEP/ADE notifications before, and during WWDC 2025 there are [new MDM migration capabilities](https://support.apple.com/guide/deployment/migrate-managed-devices-dep4acb2aa44/web), both of which assist MDM server migration, too. We hope to have more resources and documentation for people migrating in the future.

But what about longer term? The longer term *goal* hasn’t changed too much, only the amount of time available to get there—which has turned out much less than anticipated! The first [posts about MicroMDM v2 started in 2020](https://micromdm.io/blog/wwdc20-v2/). A bit after that NanoMDM was released and we [talked about MicroMDM v2](https://micromdm.io/blog/introducing-nanomdm/#micromdm-v2), too. With those goal in mind follow-on projects that implemented various parts of the Apple MDM ecosystem came next with [NanoDEP](https://github.com/micromdm/nanodep), [KMFDDM](https://github.com/jessepeterson/kmfddm), and [NanoCMD](https://github.com/micromdm/nanocmd). [NanoHUB](https://github.com/micromdm/nanohub) is the newest to join the club which unifies a few of the above projects together into a single low-level MDM server (and also has migration capabilities).

So, to say it plainly: the goal is still an eventual MicroMDM v2 which will use the Nano-suite of projects under the hood. Unfortunately there is still no timeline for that. Until that project materializes, please migrate to NanoMDM and/or the Nano-suite of components.

Finally I'd like to add a huge thank you to all of our users, contributors, partners, and fans over the years! Thanks!

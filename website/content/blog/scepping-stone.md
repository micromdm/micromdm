+++
frontpage = true
author = "Jesse Peterson"
tags = ["scep", "mda", "acme", "enrollment", "apple", "smallstep"]
date = "2024-07-12T21:54:50Z"
title = "SCEPping stone"
+++

[SCEP](https://en.wikipedia.org/wiki/Simple_Certificate_Enrollment_Protocol) has been a cornerstone of Apple MDM. While you've always technically been able to avoid SCEP in MDM, by directly embedding a device identity into an enrollment profile, SCEP has been the de facto (and most secure) way to get device identities for MDM authentication onto devices since the beginning of Apple MDM.

<!--more-->

To that end SCEP was a very early endeavor of the MicroMDM project by [Victor Vrantchen](https://github.com/groob) (the original creator of the MicroMDM project). At the time there wasn’t a Golang solution for SCEP at all so Victor set out to write one. The first commit to the [micromdm/scep](https://github.com/micromdm/scep) repository was in 2016; very near the early work on MicroMDM. The official MicroMDM server has embedded this SCEP project since the early days.

The micromdm/scep project is made up of three main components: a SCEP client CLI tool (`scepclient`), a SCEP server daemon (`scepserver`), and [a Golang library](https://pkg.go.dev/github.com/micromdm/scep/v2@v2.2.0/scep) for SCEP protocol primitives (which the other two components use). The project has been included as a [Debian package](https://packages.debian.org/stable/golang/scep) for years now. The tools and libraries are in use by other open source projects, large corporations (e.g. internally), and in commercial products—in some cases unrelated to MDM, even. We’ll mention a specific product below, in fact.

This SCEP project originated from a need to support Apple MDM enrollments. These days, though, with the introduction of support for [ACME](https://en.wikipedia.org/wiki/Automatic_Certificate_Management_Environment) and the [device-attest-01 challenge type](https://datatracker.ietf.org/doc/html/draft-acme-device-attest-01) by Apple (which they call [Managed Device Attestation](https://support.apple.com/guide/deployment/managed-device-attestation-dep28afbde6a/web)) MDM servers should be moving in this direction for certificate issuance for MDM enrollment.

If you follow the ACME (and especially `device-attest-01`/MDA) space you’re probably aware of [Smallstep](https://smallstep.com) and [step-ca](https://smallstep.com/certificates/)—an [open source](https://github.com/smallstep/certificates) Certificate Authority featuring ACME. The step-ca server had very early support for the `device-attest-01` challenge using the ACME protocol—including Apple’s MDA. This was based on the work of [Brandon Weeks](https://github.com/brandonweeks), the primary author of the IETF `device-attest-01` extension to ACME, who [integrated his work](https://github.com/brandonweeks/acme-device-attest-demo) on top of Smallstep’s step-ca. Smallstep is no stranger to open source collaboration.

Smallstep’s solution is also based on Golang. We were excited to learn that Smallstep, a pioneer in the certificate and ACME space, turned to the micromdm/scep project when they extended their step-ca product with SCEP protocol support for certificate issuance. Smallstep has also worked with us and contributed to MicroMDM and related projects as well.

With their good work and the observation that they are clearly some very smart folks and squarely in the certificate issuance space we started talking together about project ownership. The micromdm/scep project, while being popular and working well, hasn’t been very well maintained. Partly because it works decently enough as is, but mostly because the core focus of the MicroMDM projects are, well, Mobile Device Management, and not Certificate Authorities.

After further discussion it made sense to transfer ownership of the core SCEP protocol library over to Smallstep to maintain. So with enthusiastic encouragement from us, the `scep` Go package in the micromdm/scep repository was forked to [smallstep/scep](https://github.com/smallstep/scep). Separately, but relatedly, Smallstep also took on the work of forking the popular [Golang PKCS#7 library](https://github.com/smallstep/pkcs7) from Mozilla which was basically unmaintained. PKCS#7 (a.k.a. CMS) is the basis of the SCEP protocol. And with the [micromdm/scep#233 pull request](https://github.com/micromdm/scep/pull/233) the micromdm/scep project will officially remove the `scep` Golang _package_ in favor of using the Smallstep fork. We will continue to maintain the SCEP client CLI tool and SCEP server daemon for now.

I’d like to thank Smallstep, [Mike Malone](https://smallstep.com/about/), and in particular [Herman Slatman](https://github.com/hslatman) for taking on forking and maintaining of this (and other) projects! Thanks, of course, goes to Victor Vrantchan for his original creation of the micromdm/scep project as well. All of the contributors and helpful folks along the way for their support and interest in micromdm/scep as well: thanks! If you’re looking for an ACME, Apple MDA, or even SCEP solution for certificate issuance in your environment I’d suggest taking a look at [Smallstep](https://smallstep.com) and [step-ca](https://smallstep.com/docs/step-ca/)!

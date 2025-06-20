# MicroMDM v1 - a devops friendly MDM server

> [!WARNING]
> The [MicroMDM "v1" project is in maintenance mode](https://micromdm.io/blog/micromdm-v1-maintenance-mode/) and support officially ends at the end of 2025. You are encouraged to explore or migrate to [NanoMDM](https://github.com/micromdm/nanomdm) and other projects in the Nano-suite of projects.

[![CI](https://github.com/micromdm/micromdm/actions/workflows/CI.yml/badge.svg)](https://github.com/micromdm/micromdm/actions/workflows/CI.yml)

MicroMDM v1 is a Mobile Device Management server for Apple Devices, focused on giving you all the power through an API. 

# User Guide

- [Introduction](docs/user-guide/introduction.md)  
Requirements and other information you should know before using MicroMDM.

- [Quickstart](docs/user-guide/quickstart.md)  
A quick guide to get MicroMDM up and running. 

- [Enrolling Devices](docs/user-guide/enrolling-devices.md)  
Describes customizing the enrollment profile and the options available to get the profile installed on a device. Covers DEP provisioning as well as manual profile installs. 

- [API and Webhooks](docs/user-guide/api-and-webhooks.md)   
High level overview of the API used for scheduling device actions and processing the responses.

# Developer Guide

To help with development, start by reading the [CONTRIBUTING](./CONTRIBUTING.md) document, which has relevant resources. 

For a local development environment, or a demo setup, the [ngrok guide](./tools/ngrok/README.md), is the best resource to get something working.  

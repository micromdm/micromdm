# MicroMDM - a devops friendly MDM server 

[![CircleCI](https://circleci.com/gh/micromdm/micromdm/tree/main.svg?style=svg)](https://circleci.com/gh/micromdm/micromdm/tree/main)

MicroMDM is a Mobile Device Management server for Apple Devices, focused on giving you all the power through an API. 

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



# MySQL support
This version of MicroMDM can be run with Mysql support in place of default BoltDB database

## Test MicroMDM with Mysql
As root, Add your IP to your hosts file (MySQL don"t accept localhost connexions for users)
>echo $(hostname -I | cut -d" " -f1)" me.home.local" >>/etc/hosts

Get official MySQL docker image
>docker-compose -f docker-compose-dev.yaml pull db

Build micromdm with MySQL support (golang and alpine official docker images will be pulled)
>docker-compose -f docker-compose-dev.yaml build

Start MySQL and MicroMDM containers
>docker-compose -f docker-compose-dev.yaml up

Enrollement interface should now be available at this address (skip self signed certifcate warnings):
>https://me.home.local:3478/

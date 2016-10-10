mircomdm is a Mobile Device Management server for Apple Devices(primarily OS X macs).

While I intend to implement all the commands defined by Apple in the spec, the current focus is on implementing the features necessary to fit Apple's new(er) management tools (MDM, VPP, DEP) into existing enterprise environments.

This project now has a website with updated documentation - https://micromdm.io/



# Overview
**This repo is under heavy development. The current release is only for developers and expert users**

Current status

* Fetch devices from DEP
* Supports `InstallApplication` and `InstallProfile` commands
* Accepts a variety of other MDM payloads such as `OSUpdateStatus` and `DeviceInformation` but just dumps the response from the device to standard output.
* Push notifications are supported.
* Configuration profiles and applications can be grouped into a "workflow". The workflow can be assigned to a device.  
Currently the DEP enrollment step will check for a workflow but ignore it. I'll be adding this feature next.
* No SCEP/individual enrollment profiles yet. Need to have an enrollment profile on disk and pass it as a flag.

I set up a public [trello board](https://trello.com/b/js5u4DLV/micromdm-dev-board) to manage what is currently worked on and make notes.

# Getting started
Installation and configuration instructions will be maintained on the [website](https://micromdm.io/getting-started/#installation).


# Notes on architecture
* micromdm is an open source project written as an http server in [Go](https://golang.org/)
* deployed as a single binary.
* almost everything in the project is a separate library/service. `main` just wraps these together and provides configuration flags
* [PostgreSQL](http://www.postgresql.org/) for long lived data(devices, users, profiles, workflows)
* uses Redis to queue MDM Commands
* API driven - there will be an admin cli and a web ui, but the server itself is build as a RESTful API.
* exposes metrics data in [Prometheus](https://prometheus.io/) format.


# Workflows
An administrator can group a DEP enrollment profile, a list of applications and a list of configuration profiles into a workflow and assign the workflow to a device.  
If a device has an assigned workflow, `micromdm` will configure the device according to the workflow.
If you're familiar with Munki's [manifest](https://github.com/munki/munki/wiki/Manifests) feature, workflows work in a similar way.

# Build instructions

## If you know Go

1. `go get github.com/micromdm/micromdm`
2. `cd $GOPATH/src/github.com/micromdm/micromdm`
3. `glide install` install the necessary dependencies into /vendor folder
4. `go build` or `go install`

## If you're new to Go
Go is a bit different from other languages in its requirements for how it expects its programmers to organize Go code on a system.
First, Go expects you to choose a folder, called a workspace(you can name it anything you'd like). The path to this folder must always be set in an environment variable - `GOPATH`(example: `GOPATH=/Users/groob/code/go`)  
Your `GOPATH` must have thee subfolders - `bin`, `pkg` and `src`, and any code you create must live inside the `src` folder. It's also helpful to add `$GOPATH/bin` to your environment's `PATH` as that is where `go install` will place go binaries that you build.

A few helpful resources for getting started with Go.

* [Writing, building, installing, and testing Go code](https://www.youtube.com/watch?v=XCsL89YtqCs)
* [Resources for new Go programmers](http://dave.cheney.net/resources-for-new-go-programmers)
* [How I start](https://howistart.org/posts/go/1)
* [How to write Go code](https://golang.org/doc/code.html)
* [GOPATH - go wiki page](https://github.com/golang/go/wiki/GOPATH)

To build MicroMDM you will need to:  

1. Download and install [`Go`](https://golang.org/dl/)  
2. Install [`glide`](https://github.com/Masterminds/glide)
3. Set the `GOPATH` as explained above.
4. `mkdir -p $GOPATH/src/github.com/micromdm`
5. `git clone` the project into the above folder.  
The repo must always be in the folder `$GOPATH/src/github.com/micromdm/micromdm` even if you forked the project. Add a git remote to your fork.  
6. `glide install` The glide command will download and install all necessary dependencies for the project to compile.
7. `go build` or `go install`
8. File an issue or a pull request if the instructions were unclear.


## Makefile
The project has a Makefile and will build the project for you assuming you have `GOPATH` set correctly.
* run `make` to create a new build.
* `make deps` will install the necessary dependencies. after that you can use `go build`, `go test` etc.
* run `make docker` to build a docker container from the local source.  

## Docker container for Redis and PostgreSQL
If you want to run locally for testing/development, an easy way to run PostgreSQL and Redis is by using `docker-compose`
`docker-compose -f compose-pg.yml up`

## Dockerfiles for development and release.
* `Dockerfile` will build the latest release(by downloading the binaries.  
This is equivalent to `docker pull micromdm/micromdm:latest`

* `Dockerfile.dev` builds the latest version from the local source.
`docker build -f Dockerfile.dev -t micromdm .`

`docker pull micromdm/micromdm:dev` to get the latest version built from master.

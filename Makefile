all: build

.PHONY: build

ifndef ($(GOPATH))
	GOPATH = $(HOME)/go
endif

PATH := $(GOPATH)/bin:$(PATH)
VERSION = $(shell git describe --tags --always --dirty)
BRANCH = $(shell git rev-parse --abbrev-ref HEAD)
REVISION = $(shell git rev-parse HEAD)
REVSHORT = $(shell git rev-parse --short HEAD)
USER = $(shell whoami)
GOVERSION = $(shell go version | awk '{print $$3}')
NOW	= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
SHELL = /bin/bash

ifneq ($(OS), Windows_NT)
	CURRENT_PLATFORM = linux
	ifeq ($(shell uname), Darwin)
		SHELL := /bin/bash
		CURRENT_PLATFORM = darwin
	endif
else
	CURRENT_PLATFORM = windows
endif

BUILD_VERSION = "\
	-X github.com/micromdm/micromdm/vendor/github.com/micromdm/go4/version.appName=${APP_NAME} \
	-X github.com/micromdm/micromdm/vendor/github.com/micromdm/go4/version.version=${VERSION} \
	-X github.com/micromdm/micromdm/vendor/github.com/micromdm/go4/version.branch=${BRANCH} \
	-X github.com/micromdm/micromdm/vendor/github.com/micromdm/go4/version.buildUser=${USER} \
	-X github.com/micromdm/micromdm/vendor/github.com/micromdm/go4/version.buildDate=${NOW} \
	-X github.com/micromdm/micromdm/vendor/github.com/micromdm/go4/version.revision=${REVISION} \
	-X github.com/micromdm/micromdm/vendor/github.com/micromdm/go4/version.goVersion=${GOVERSION}"

WORKSPACE = ${GOPATH}/src/github.com/micromdm/micromdm
check-deps:
ifneq ($(shell test -e ${WORKSPACE}/Gopkg.lock && echo -n yes), yes)
	@echo "folder is clonded in the wrong place, copying to a Go Workspace"
	@echo "See: https://golang.org/doc/code.html#Workspaces"
	@git clone git@github.com:micromdm/micromdm ${WORKSPACE}
	@echo "cd to ${WORKSPACE} and run make deps again."
	@exit 1
endif
ifneq ($(shell pwd), $(WORKSPACE))
	@echo "cd to ${WORKSPACE} and run make deps again."
	@exit 1
endif

deps: check-deps
	go get -u github.com/golang/dep/...
	dep ensure -vendor-only

test:
	go test -cover -race -v $(shell go list ./... | grep -v /vendor/)

build: micromdm mdmctl

clean:
	rm -rf build/
	rm -f *.zip

.pre-build:
	mkdir -p build/darwin
	mkdir -p build/linux

INSTALL_STEPS := \
	install-mdmctl \
	install-micromdm

install-local: $(INSTALL_STEPS)

.pre-mdmctl:
	$(eval APP_NAME = mdmctl)

mdmctl: .pre-build .pre-mdmctl
	go build -i -o build/$(CURRENT_PLATFORM)/mdmctl -ldflags ${BUILD_VERSION} ./cmd/mdmctl

install-mdmctl: .pre-mdmctl
	go install -ldflags ${BUILD_VERSION} ./cmd/mdmctl

APP_NAME = micromdm

.pre-micromdm:
	$(eval APP_NAME = micromdm)

micromdm: .pre-build .pre-micromdm
	go build -i -o build/$(CURRENT_PLATFORM)/micromdm -ldflags ${BUILD_VERSION} ./

install-micromdm: .pre-micromdm
	go install -ldflags ${BUILD_VERSION}

xp-micromdm: .pre-build .pre-micromdm
	GOOS=darwin go build -i -o build/darwin/micromdm -ldflags ${BUILD_VERSION}
	GOOS=linux CGO_ENABLED=0 go build -i -o build/linux/micromdm  -ldflags ${BUILD_VERSION}

release-zip: xp-micromdm mdmctl
	zip -r micromdm_${VERSION}.zip build/

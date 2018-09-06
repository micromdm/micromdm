all: build

.PHONY: build

ifeq ($(GOPATH),)
	PATH := $(HOME)/go/bin:$(PATH)
else
	PATH := $(GOPATH)/bin:$(PATH)
endif

export GO111MODULE=on

VERSION = $(shell git describe --tags --always --dirty)
BRANCH = $(shell git rev-parse --abbrev-ref HEAD)
REVISION = $(shell git rev-parse HEAD)
REVSHORT = $(shell git rev-parse --short HEAD)
USER = $(shell whoami)
GOVERSION = $(shell go version | awk '{print $$3}')
NOW	= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
SHELL = /bin/bash

ifndef ($(DOCKER_IMAGE_NAME))
	DOCKER_IMAGE_NAME = micromdm/micromdm
endif
DOCKER_IMAGE_TAG = $(shell echo ${VERSION} | sed 's/^v//')

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
	-X github.com/micromdm/go4/version.appName=${APP_NAME} \
	-X github.com/micromdm/go4/version.version=${VERSION} \
	-X github.com/micromdm/go4/version.branch=${BRANCH} \
	-X github.com/micromdm/go4/version.buildUser=${USER} \
	-X github.com/micromdm/go4/version.buildDate=${NOW} \
	-X github.com/micromdm/go4/version.revision=${REVISION} \
	-X github.com/micromdm/go4/version.goVersion=${GOVERSION}"

gomodcheck: 
	@go help mod > /dev/null || (@echo micromdm requires Go version 1.11 or higher && exit 1)

deps: gomodcheck
	@go mod download

test:
	go test -cover -race ./...

build: micromdm mdmctl

clean:
	rm -rf build/
	rm -f *.zip

.pre-build:
	mkdir -p build/darwin
	mkdir -p build/linux

install-local: \
	install-mdmctl \
	install-micromdm

.pre-mdmctl:
	$(eval APP_NAME = mdmctl)

mdmctl: .pre-build .pre-mdmctl
	go build -i -o build/$(CURRENT_PLATFORM)/mdmctl -ldflags ${BUILD_VERSION} ./cmd/mdmctl

xp-mdmctl: .pre-build .pre-mdmctl
	GOOS=darwin go build -i -o build/darwin/mdmctl -ldflags ${BUILD_VERSION} ./cmd/mdmctl
	GOOS=linux CGO_ENABLED=0 go build -i -o build/linux/mdmctl  -ldflags ${BUILD_VERSION} ./cmd/mdmctl

install-mdmctl: .pre-mdmctl
	go install -ldflags ${BUILD_VERSION} ./cmd/mdmctl

APP_NAME = micromdm

.pre-micromdm:
	$(eval APP_NAME = micromdm)

micromdm: .pre-build .pre-micromdm
	go build -i -o build/$(CURRENT_PLATFORM)/micromdm -ldflags ${BUILD_VERSION} ./cmd/micromdm

install-micromdm: .pre-micromdm
	go install -ldflags ${BUILD_VERSION} ./cmd/micromdm

xp-micromdm: .pre-build .pre-micromdm
	GOOS=darwin go build -i -o build/darwin/micromdm -ldflags ${BUILD_VERSION} ./cmd/micromdm
	GOOS=linux CGO_ENABLED=0 go build -i -o build/linux/micromdm  -ldflags ${BUILD_VERSION} ./cmd/micromdm

release-zip: xp-micromdm xp-mdmctl
	zip -r micromdm_${VERSION}.zip build/

docker-build:
	GOOS=linux CGO_ENABLED=0 go build -o build/linux/micromdm  -ldflags ${BUILD_VERSION} ./cmd/micromdm
	docker build -t ${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG} .

docker-tag: docker-build
	docker tag ${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG} ${DOCKER_IMAGE_NAME}:latest

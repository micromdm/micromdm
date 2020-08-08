.PHONY: build

GO := go
WORKSPACE := $(dir $(shell ${GO} env GOMOD))
BUILD_DIR := $(WORKSPACE)/build

export GOBIN=$(BUILD_DIR)

VERSION = $(shell git describe --tags --always --dirty)
BRANCH = $(shell git rev-parse --abbrev-ref HEAD)
REVISION = $(shell git rev-parse HEAD)
REVSHORT = $(shell git rev-parse --short HEAD)
USER = $(shell whoami)
NOW	= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

BUILD_VERSION = "\
	-X micromdm.io/v2/pkg/version.appName=${APP_NAME} \
	-X micromdm.io/v2/pkg/version.version=${VERSION} \
	-X micromdm.io/v2/pkg/version.branch=${BRANCH} \
	-X micromdm.io/v2/pkg/version.buildUser=${USER} \
	-X micromdm.io/v2/pkg/version.buildDate=${NOW} \
	-X micromdm.io/v2/pkg/version.revision=${REVISION}"

download:
	@echo "Download dependencies"
	@${GO} mod download

install-tools: download
	@echo "Installing tools"
	@cat cmd/tools/tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % ${GO} install %

micromdm:
	$(eval APP_NAME = micromdm)
	${GO} install -race -ldflags ${BUILD_VERSION} ./cmd/micromdm

all: build

.PHONY: build

ifeq ($(GOPATH),)
	PATH := $(HOME)/go/bin:$(PATH)
else
	PATH := $(GOPATH)/bin:$(PATH)
endif

export GO111MODULE=on

#Specify a minimum version for macos otherwise notarization will fail
CGO_LDFLAGS=-mmacosx-version-min=10.12 

VERSION = $(shell git describe --tags --always --dirty)
BRANCH = $(shell git rev-parse --abbrev-ref HEAD)
REVISION = $(shell git rev-parse HEAD)
REVSHORT = $(shell git rev-parse --short HEAD)
USER = $(shell whoami)
GOVERSION = $(shell go version | awk '{print $$3}')
NOW	= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
SHELL = /bin/sh

ifndef ($(DOCKER_IMAGE_NAME))
	DOCKER_IMAGE_NAME = micromdm/micromdm
endif
DOCKER_IMAGE_TAG = $(shell echo ${VERSION} | sed 's/^v//')

ifneq ($(OS), Windows_NT)
	CURRENT_PLATFORM = linux
	ifeq ($(shell uname), Darwin)
		SHELL := /bin/sh
		CURRENT_PLATFORM = darwin
	endif
else
	CURRENT_PLATFORM = windows
endif

ifeq ($(PG_HOST),)
PG_HOST := localhost
endif

ifeq ($(MYSQL_HOST),)
MYSQL_HOST := 127.0.0.1
endif

BUILD_VERSION = "\
	-X github.com/micromdm/go4/version.appName=${APP_NAME} \
	-X github.com/micromdm/go4/version.version=${VERSION} \
	-X github.com/micromdm/go4/version.branch=${BRANCH} \
	-X github.com/micromdm/go4/version.buildUser=${USER} \
	-X github.com/micromdm/go4/version.buildDate=${NOW} \
	-X github.com/micromdm/go4/version.revision=${REVISION} \
	-X github.com/micromdm/go4/version.goVersion=${GOVERSION} \
	-X github.com/micromdm/micromdm/dep.version=${VERSION}"

gomodcheck:
	@go help mod > /dev/null || (@echo micromdm requires Go version 1.11 or higher && exit 1)

deps: gomodcheck
	@go mod download

test: db-migrate-test db-mysql-migrate-test
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
	go build -o build/$(CURRENT_PLATFORM)/mdmctl -ldflags ${BUILD_VERSION} ./cmd/mdmctl

xp-mdmctl: .pre-build .pre-mdmctl
	GOOS=darwin go build -o build/darwin/mdmctl -ldflags ${BUILD_VERSION} ./cmd/mdmctl
	GOOS=linux CGO_ENABLED=0 go build -o build/linux/mdmctl  -ldflags ${BUILD_VERSION} ./cmd/mdmctl

install-mdmctl: .pre-mdmctl
	go install -ldflags ${BUILD_VERSION} ./cmd/mdmctl

APP_NAME = micromdm

.pre-micromdm:
	$(eval APP_NAME = micromdm)

micromdm: .pre-build .pre-micromdm
	go build -o build/$(CURRENT_PLATFORM)/micromdm -ldflags ${BUILD_VERSION} ./cmd/micromdm

install-micromdm: .pre-micromdm
	go install -ldflags ${BUILD_VERSION} ./cmd/micromdm

xp-micromdm: .pre-build .pre-micromdm
	GOOS=darwin go build -o build/darwin/micromdm -ldflags ${BUILD_VERSION} ./cmd/micromdm
	GOOS=linux CGO_ENABLED=0 go build -o build/linux/micromdm  -ldflags ${BUILD_VERSION} ./cmd/micromdm

release-zip: xp-micromdm xp-mdmctl
	zip -r micromdm_${VERSION}.zip build/

docker-build:
	GOOS=linux CGO_ENABLED=0 go build -o build/linux/micromdm  -ldflags ${BUILD_VERSION} ./cmd/micromdm
	docker build -t ${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG} .

docker-tag: docker-build
	docker tag ${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG} ${DOCKER_IMAGE_NAME}:latest

ngrok:
	@./tools/ngrok/screen

docker-compose:
	docker-compose -f docker-compose-dev.yaml up -d





# PostgreSQL
db-psql-test:
	$(call psql_db,micromdm_test)

db-psql:
	$(call psql_db,micromdm)

define psql_db
	PGPASSWORD=micromdm psql --host=${PG_HOST} --port=5432 --username=micromdm --dbname=$(1)
endef

db-reset-test:
	$(call psql_exec,'DROP DATABASE IF EXISTS micromdm_test;')
	$(call psql_exec,'CREATE DATABASE micromdm_test;')

define psql_exec
	PGPASSWORD=micromdm psql --host=${PG_HOST} --port=5432 --username=micromdm -c $(1)
endef

db-migrate-test:
	$(call goose_up,micromdm_test)

db-migrate:
	$(call goose_up,micromdm)

define goose_up
	cd ./pg/migrations && goose postgres "host=${PG_HOST} port=5432 user=micromdm dbname=$(1) password=micromdm sslmode=disable" up
endef




# MySQL

# Make sure user can access db:
# mysql -u root -p
# mysql> CREATE USER 'micromdm'@'localhost' IDENTIFIED BY 'micromdm';
# mysql> GRANT ALL PRIVILEGES ON *.* TO 'micromdm'@'localhost' IDENTIFIED BY 'micromdm';
db-mysql-test:
	$(call mysql_db,micromdm_test)

db-mysql:
	$(call mysql_db,micromdm)

define mysql_db
	MYSQLPASSWORD=micromdm mysql --host=${MYSQL_HOST} --port=3306 --user=micromdm --password=micromdm $(1)
endef

db-mysql-reset-test:
	$(call mysql_exec,'SHOW DATABASES;')
	$(call mysql_exec,'DROP DATABASE IF EXISTS micromdm_test;')
	$(call mysql_exec,'CREATE DATABASE micromdm_test;')
	$(call mysql_exec,'SET sql_mode = "";')
	
define mysql_exec
	mysql --host=${MYSQL_HOST} --port=3306 --user=micromdm --password=micromdm --execute=$(1)
endef

db-mysql-migrate-test:
	$(call mysql-goose_up,micromdm_test)

db-mysql-migrate:
	$(call mysql-goose_up,micromdm)

define mysql-goose_up
	cd ./mysql/migrations && goose mysql "micromdm:micromdm@/$(1)?parseTime=true" up
endef

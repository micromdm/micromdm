#!/bin/bash

VERSION="0.1.2.0-dev"
NAME=micromdm

echo "Building $NAME version $VERSION"

mkdir -p build

build() {
  echo -n "=> $1-$2: "
  GOOS=$1 GOARCH=$2 CGO_ENABLED=0 go build -o build/$NAME-$1-$2 -ldflags "\
      -X github.com/micromdm/micromdm/version.version=${VERSION} \
      -X github.com/micromdm/micromdm/version.gitBranch=$(git rev-parse --abbrev-ref HEAD) \
      -X github.com/micromdm/micromdm/version.goVersion=$(go version | awk '{print $3}') \
      -X github.com/micromdm/micromdm/version.buildTime=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
      -X github.com/micromdm/micromdm/version.gitRev=$(git rev-parse HEAD)" ./main.go
  du -h build/$NAME-$1-$2
}

build "darwin" "amd64"
build "linux" "amd64"

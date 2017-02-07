#!/usr/bin/env bash

set -eo pipefail

usage() {
  base="$(basename "$0")"
  cat <<EOUSAGE
Usage: ${base} [args]
  -T,--test             : Run tests
  -B,--build            : Build a release
EOUSAGE
}

if [ $# -eq 0 ]; then
  usage
fi

# Flag parsing
while [[ $# -gt 0 ]]; do
  opt="$1"
  case "${opt}" in
    -B|--build)
      build=1
      shift
      ;;
    -T|--tests)
      tests=1
      shift
      ;;
    *)
      echo "Error: Unknown option: ${opt}"
      usage
      exit 1
      ;;
  esac
done

VERSION="0.1.0.1-dev"
NAME=micromdm

build_release() {
  echo -n "=> $1-$2: "
  GOGC=500 GOOS=$1 GOARCH=$2 CGO_ENABLED=0 go build -i -o build/$NAME-$1-$2 -ldflags "-X main.Version=$VERSION -X main.gitHash=`git rev-parse HEAD`" ./main.go
  du -h build/$NAME-$1-$2
}

build=${build:-0}
if [ ${build} -eq 1 ]; then
    mkdir -p build
    build_release "darwin" "amd64"
    build_release "linux" "amd64"
  exit 0
fi

tests=${tests:-0}
if [ ${tests} -eq 1 ]; then
    echo "not implemented"
  exit 0
fi

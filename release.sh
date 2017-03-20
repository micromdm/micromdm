VERSION="$(git describe --tags --always --dirty)"
NAME=micromdm
USER=$(whoami)
BRANCH=$(git rev-parse --abbrev-ref HEAD)
NOW=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
REVISION=$(git rev-parse HEAD)
GOVERSION=$(go version | awk '{print $3}')

echo "Building $NAME version $VERSION"

mkdir -p build

build() {
  echo -n "=> $1-$2: "
  GOOS=$1 GOARCH=$2 CGO_ENABLED=0 go build -o build/$NAME-$1-$2 -ldflags "\
      -X github.com/micromdm/micromdm/version.version=${VERSION}\
      -X github.com/micromdm/micromdm/version.branch=${BRANCH}\
      -X github.com/micromdm/micromdm/version.buildUser=${USER}\
      -X github.com/micromdm/micromdm/version.buildDate=${NOW}\
      -X github.com/micromdm/micromdm/version.revision=${REVISION}\
      -X github.com/micromdm/micromdm/version.goVersion=${GOVERSION}\
      " ./*.go
  du -h build/$NAME-$1-$2
}

build "darwin" "amd64"
build "linux" "amd64"

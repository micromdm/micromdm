VERSION="$(git describe --tags --always --dirty)"
USER=$(whoami)
BRANCH=$(git rev-parse --abbrev-ref HEAD)
NOW=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
REVISION=$(git rev-parse HEAD)
GOVERSION=$(go version | awk '{print $3}')

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
      " $3
  du -h build/$NAME-$1-$2
}

NAME=micromdm
echo "Building $NAME version $VERSION"
build "darwin" "amd64" "./*.go"
build "linux" "amd64" "./*.go"

NAME=mdmctl
echo "\nBuilding $NAME version $VERSION"
build "darwin" "amd64" "./cmd/mdmctl/*.go"
build "linux" "amd64" "./cmd/mdmctl/*.go"

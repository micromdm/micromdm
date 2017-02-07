FROM golang:1.7.1-alpine

RUN apk add --update \
        git \
		bash \
		musl-dev \
		openssl


COPY ./builder.sh /builder.sh

VOLUME   /go/src/github.com/micromdm/micromdm
WORKDIR  /go/src/github.com/micromdm/micromdm

ENTRYPOINT ["/builder.sh"]


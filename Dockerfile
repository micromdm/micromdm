FROM golang:latest as builder

WORKDIR /go/src/github.com/micromdm/micromdm/

ENV CGO_ENABLED=0 \
	GOARCH=amd64 \
	GOOS=linux

COPY . .

RUN make deps
RUN make


FROM alpine:latest

RUN apk --update add ca-certificates git
WORKDIR /go/src/app
COPY . .
ENV GO111MODULE on
RUN go mod download
RUN GOOS=linux CGO_ENABLED=0 go install ./cmd/micromdm
RUN mkdir /data; chmod 777 /data

COPY docker-entrypoint.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
EXPOSE 8080 8443
CMD ["micromdm", "serve"]

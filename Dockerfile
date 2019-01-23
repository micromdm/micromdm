FROM golang:1.11-alpine

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

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
COPY --from=builder /go/src/github.com/micromdm/micromdm/build/linux/micromdm .
COPY --from=builder /go/src/github.com/micromdm/micromdm/build/linux/mdmctl .

COPY docker-entrypoint.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
EXPOSE 8080 8443
CMD ["micromdm", "serve"]

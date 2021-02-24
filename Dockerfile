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
RUN apk add openssl

RUN mkdir repo
RUN mkdir /data; chmod 777 /data
COPY --from=builder /go/src/github.com/micromdm/micromdm/build/linux/micromdm /usr/bin/
COPY --from=builder /go/src/github.com/micromdm/micromdm/build/linux/mdmctl /usr/bin/

RUN DNSNAME=me.home.local;  (cat /etc/ssl/openssl.cnf ; printf "\n[SAN]\nsubjectAltName=DNS:$DNSNAME\n") | openssl req -new -newkey rsa:2048 -days 365 -nodes -x509 -sha256 -keyout server.key -out server.crt -subj "/CN=$DNSNAME" -reqexts SAN -extensions SAN -config /dev/stdin
RUN sh -c 'echo "127.0.0.1 me.home.local" >> /etc/hosts'

EXPOSE 8080 8443
CMD ["micromdm", "serve"]

FROM golang:1.11-alpine

RUN apk --update add ca-certificates git
WORKDIR /go/src/app
COPY . .
ENV GO111MODULE on
RUN go mod download
RUN GOOS=linux CGO_ENABLED=0 go install ./cmd/micromdm

EXPOSE 80 443
VOLUME ["/var/db/micromdm"]
CMD ["micromdm", "serve"]

FROM alpine:3.7

RUN apk --update add ca-certificates

COPY ./build/linux/micromdm /usr/bin

EXPOSE 80 443
VOLUME ["/var/db/micromdm"]
CMD ["micromdm", "serve"]

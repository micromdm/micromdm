FROM alpine:3.6

RUN apk --update add ca-certificates && \
	update-ca-certificates

COPY ./build/linux/micromdm /usr/bin

EXPOSE 80 443
VOLUME ["/var/db/micromdm"]
CMD ["micromdm", "serve"]

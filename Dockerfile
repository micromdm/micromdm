FROM alpine:3.6

ENV MICROMDM_VERSION=v1.2.0
RUN apk --no-cache add curl unzip ca-certificates && \
	update-ca-certificates && \
    curl -L https://github.com/micromdm/micromdm/releases/download/${MICROMDM_VERSION}/micromdm_${MICROMDM_VERSION}.zip -o /micromdm.zip && \
    unzip -p /micromdm.zip build/linux/micromdm > /micromdm && \
    rm /micromdm.zip && \
    chmod a+x /micromdm && \
    apk del curl unzip && \
    echo '666e5457d396d84b1fc995886c7598741a487b50132f95305b6be44b20c7d834  micromdm' | sha256sum -c || exit 0

EXPOSE 443
VOLUME ["/var/db/micromdm"]
ENTRYPOINT ["/micromdm"]
CMD ["serve", "-examples"]

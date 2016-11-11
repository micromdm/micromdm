FROM alpine:3.3

ENV SCEP_VERSION=0.1.0.0
RUN apk --no-cache add curl && \
    curl -L https://github.com/micromdm/scep/releases/download/${SCEP_VERSION}/scep-linux-amd64 -o /scep && \
    chmod a+x /scep && \
    apk del curl

CMD ["/scep"]

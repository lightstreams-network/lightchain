FROM alpine:3.7
MAINTAINER Gabriel Garrido <gabriel@lightstreams.io>

RUN apk update && \
    apk upgrade && \
    apk --no-cache add curl jq file

VOLUME [ /lightstreams ]
WORKDIR /lightstreams

EXPOSE 26656 26657

ENTRYPOINT ["/usr/bin/wrapper.sh"]
CMD ["node", "--proxy_app", "kvstore"]
STOPSIGNAL SIGTERM

COPY wrapper.sh /usr/bin/wrapper.sh


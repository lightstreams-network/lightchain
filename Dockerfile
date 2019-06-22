FROM golang:1.10-stretch

MAINTAINER Gabriel Garrido

RUN apt update
RUN apt install rsync -y

COPY build/lightchain /usr/bin/lightchain
RUN chmod a+x /usr/bin/lightchain

COPY ./scripts/docker.sh /root/entrypoint.sh
RUN chmod a+x /root/entrypoint.sh
ENTRYPOINT ["/root/entrypoint.sh"]

EXPOSE 8545 26657 26656

FROM golang:1.10-stretch

MAINTAINER Gabriel Garrido

RUN apt update
RUN apt install rsync wget -y

RUN wget "https://s3.eu-central-1.amazonaws.com/lightstreams-public/lightchain/latest/lightchain-linux-amd64" -O "/usr/bin/lightchain"
RUN chmod a+x /usr/bin/lightchain

COPY ./scripts/docker.sh /root/entrypoint.sh
RUN chmod a+x /root/entrypoint.sh
ENTRYPOINT ["/root/entrypoint.sh"]

EXPOSE 8545 26657 26656

FROM golang:1.10-stretch

MAINTAINER Gabriel Garrido

RUN apt-get update
RUN apt-get install -y vim apt-utils git

## Install project dependencies
RUN go get -u github.com/golang/dep/cmd/dep

# Install Lightchain
RUN mkdir -p $GOPATH/src/github.com/lightstreams-network/lightchain

COPY . $GOPATH/src/github.com/lightstreams-network/lightchain

RUN cd $GOPATH/src/github.com/lightstreams-network/lightchain && \
	make get_vendor_deps && \
	make install

RUN mkdir -p /srv/lightchain && \
	lightchain init --datadir=/srv/lightchain --force

RUN rm -rf $GOPATH/src/github.com/lightstreams-network/lightchain

COPY ./scripts/docker.sh /root/init.sh
RUN chmod a+x /root/init.sh
ENTRYPOINT ["/root/init.sh"]

EXPOSE 8545 26657 26656

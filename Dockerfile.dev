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

COPY ./scripts/docker.sh /root/entrypoint.sh
RUN chmod a+x /root/entrypoint.sh
ENTRYPOINT ["/root/entrypoint.sh"]

EXPOSE 8545 26657 26656

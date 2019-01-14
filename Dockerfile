FROM golang:1.10-stretch

MAINTAINER Gabriel Garrido

ARG ssh_prv_key

## Authorize SSH Host
RUN mkdir -p /root/.ssh && \
    chmod 0700 /root/.ssh && \
    touch /root/.ssh/known_hosts && \
    ssh-keyscan github.com > /root/.ssh/known_hosts

# Add the keys and set permissions
RUN echo "-----BEGIN RSA PRIVATE KEY-----" > /root/.ssh/id_rsa && \
	echo "$ssh_prv_key" | sed -E -e 's/[[:blank:]]+/\n/g' >> /root/.ssh/id_rsa && \
	echo "-----END RSA PRIVATE KEY-----"  >> /root/.ssh/id_rsa && \ 
	chmod 600 /root/.ssh/id_rsa

RUN apt-get update
RUN apt-get install -y vim apt-utils git

## Install project dependencies
RUN go get github.com/Masterminds/glide
RUN go get github.com/gogo/protobuf/proto
RUN go get -u github.com/golang/dep/cmd/dep

# Install Geth
RUN mkdir -p $GOPATH/src/github.com/ethereum/
RUN git clone -b v1.8.18 --single-branch git@github.com:ethereum/go-ethereum.git $GOPATH/src/github.com/ethereum/go-ethereum
RUN cd ${GOPATH}/src/github.com/ethereum/go-ethereum && \
	make geth && \
	cp -f build/bin/geth ${GOPATH}/bin/

# Install Lightchain
RUN mkdir -p $GOPATH/src/github.com/lightstreams-network/lightchain
COPY . $GOPATH/src/github.com/lightstreams-network/lightchain
RUN cd $GOPATH/src/github.com/lightstreams-network/lightchain && \
	make make get_vendor_deps && \
	make install

# Remove SSH keys and repos
RUN rm -rf /root/.ssh/

# Geth ws, rpc, p2p
EXPOSE 8545 26657 26656

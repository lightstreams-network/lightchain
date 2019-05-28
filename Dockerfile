FROM golang:1.11-stretch as builder

LABEL MAINTAINER "Gabriel Garrido"

ENV GOOS=linux
ENV GOARCH=amd64

RUN apt-get update && \
	apt-get install -y apt-utils git

WORKDIR /app

COPY go.mod .
COPY go.sum .
# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download

COPY . .

# Build the binary
RUN make build

# deployment stage
FROM scratch
WORKDIR /app
COPY --from=builder /app/build/lightchain /app/build/lightchain

CMD ["/app/build/lightchain run --datadir=/srv/lightchain --rpc --rpcaddr=0.0.0.0 --rpcport=8545 --rpcapi=eth,net,web3,personal,admin"]

EXPOSE 8545 26657 26656

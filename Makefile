.PHONY: all clean build install

DEP_REPO = github.com/golang/dep/cmd/dep
GO_ETHEREUM_REPO = github.com/ethereum/go-ethereum
TENDERMINT_REPO = github.com/tendermint/tendermint

DEP_BIN_PATH := $(shell command -v dep 2> /dev/null)

BUILD_TAGS? := lightchain

VERSION_TAG := 0.1.0
BUILD_DEBUG_FLAGS = -gcflags=all="-N -l"

all: get_vendor_deps install

check-tools:
ifdef DEP_BIN_PATH
	@echo "DEP is correctly installed"
else
	@echo "Not DEP found. Visit http://${DEP_REPO} and follow installation instructions."
endif

ifdef GETH_BIN_PATH
	@echo "GETH is correctly installed"
else
	@echo "Not GETH found. Visit http://${GO_ETHEREUM_REPO} and follow installation instructions."
endif

ifndef TENDERMINT_BIN_PATH
	@echo "TENDERMINT_REPO is correctly installed"
else
	@echo "Not DEP found. Visit http://${TENDERMINT_REPO} and follow installation instructions."
endif

install:
	CGO_ENABLED=1 go install ./cmd/lightchain

build:
	CGO_ENABLED=1 go build -o ./build/lightchain ./cmd/lightchain

### Development ###

build-dev:
	CGO_ENABLED=1 go build $(BUILD_DEBUG_FLAGS) -o ./build/lightchain ./cmd/lightchain

### Tooling ###

get_vendor_deps:
	@rm -rf vendor
	@echo "--> dep ensure"
	@dep ensure
	@rm -rf vendor/github.com/ethereum/go-ethereum/vendor
	@rm -rf vendor/github.com/tendermint/tendermint/vendor

clean:
	@rm build/lightchain

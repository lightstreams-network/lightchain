.PHONY: all clean build install docker

DEP_REPO = github.com/golang/dep/cmd/dep
DEP_BIN_PATH := $(shell command -v dep 2> /dev/null)

BUILD_TAGS? := lightchain
BUILD_DEBUG_FLAGS = -gcflags=all="-N -l"

define VERSION_TAG
	$(shell git ls-remote git@github.com:lightstreams-network/lightchain.git HEAD | cut -f1 | cut -c1-9)
endef



all: get_vendor_deps install

check-tools:
ifdef DEP_BIN_PATH
	@echo "DEP is correctly installed"
else
	@echo "Not DEP found. Visit http://${DEP_REPO} and follow installation instructions."
endif

install:
	CGO_ENABLED=1 go install ./cmd/lightchain

build:
	CGO_ENABLED=1 go build -o ./build/lightchain ./cmd/lightchain

clean:
	@rm build/lightchain

### Development ###

build-dev:
	CGO_ENABLED=1 go build $(BUILD_DEBUG_FLAGS) -o ./build/lightchain ./cmd/lightchain

### Tooling ###

get_vendor_deps:
	@rm -rf vendor
	@echo "--> dep ensure"
	@dep ensure

### Docker ###
docker:
	@echo "Build docker image"
	docker build -t lightchain:latest --build-arg version="$(VERSION_TAG)" .


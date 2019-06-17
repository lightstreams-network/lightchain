DEP_REPO = github.com/golang/dep/cmd/dep
DEP_BIN_PATH := $(shell command -v dep 2> /dev/null)

BUILD_TAGS? := lightchain
BUILD_DEBUG_FLAGS = -gcflags=all="-N -l"

define VERSION_TAG
	$(shell git ls-remote git@github.com:lightstreams-network/lightchain.git HEAD | cut -f1 | cut -c1-9)
endef

.PHONY: help
help: ## Prints this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: all
all: get_vendor_deps install ## Downloads dependencies and install lightchain

.PHONY: check-tools
check-tools: ## Check that required tools are installed
ifdef DEP_BIN_PATH
	@echo "DEP is correctly installed"
else
	@echo "Not DEP found. Visit http://${DEP_REPO} and follow installation instructions."
endif

.PHONY: install
install: ## Install lightchain
	go install ./cmd/lightchain

.PHONY: build
build: ## Build lightchain
	go build -o ./build/lightchain ./cmd/lightchain

.PHONY: clean
clean: ## Clean binaries
	@rm build/lightchain

.PHONY: build-dev
build-dev: ## (Dev) Build lightchain
	CGO_ENABLED=1 go build $(BUILD_DEBUG_FLAGS) -o ./build/lightchain ./cmd/lightchain

.PHONY: gen-bindings
gen-bindings: ## (Dev) Generate bindings
	abigen --sol ./distribution/distribution.sol --pkg distribution --out ./distribution/distribution_bindings.go

.PHONY: get_vendor_deps
get_vendor_deps: check-tools ## Download dependencies
	@rm -rf vendor
	@echo "--> dep ensure"
	@dep ensure

.PHONY: docker
docker: ## Build docker image for lightchain
	@echo "Build docker image"
	go build -o ./build/lightchain ./cmd/lightchain
	docker build --no-cache -t lightchain:latest -f ./Dockerfile .

.PHONY: docker-dev
docker-dev: ## Build docker image for lightchain
	@echo "Build docker image"
	docker build -t lightchain:latest-dev -f ./Dockerfile.dev --build-arg version="$(VERSION_TAG)"  .

.PHONY: gen-bindings-governance
gen-bindings-governance:
	solc github.com/lightstreams-network/lightchain/=${GOPATH}/src/github.com/lightstreams-network/lightchain/ --abi --bin governance/contracts/ValidatorSet.sol -o governance/build/ --overwrite
	abigen --bin=${GOPATH}/src/github.com/lightstreams-network/lightchain/governance/build/ValidatorSet.bin --abi=${GOPATH}/src/github.com/lightstreams-network/lightchain/governance/build/ValidatorSet.abi --pkg=bindings --out=${GOPATH}/src/github.com/lightstreams-network/lightchain/governance/bindings/ValidatorSet.go --type=ValidatorSet

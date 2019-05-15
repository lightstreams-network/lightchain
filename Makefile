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

.PHONY: docker
docker: ## Build docker image for lightchain
	@echo "Build docker image"
	docker build -t lightchain:latest --build-arg version="$(VERSION_TAG)" .

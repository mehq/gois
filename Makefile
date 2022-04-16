GO111MODULE ?= on
export GO111MODULE
GOPROXY ?= https://proxy.golang.org
export GOPROXY

BINARY ?= gois

all: check test build

build: ## Build for current platform
	@go build -o bin/$(BINARY)

check: deps ## Check source code issues
	@golangci-lint run ./...

clean: ## Remove build files
	@go clean
	@rm -fr ./.build ./.release ./.tarballs ./.tmp ./bin ./dist ./gois ./vendor ./coverage.txt ||:

deps: ## Ensures fresh go.mod and go.sum.
	@go mod tidy
	@go mod verify

fmt: ## Format code
	@gofmt -w .
	@goimports -local github.com/mehq/gois -w .

help: ## Show this help
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    %-20s%s\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  %s\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)

test: ## Run tests
	@go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

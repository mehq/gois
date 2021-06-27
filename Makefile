GOCMD ?= go
GOPATH ?= $(shell go env GOPATH)
GOBIN ?= $(firstword $(subst :, ,${GOPATH}))/bin

GO111MODULE ?= on
export GO111MODULE
GOPROXY ?= https://proxy.golang.org
export GOPROXY

DIST_DIR = dist
BINARY_NAME = gomage
BUILD_NUMBER ?= $(shell git rev-list HEAD --count)
VERSION ?= X.X.X
LDFLAGS = -ldflags "-X main.buildNumber=${BUILD_NUMBER} -X main.programVersion=${VERSION} -X main.programName=${BINARY_NAME}"

all: test build

build: ## Build for current platform
	@$(GOCMD) build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)

check: deps ## Check source code issues
	@bash scripts/gocheck.sh

clean: ## Remove build/ci related file
	@$(GOCMD) clean
	@rm -fr ./dist ||:
	@rm coverage.txt ||:

crossbuild: ## Build for multiple platforms
	@GOOS=darwin GOARCH=amd64 $(GOCMD) build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-$(VERSION).darwin-amd64/$(BINARY_NAME)
	@GOOS=freebsd GOARCH=amd64 $(GOCMD) build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-$(VERSION).freebsd-amd64/$(BINARY_NAME)
	@GOOS=linux GOARCH=amd64 $(GOCMD) build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-$(VERSION).linux-amd64/$(BINARY_NAME)
	@GOOS=linux GOARCH=arm64 $(GOCMD) build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-$(VERSION).linux-arm64/$(BINARY_NAME)
	@GOOS=windows GOARCH=amd64 $(GOCMD) build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-$(VERSION).windows-amd64/$(BINARY_NAME).exe

deps: ## Ensures fresh go.mod and go.sum.
	@$(GOCMD) mod tidy
	@$(GOCMD) mod verify

help: ## Show this help
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    %-20s%s\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  %s\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)

install: ## Install the program
	@GOBIN=$(GOBIN) $(GOCMD) install -v $(LDFLAGS)

release: clean crossbuild ## Build release artifacts
	@cp CHANGELOG.md LICENSE README.md $(DIST_DIR)/$(BINARY_NAME)-$(VERSION).darwin-amd64/
	@cp CHANGELOG.md LICENSE README.md $(DIST_DIR)/$(BINARY_NAME)-$(VERSION).freebsd-amd64/
	@cp CHANGELOG.md LICENSE README.md $(DIST_DIR)/$(BINARY_NAME)-$(VERSION).linux-amd64/
	@cp CHANGELOG.md LICENSE README.md $(DIST_DIR)/$(BINARY_NAME)-$(VERSION).linux-arm64/
	@cp CHANGELOG.md LICENSE README.md $(DIST_DIR)/$(BINARY_NAME)-$(VERSION).windows-amd64/
	@cd $(DIST_DIR)/$(BINARY_NAME)-$(VERSION).darwin-amd64/; tar -zcf ../$(BINARY_NAME)-$(VERSION).darwin-amd64.tar.gz *
	@cd $(DIST_DIR)/$(BINARY_NAME)-$(VERSION).freebsd-amd64/; tar -zcf ../$(BINARY_NAME)-$(VERSION).freebsd-amd64.tar.gz *
	@cd $(DIST_DIR)/$(BINARY_NAME)-$(VERSION).linux-amd64/; tar -zcf ../$(BINARY_NAME)-$(VERSION).linux-amd64.tar.gz *
	@cd $(DIST_DIR)/$(BINARY_NAME)-$(VERSION).linux-arm64/; tar -zcf ../$(BINARY_NAME)-$(VERSION).linux-arm64.tar.gz *
	@cd $(DIST_DIR)/$(BINARY_NAME)-$(VERSION).windows-amd64/; tar -zcf ../$(BINARY_NAME)-$(VERSION).windows-amd64.tar.gz *

test: ## Run tests
	@$(GOCMD) test -v -race -coverprofile=coverage.txt -covermode=atomic -tags test

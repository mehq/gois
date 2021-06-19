BINARY = gomage
VERSION = 0.0.1
BUILD = `git rev-parse HEAD`

PREFIX ?= /usr/local
BINDIR ?= $(PREFIX)/bin

# Setup linker flags option for build that interoperate with variable names in src code
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"

.PHONY: all test build vendor

all: check test install clean

build: ## Build project for current platform
	mkdir -p build/bin
	GO111MODULE=on go build $(LDFLAGS) -mod vendor -o build/bin/$(BINARY) .

check: ## Check source code issues
	bash scripts/gocheck.sh

clean: ## Remove build/ci related file
	rm -fr ./build
	rm coverage.txt

help: ## Show this help.
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)

install: build ## Build and install the program
	install -d $(DESTDIR)$(BINDIR)
	install -m 755 build/bin/$(BINARY) $(DESTDIR)$(BINDIR)

test: ## Run tests
	go test -v -race -coverprofile=coverage.txt -covermode=atomic -tags test

vendor:
	go mod vendor

run: ## Run the program
	go run .

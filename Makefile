GOCMD ?= go
GOBUILD = $(GOCMD) build -v
GOCLEAN = $(GOCMD) clean -v
GOTEST = $(GOCMD) test -v
PREFIX ?= /usr/local
BINDIR ?= $(PREFIX)/bin
BINARY_NAME = gomage

BUILD ?= $(shell git rev-list HEAD --count)
VERSION ?= $(shell git describe --tags --abbrev=0 2>/dev/null)
LDFLAGS = -ldflags "-X main.buildNumber=${BUILD} -X main.programVersion=${VERSION} -X main.programName=${BINARY_NAME}"

all: test build

build: ## Build project for current platform
	$(GOBUILD) $(LDFLAGS) -mod vendor -o dist/$(BINARY_NAME) .

check: ## Check source code issues
	bash scripts/gocheck.sh

clean: ## Remove build/ci related file
	$(GOCLEAN)
	rm -fr ./dist
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
	install -m 755 dist/$(BINARY_NAME) $(DESTDIR)$(BINDIR)

test: ## Run tests
	$(GOTEST) -race -coverprofile=coverage.txt -covermode=atomic -tags test

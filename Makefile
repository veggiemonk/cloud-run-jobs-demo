SHELL := /bin/bash

# CHECK EXAMPLES: https://makefiletutorial.com/

# $(shell git rev-parse --short HEAD)
# VERSION := 1.0
GO   	= $(shell which go)
BIN  	= $(CURDIR)/bin

## TOOLS
GOLANGCI_LINT = $(or $(shell which golangci-lint), $(error "Missing dependency - no golangci-lint in PATH."))
PRESENT = $(or $(shell which present), $(error "Missing dependency - no present in PATH. Run `go get -u golang.org/x/tools/cmd/present`"))

## PROJECT VARIABLES
GITSHA        =$(shell git rev-parse --short HEAD)
DATE		  =$(shell TZ=UTC date +"%Y-%m-%dT%TZ")
VERSION       =$(shell echo $(DATE)-$(GITSHA))

## BUILD FLAGS
GOBUILD = -a -v -trimpath='true' -buildmode='exe' -buildvcs='true' -compiler='gc' -mod='vendor'
LDFLAGS = -X main.gitSha=$(GITSHA) -X main.buildTime=$(DATE)

# ==============================================================================
# MAIN TARGETS

.PHONY: bin
bin: ## Build the binary file
	@echo "Building binaries..."
	rm -rf "$(BIN)"
	mkdir -p "$(BIN)"
	@echo "VERSION: $(VERSION)"
	$(GO) build $(GOBUILD) -ldflags "$(LDFLAGS)" -o "$(BIN)/batch" *.go

.PHONY: dev
dev: ## Build the binary file
	@echo "Building binaries..."
	rm -rf "$(BIN)"
	mkdir -p "$(BIN)"
	@echo "VERSION: $(VERSION)"
	$(GO) build -mod='vendor' -o "$(BIN)/batch"  *.go

.PHONY: lint
lint: ## Lint the files
	golangci-lint run ./...
	@echo "GOLANGCI-LINT DONE"
	@echo
	staticcheck -checks=all ./...

# ==============================================================================
# Running tests within the local computer

.PHONY: test
test: ## Run all tests
	go test -v ./... -count=1

# ==============================================================================
# Presentation

.PHONY: presentation
presentation: ## Run golang.org/x/tools/present
	present -base slides/ -content slides/



# ==============================================================================
# Modules support

.PHONY: deps-reset
deps-reset: ## Reset all dependencies
	git checkout -- go.mod
	$(GO) mod tidy
	$(GO) mod vendor

.PHONY: tidy
tidy: ## Update go.mod and go.sum
	$(GO) mod tidy
	$(GO) mod vendor

.PHONY: deps-upgrade
deps-upgrade: ## Upgrade all dependencies
	$(GO) get -u -v ./...
	$(GO) mod tidy
	$(GO) mod vendor

.PHONY: deps-cleancache
deps-cleancache: ## Clean the module download cache
	$(GO) clean -modcache

# To remove deps, example: go get github.com/pkg/errors@none


# ==============================================================================
# meta targets

.PHONY: list
list:
	@LC_ALL=C $(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/(^|\n)# Files(\n|$$)/,/(^|\n)# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$'


.PHONY: help
help: ## Show help
	@echo Please specify a build target. The choices are:
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "    \033[36m%-30s\033[0m %s\n", $$1, $$2}'

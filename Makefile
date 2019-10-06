# Project information
VERSION? := $(shell git describe --tags)
BUILD := $(shell git rev-parse --short HEAD)

# Go build variables
GOBASE := $(shell pwd)
GOPATH := $(GOBASE)/vendor:$(GOBASE)
GOBIN := $(GOBASE)/bin
GOFILES := $(wildcard *.go)

CMD := $(GOBASE)/cmd

# Linker flags
LDFLAGS=-v -ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

.PHONY: all build clean lint

all: build

dep: clean ## Confirm vendor directory with dependencies
	@printf "\033[36m%-30s\033[0m\n" "Confirm dependencies"
	@-go mod vendor

lint: ## Lint the files
	@printf "\033[36m%-30s\033[0m\n" "Lint source code"
	@golint pkg/...
	@golint internal/...
	@golint cmd/...

build: dep test ## Build the binary files
	@printf "\033[36m%-30s\033[0m\n" "Build binaries"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go build $(LDFLAGS) -o $(GOBIN)/checker $(CMD)/checker
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go build $(LDFLAGS) -o $(GOBIN)/cli $(CMD)/cli

test: lint ## Test the library
	@printf "\033[36m%-30s\033[0m\n" "Perform covered tests"
	@go test -race -timeout 10000ms ./... -coverprofile artifacts/cover.out
	@go tool cover -html=artifacts/cover.out -o artifacts/cover.html
	@go tool cover -func=artifacts/cover.out

clean: ## Remove previous build
	@printf "\033[36m%-30s\033[0m\n" "Clean"
	@-chmod -R 777 ./vendor
	@-rm -rf ./bin
	@-rm -rf ./vendor
	@GO111MODULE=off go clean
	@GO111MODULE=off go clean -modcache

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

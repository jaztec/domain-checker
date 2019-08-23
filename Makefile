# Project information
VERSION? := $(shell git describe --tags)
BUILD := $(shell git rev-parse --short HEAD)
PROJECTNAME := $(shell basename "$(PWD)")

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

build: dep test ## Build the binary file
	@printf "\033[36m%-30s\033[0m\n" "Build binaries"
	@mkdir -p ./bin/limbo
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go build $(LDFLAGS) -o $(GOBIN)/limbo/limbo $(CMD)/$(PROJECTNAME)

race: lint ## Run race tests on the library
	@printf "\033[36m%-30s\033[0m\n" "Perform race tests"
	@go test ./... -race -timeout 10000ms

test: lint race ## Test the library
	@printf "\033[36m%-30s\033[0m\n" "Perform covered tests"
	@go test ./... -coverprofile artifacts/cover.out
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

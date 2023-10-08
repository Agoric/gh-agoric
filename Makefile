BIN=
BINARY_NAME=gh-agoric
EXTENSION_NAME=agoric
SRC=$(shell find . -name "*.go")

ifeq (, $(shell which golangci-lint))
$(warning "could not find golangci-lint in $(PATH), run: curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh")
endif

.PHONY: build run fmt lint test install_deps clean

default: all

all: fmt test build

build:
	go build -o $(BINARY_NAME)

install: build
	gh extension remove $(EXTENSION_NAME)
	gh extension install .

run: install
	gh agoric

fmt:
	gofmt -w $(SRC)

lint:
	golangci-lint run

test: install_deps
	go test -v ./...

install_deps:
	go get -v ./...

clean:
	rm -rf $(BINARY_NAME)
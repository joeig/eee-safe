GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOCOVER=$(GOCMD) tool cover
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOFMT=gofmt
BINARY_NAME=eee-safe
GOFILES=$(shell find . -type f -name '*.go' -not -path "./vendor/*")

.DEFAULT_GOAL := all
.PHONY: all build build-linux-amd64 coverage test check-fmt fmt clean run run-debug

all: check-fmt test coverage build

build:
	mkdir -p ./out
	$(GOBUILD) -o ./out/$(BINARY_NAME) -v ./cmd/eee-safe

build-linux-amd64:
	mkdir -p ./out
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o ./out/$(BINARY_NAME)_linux_amd64 -v ./cmd/eee-safe

coverage:
	$(GOCOVER) -func=./out/coverage.out

test:
	mkdir -p ./out
	mkdir -p cmd/eee-safe/threema-backups
	$(GOTEST) -v ./... -covermode=count -coverprofile=./out/coverage.out

check-fmt:
	$(GOFMT) -d ${GOFILES}

fmt:
	$(GOFMT) -w ${GOFILES}

clean:
	$(GOCLEAN)
	rm -rf ./out

run: build
	./out/$(BINARY_NAME) -config=configs/config.yml

run-debug: build
	./out/$(BINARY_NAME) -config=configs/config.yml -debug

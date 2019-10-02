GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOCOVER=$(GOCMD) tool cover
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOFMT=gofmt
GODEP=dep
BINARY_NAME=eee-safe
GOFILES=$(shell find . -type f -name '*.go' -not -path "./vendor/*")

.DEFAULT_GOAL := all
.PHONY: all build build-linux-amd64 coverage test check-fmt fmt clean run deps

all: check-fmt test coverage build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/eee-safe

build-linux-amd64:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)_linux_amd64 -v ./cmd/eee-safe

coverage:
	$(GOCOVER) -func=c.out

test:
	mkdir -p cmd/eee-safe/threema-backups
	$(GOTEST) -v ./... -covermode=count -coverprofile=c.out

check-fmt:
	$(GOFMT) -d ${GOFILES}

fmt:
	$(GOFMT) -w ${GOFILES}

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME)_linux_amd64

run:
	$(GOBUILD) -o $(BINARY_NAME) -v
	./$(BINARY_NAME)

run-debug:
	$(GOBUILD) -o $(BINARY_NAME) -v
	./$(BINARY_NAME) -debug

deps:
	$(GODEP) ensure

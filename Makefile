.PHONY: build install clean test lint fmt generate

VERSION ?= dev
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

build:
	go build $(LDFLAGS) -o notte ./cmd/notte

install:
	go install $(LDFLAGS) ./cmd/notte

clean:
	rm -f notte
	go clean

test:
	go test -v ./...

lint:
	golangci-lint run

fmt:
	goimports -w .
	gofumpt -w .

generate:
	./scripts/generate.sh

.DEFAULT_GOAL := build

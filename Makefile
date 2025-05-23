VERSION := $(shell git describe --tags)

PHONY: build
build: test
	go build -o bin/synth -ldflags="-X github.com/iljarotar/synth/cmd.version=$(VERSION)"

PHONY: test
test:
	go test ./...

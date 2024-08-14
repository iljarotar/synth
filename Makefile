VERSION := $(shell git describe --tags --exact-match 2> /dev/null || git symbolic-ref -q --short HEAD || git rev-parse --short HEAD)

PHONY: build
build:
	go build -o bin/synth -ldflags="-X github.com/iljarotar/synth/cmd.version=$(VERSION)"

PHONY: test
test:
	go test ./...

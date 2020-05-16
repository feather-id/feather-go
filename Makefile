SHELL := /bin/bash
PKGS := $(shell go list ./... | grep -v /vendor | grep -v /config)
.DEFAULT_GOAL := all

# Make the repo
all: test

# Run tests
test:
	go test -cover $(PKGS)

coverage:
	go test -cover $(PKGS) -covermode=count -coverprofile=combined.coverprofile ./...

clean:
	find . -name \*.coverprofile -delete

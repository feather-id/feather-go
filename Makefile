SHELL := /bin/bash
PKGS := $(shell go list ./... | grep -v /vendor | grep -v /config)
.DEFAULT_GOAL := all

# Make the repo
all: clean mock test

# Cleanup builds and generated code
clean:
	go clean
	rm -rf .gen
	rm -rf tmp

# Generate mocks for testing
mock:
	for SRC in $(shell find .$(MDIR) -not -name "*_test.go" -not -name main.go -path "./*.go") ; do \
		mkdir -p .gen/mock/$$SRC ; \
		rm -rf .gen/mock/$$SRC ; \
		mockgen -source=$$SRC -destination=.gen/mock/$$SRC ; \
	done

# Run tests
test:
	go test -cover $(PKGS)

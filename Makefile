GO?=go
GOFLAGS?=-mod=vendor
BINARY?=gorunner

.PHONY: install
install:
	$(GO) mod tidy
	$(GO) mod vendor

.PHONY: build
## build: build executable file
build:
	$(GO) build -o ./build/$(BINARY) ./*.go

.PHONY: build-test
## build-test: build tests
build-test:
	$(MAKE) -C test build

.PHONY: test
## test: run tests
test:
	$(MAKE) -C test test

.PHONY: clean
## clean: clean build and dependencies files
clean:
	if [ -e ./build ] ; then rm -r ./build; fi
	if [ -e ./vendor ] ; then rm -r ./vendor; fi

.PHONY: help
## help: prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.DEFAULT_GOAL := help

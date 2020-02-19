APP=gorunner
BUILD_DIR := "."
DOCKER_BUILDKIT=1
DOCKER_IMAGE?=dundergitcall/go-runner

ifndef COMMIT
    export COMMIT=$(shell git rev-list -1 HEAD)
endif
ifndef BUILD_DATE
    export BUILD_DATE=$(shell date +"%Y-%m-%d_%T%z")
endif
ifndef VERSION
    export VERSION=$(BUILD_DATE)
endif

GO?=go
GOARGS?=
CGO_ENABLED?=0
BINARY?=${APP}
LDFLAGS:="${LDFLAGS} -X main.buildCommit=$(COMMIT) -X main.buildVersion=$(VERSION) -X main.buildDate=$(BUILD_DATE)"

.PHONY: install
install:
	$(GO) mod tidy
	$(GO) mod vendor

.PHONY: build
## build: build executable file
build:
	$(GO) build $(GOARGS) -mod=vendor -ldflags $(LDFLAGS) $(GORACE) -o ./build/$(BINARY) ./*.go

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

.PHONY: check-docker-envs
check-docker-envs:
ifndef DOCKER_REGISTRY
	$(error DOCKER_REGISTRY is undefined)
endif
ifndef DOCKER_IMAGE
	$(error DOCKER_IMAGE is undefined)
endif

.PHONY: help
## help: prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.DEFAULT_GOAL := help

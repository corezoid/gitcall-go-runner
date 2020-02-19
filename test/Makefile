.PHONY: build
## build: build executable file
build:
	CGO_ENABLED=1 go build -race -mod=vendor -o ../build/gorunner-success ./apps/success/*.go
	CGO_ENABLED=1 go build -race -mod=vendor -o ../build/gorunner-error ./apps/error/*.go

.PHONY: test
## test: run tests
test:
	@go test -mod=vendor -race -count=1 -cover ./tests/... -v

.PHONY: help
## help: prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.DEFAULT_GOAL := help

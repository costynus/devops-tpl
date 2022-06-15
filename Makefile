# HELP =================================================================================================================
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
.PHONY: help

run-agent: ### run agent
	go mod tidy & go mod download &&\
	go run ./cmd/agent/main.go
.PHONY: run-agent

run-server: ### run server
	go mod tidy & go mod download &&\
	go run ./cmd/server/main.go
.PHONY: run-server

unit-test: ### run unit-test
	go test -v -cover -race ./internal/...
.PHONY: test

test: ### run test
	make unit-test
	# TODO: make integration-test
.PHONY: test

test-agent:
	go test -v -cover -race ./internal/app/agent/...
.PHONY: test-agent

test-server:
	go test -v -cover -race ./internal/app/server/...
.PHONY: test-server
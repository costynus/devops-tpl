ENV_FILE = .env.local
include ${ENV_FILE}
export

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

compose-up-server: ### Run server with docker
	docker-compose up devops_server
.PHONY: compose-up-server

compose-down: ### Down docker
	docker-compose down --remove-orphans
.PHONY: compose-down-server

compose-up-integration-test: ### Run docker-compose with integration test
	docker-compose -f docker-compose.test.yaml up --build --abort-on-container-exit --exit-code-from integration
.PHONY: compose-up-integration-test

build: ### BUILD local with docker
	docker-compose -f docker-compose.yaml build
.PHONY: build

run: ### RUN local with docker
	docker-compose -f docker-compose.yaml up
.PHONY: run

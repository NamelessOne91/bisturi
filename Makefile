export GO_EXECUTABLE_PATH := $(shell which go)

build:
	@cd cmd/bisturi/ && go build -o ../../bin/bisturi
	@echo "Grant the capability to create raw sockets to the binary executable ..."
	@sudo setcap cap_net_raw=eip ./bin/bisturi

run: build
	@./bin/bisturi

run-debug: build
	@export BISTURI_DEBUG=y && ./bin/bisturi

test:
	@$$GO_EXECUTABLE_PATH test -v -race ./...

coverage:
	@$$GO_EXECUTABLE_PATH test -v -race --cover --coverprofile=cover.profile ./...

coverage-report: coverage
	@go tool cover -html=cover.profile
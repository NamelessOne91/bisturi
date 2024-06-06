export GO_EXECUTABLE_PATH := $(shell which go)

build:
	@cd cmd/bisturi/ && go build -o ../../bin/bisturi

run: build
	@./bin/bisturi

test:
	@$$GO_EXECUTABLE_PATH test -v ./...
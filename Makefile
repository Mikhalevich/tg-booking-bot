MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
ROOT := $(dir $(MKFILE_PATH))
GOBIN ?= $(ROOT)/tools/bin
ENV_PATH = PATH=$(GOBIN):$(PATH)
BIN_PATH ?= $(ROOT)/bin
LINTER_NAME := golangci-lint
LINTER_VERSION := v1.60.2

.PHONY: all build test compose-up vendor install-linter lint tools generate

all: build

build:
	go build -mod=vendor -o $(BIN_PATH)/bot ./cmd/schedule/main.go

test:
	go test ./...

compose-up:
	docker-compose -f ./script/docker/docker-compose.yml up --build

vendor:
	go mod tidy
	go mod vendor

install-linter:
	if [ ! -f $(GOBIN)/$(LINTER_VERSION)/$(LINTER_NAME) ]; then \
		echo INSTALLING $(GOBIN)/$(LINTER_VERSION)/$(LINTER_NAME) $(LINTER_VERSION) ; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN)/$(LINTER_VERSION) $(LINTER_VERSION) ; \
		echo DONE ; \
	fi

lint: install-linter
	$(GOBIN)/$(LINTER_VERSION)/$(LINTER_NAME) run --config .golangci.yml

tools: install-linter
	@if [ ! -f $(GOBIN)/mockgen ]; then\
		echo "Installing mockgen";\
		GOBIN=$(GOBIN) go install go.uber.org/mock/mockgen@v0.5.0;\
	fi

generate: tools
	$(ENV_PATH) go generate ./...

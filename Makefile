.PHONY: deps build test dev all clean

PLUGIN_NAME := vault-plugin-secrets-cloudflare
BIN_DIR := bin

DEPS_TARGETS := go.mod go.sum

$(BIN_DIR):
	@mkdir -p $(BIN_DIR)

deps: $(DEPS_TARGETS)
	go mod download

build: deps | $(BIN_DIR)
	CGO_ENABLED=0 go build -o $(BIN_DIR)/$(PLUGIN_NAME) cmd/cloudflare/main.go

test: deps
	CGO_ENABLED=0 go test ./...

dev: build
	bash ./scripts/dev.sh

all: build test

clean:
	rm -rf $(BIN_DIR)

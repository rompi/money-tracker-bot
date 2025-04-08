# Project settings
APP_NAME=telegram-bot-go
MAIN_PKG=./cmd/telebot
BINARY_NAME=bot

# Environment variables (optional override via CLI)
TELEGRAM_BOT_TOKEN ?=

# Default target
.PHONY: all
all: run

## Build the binary
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) $(MAIN_PKG)

## Run the bot
.PHONY: run
run:
	@echo "Running bot..."
	go run $(MAIN_PKG)

## Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)

## Run tests
.PHONY: test
test:
	@echo "Running tests..."
	go test ./...

## Format code
.PHONY: fmt
fmt:
	go fmt ./...

## Lint code
.PHONY: lint
lint:
	go vet ./...

## Show help
.PHONY: help
help:
	@echo "Usage:"
	@echo "  make build         Build the binary"
	@echo "  make run           Run the bot (requires TELEGRAM_BOT_TOKEN)"
	@echo "  make clean         Remove build artifacts"
	@echo "  make test          Run tests"
	@echo "  make fmt           Format code"
	@echo "  make lint          Lint code"

.PHONY: all build test clean generate update-grammar install release

# Version info
VERSION := $(shell cat VERSION)
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT) -X main.BuildTime=$(BUILD_TIME)"

# Variables
ANTLR_VERSION := 4.13.1
ANTLR_JAR := antlr-$(ANTLR_VERSION)-complete.jar
ANTLR_URL := https://www.antlr.org/download/$(ANTLR_JAR)
GRAMMAR_DIR := grammar
GEN_DIR := internal/gen
BINARY := solast

# Default target
all: build test

# Download ANTLR JAR if not present
$(ANTLR_JAR):
	@echo "Downloading ANTLR $(ANTLR_VERSION)..."
	curl -O $(GRAMMAR_DIR)/$(ANTLR_JAR) $(ANTLR_URL)

# Generate Go parser from ANTLR grammar
generate: $(ANTLR_JAR)
	@echo "Generating parser from grammar..."
	@mkdir -p $(GEN_DIR)
	java -jar $(ANTLR_JAR) -Dlanguage=Go -visitor -package gen -o $(GEN_DIR) $(GRAMMAR_DIR)/SolidityLexer.g4 $(GRAMMAR_DIR)/SolidityParser.g4
	@echo "Parser generated successfully"

# Update grammar from official Solidity repository
update-grammar:
	@echo "Updating grammar from ethereum/solidity..."
	@mkdir -p $(GRAMMAR_DIR)
	curl -o $(GRAMMAR_DIR)/SolidityLexer.g4 https://raw.githubusercontent.com/ethereum/solidity/develop/docs/grammar/SolidityLexer.g4
	curl -o $(GRAMMAR_DIR)/SolidityParser.g4 https://raw.githubusercontent.com/ethereum/solidity/develop/docs/grammar/SolidityParser.g4
	@echo "Grammar updated from official Solidity repo."
	@echo "Note: Manual adjustments may be needed for Go target compatibility."

# Build the CLI binary with version info
build:
	@echo "Building $(BINARY) v$(VERSION)..."
	go build $(LDFLAGS) -o $(BINARY) ./cmd/solast

# Build without version info (dev)
build-dev:
	@echo "Building $(BINARY) (dev)..."
	go build -o $(BINARY) ./cmd/solast

# Install the CLI with version info
install:
	@echo "Installing $(BINARY) v$(VERSION)..."
	go install $(LDFLAGS) ./cmd/solast

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Run tests verbose
test-v:
	@echo "Running tests (verbose)..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Build release binaries for multiple platforms
release: clean
	@echo "Building release binaries v$(VERSION)..."
	@mkdir -p dist
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY)-linux-amd64 ./cmd/solast
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY)-linux-arm64 ./cmd/solast
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY)-darwin-amd64 ./cmd/solast
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY)-darwin-arm64 ./cmd/solast
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY)-windows-amd64.exe ./cmd/solast
	@echo "Release binaries built in dist/"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -f $(BINARY)
	rm -f coverage.out coverage.html
	rm -rf $(GEN_DIR)
	rm -rf dist

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	go vet ./...

# Show version (always run, not a file target)
.PHONY: version
version:
	@echo "$(VERSION)"

# Help
help:
	@echo "solast-go v$(VERSION)"
	@echo ""
	@echo "Available targets:"
	@echo "  all            - Build and test"
	@echo "  build          - Build the CLI binary with version info"
	@echo "  build-dev      - Build the CLI binary (dev, no version)"
	@echo "  install        - Install the CLI"
	@echo "  test           - Run tests"
	@echo "  test-v         - Run tests (verbose)"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  release        - Build release binaries for all platforms"
	@echo "  clean          - Remove build artifacts"
	@echo "  fmt            - Format Go code"
	@echo "  lint           - Run linter"
	@echo "  version        - Show version"
	@echo "  generate       - Generate Go parser from ANTLR grammar"
	@echo "  update-grammar - Download latest grammar from upstream"

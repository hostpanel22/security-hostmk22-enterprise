# HostMK Security Framework - Makefile
# Usage: make <target>

.PHONY: help build run dev test lint clean docker docker-build docker-run docker-down

# ============================================================================
# Default Target
# ============================================================================
help: ## Show this help message
	@echo "HostMK Security Framework - Makefile"
	@echo "===================================="
	@echo ""
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ""

# ============================================================================
# Environment Variables
# ============================================================================
GOPATH ?= $(shell go env GOPATH)
BINARY_NAME ?= hostmk-security
BINARY_PATH ?= ./bin/$(BINARY_NAME)
DOCKER_IMAGE ?= hostmk/security-framework
DOCKER_TAG ?= latest
VERSION ?= 1.0.0
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE ?= $(shell date +%Y-%m-%dT%H:%M:%SZ)

# ============================================================================
# Build
# ============================================================================
build: go-build ## Build the application

.go-build: go-mod-tidy go-vet go-fmt
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
		-ldflags="-s -w -X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.Date=$(DATE) -X main.BuiltBy=make" \
		-o $(BINARY_PATH) \
		./cmd/server
	@echo "Build complete: $(BINARY_PATH)"

build-all: go-build build-linux build-darwin build-windows ## Build for all platforms

build-linux: go-mod-tidy
	@echo "Building for Linux..."
	@mkdir -p bin/linux
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
		-ldflags="-s -w -X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.Date=$(DATE) -X main.BuiltBy=make" \
		-o bin/linux/$(BINARY_NAME)-linux-amd64 \
		./cmd/server
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build \
		-ldflags="-s -w -X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.Date=$(DATE) -X main.BuiltBy=make" \
		-o bin/linux/$(BINARY_NAME)-linux-arm64 \
		./cmd/server
	@echo "Linux builds complete"

build-darwin: go-mod-tidy
	@echo "Building for Darwin..."
	@mkdir -p bin/darwin
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build \
		-ldflags="-s -w -X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.Date=$(DATE) -X main.BuiltBy=make" \
		-o bin/darwin/$(BINARY_NAME)-darwin-amd64 \
		./cmd/server
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build \
		-ldflags="-s -w -X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.Date=$(DATE) -X main.BuiltBy=make" \
		-o bin/darwin/$(BINARY_NAME)-darwin-arm64 \
		./cmd/server
	@echo "Darwin builds complete"

build-windows: go-mod-tidy
	@echo "Building for Windows..."
	@mkdir -p bin/windows
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build \
		-ldflags="-s -w -X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.Date=$(DATE) -X main.BuiltBy=make" \
		-o bin/windows/$(BINARY_NAME)-windows-amd64.exe \
		./cmd/server
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build \
		-ldflags="-s -w -X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.Date=$(DATE) -X main.BuiltBy=make" \
		-o bin/windows/$(BINARY_NAME)-windows-386.exe \
		./cmd/server
	@echo "Windows builds complete"

# ============================================================================
# Run
# ============================================================================
run: go-build ## Run the application
	@echo "Starting $(BINARY_NAME)..."
	@$(BINARY_PATH) serve

run-dev: go-mod-tidy ## Run in development mode
	@echo "Starting $(BINARY_NAME) in development mode..."
	@go run ./cmd/server serve

run-migrate: go-mod-tidy ## Run database migrations
	@echo "Running database migrations..."
	@go run ./cmd/server migrate up

run-seed: go-mod-tidy ## Seed the database
	@echo "Seeding database..."
	@go run ./cmd/server seed

# ============================================================================
# Development
# ============================================================================
dev: run-dev ## Start development server

start: run-dev ## Start development server (alias)

# ============================================================================
# Test
# ============================================================================
test: go-test ## Run all tests

go-test: go-mod-tidy go-fmt go-vet
	@echo "Running tests..."
	@go test -race -v -coverprofile=coverage.out -covermode=atomic ./...

test-short: go-mod-tidy ## Run tests without race detector
	@echo "Running tests (short)..."
	@go test -v -short ./...

test-cover: go-mod-tidy ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -race -v -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -html=coverage.out -o coverage.html

test-unit: go-mod-tidy ## Run unit tests only
	@echo "Running unit tests..."
	@go test -v -run "^Test[^I]" ./...

test-integration: go-mod-tidy ## Run integration tests only
	@echo "Running integration tests..."
	@go test -v -run "^TestI" ./...

# ============================================================================
# Lint & Format
# ============================================================================
lint: go-lint go-fmt go-vet go-cyclo ## Run all linters

go-lint: go-mod-tidy
	@echo "Running golangci-lint..."
	@golangci-lint run --config=.golangci.yml --timeout=5m

go-fmt: go-mod-tidy ## Format code
	@echo "Formatting code..."
	@gofmt -w -s .
	@echo "Code formatted"

go-vet: go-mod-tidy ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

go-cyclo: go-mod-tidy ## Check cyclomatic complexity
	@echo "Checking cyclomatic complexity..."
	@gocyclo -over 15 .

# ============================================================================
# Security
# ============================================================================
security: go-security govulncheck gosec trivy ## Run all security checks

go-security: go-mod-tidy ## Run Go security checks
	@echo "Running Go security checks..."
	@golangci-lint run --config=.golangci.yml --timeout=5m --enable gosec

govulncheck: go-mod-tidy ## Run Go vulnerability checker
	@echo "Running govulncheck..."
	@go install golang.org/x/vuln/cmd/govulncheck@latest 2>/dev/null || true
	@govulncheck ./...

gosec: go-mod-tidy ## Run Gosec security scanner
	@echo "Running gosec..."
	@gosec -include=G101,G102,G103,G104,G106,G107,G108,G201,G202,G203,G204,G301,G302,G304,G306,G307,G401,G402,G403,G404,G501,G502,G503,G504,G505,G601 ./...

trivy: go-mod-tidy ## Run Trivy vulnerability scanner
	@echo "Running trivy..."
	@trivy fs --security-checks vuln,config --exit-code 1 --ignore-unfixed .

# ============================================================================
# Clean
# ============================================================================
clean: go-clean docker-clean ## Clean build artifacts

go-clean: ## Clean Go build artifacts
	@echo "Cleaning Go artifacts..."
	@rm -rf bin/
	@rm -rf coverage.out coverage.html
	@rm -rf dist/
	@go clean -cache -testcache

docker-clean: ## Clean Docker artifacts
	@echo "Cleaning Docker artifacts..."
	@docker system prune -f --volumes
	@rm -rf .docker

# ============================================================================
# Docker
# ============================================================================
docker: docker-build docker-run ## Build and run Docker container

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg DATE=$(DATE) \
		-t $(DOCKER_IMAGE):$(DOCKER_TAG) \
		-t $(DOCKER_IMAGE):$(COMMIT) \
		.
	@echo "Docker image built: $(DOCKER_IMAGE):$(DOCKER_TAG)"

docker-run: docker-build ## Run Docker container
	@echo "Running Docker container..."
	@docker run -d \
		--name hostmk-security \
		--restart unless-stopped \
		-p 8080:8080 \
		-p 9090:9090 \
		-e APP_ENV=production \
		-e DATABASE_URL=postgresql://postgres:postgres@host.docker.internal:5432/security \
		-e REDIS_URL=redis://host.docker.internal:6379/0 \
		-e JWT_SECRET=$(JWT_SECRET) \
		-e ENCRYPTION_KEY=$(ENCRYPTION_KEY) \
		--security-opt no-new-privileges \
		--cap-drop=ALL \
		--cap-add=NET_BIND_SERVICE \
		--read-only \
		--tmpfs /tmp:size=100M,mode=1777 \
		$(DOCKER_IMAGE):$(DOCKER_TAG)

docker-down: ## Stop and remove Docker container
	@echo "Stopping Docker container..."
	@docker stop hostmk-security 2>/dev/null || true
	@docker rm hostmk-security 2>/dev/null || true

docker-logs: ## Show Docker container logs
	@docker logs -f hostmk-security 2>/dev/null || echo "Container not running"

docker-shell: ## Open shell in Docker container
	@docker exec -it hostmk-security sh 2>/dev/null || echo "Container not running"

# ============================================================================
# Database
# ============================================================================
db-migrate: go-mod-tidy ## Run database migrations
	@echo "Running database migrations..."
	@go run ./cmd/server migrate up

db-migrate-down: go-mod-tidy ## Rollback database migrations
	@echo "Rolling back database migrations..."
	@go run ./cmd/server migrate down

db-seed: go-mod-tidy ## Seed the database
	@echo "Seeding database..."
	@go run ./cmd/server seed

# ============================================================================
# Dependencies
# ============================================================================
go-mod-tidy: ## Tidy Go modules
	@echo "Tidying Go modules..."
	@go mod tidy

go-mod-download: ## Download Go modules
	@echo "Downloading Go modules..."
	@go mod download

go-mod-verify: ## Verify Go modules
	@echo "Verifying Go modules..."
	@go mod verify

# ============================================================================
# Tools
# ============================================================================
tools: go-tools ## Install development tools

go-tools: ## Install Go development tools
	@echo "Installing Go development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install golang.org/x/vuln/cmd/govulncheck@latest
	@go install github.com/securego/gosec/v2/cmd/gosec@latest
	@go install github.com/aquasecurity/trivy@latest
	@echo "Tools installed"

# ============================================================================
# Documentation
# ============================================================================
docs: ## Generate documentation
	@echo "Generating documentation..."
	@mkdir -p docs
	@go doc -all ./... > docs/api.md
	@echo "Documentation generated in docs/"

docs-serve: ## Serve documentation locally
	@echo "Serving documentation..."
	@python3 -m http.server 8000 --directory docs 2>/dev/null || echo "Python not found"

# ============================================================================
# CI/CD
# ============================================================================
ci: test lint security ## Run CI pipeline locally

# ============================================================================
# Version
# ============================================================================
version: ## Show version information
	@echo "HostMK Security Framework"
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"
	@echo "Date: $(DATE)"

# Define the default target that runs when you just type 'make'
.DEFAULT_GOAL := help

.PHONY: help
help: ## Display this help message
	@echo "Usage:"
	@echo "  make <target>"
	@echo ""
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ { printf "  %-15s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

.PHONY: fmt-check
fmt-check: ## Check Go code formatting without modifying files (fails if unformatted)
	@echo "Checking code formatting..."
	@ERRORS=$$(gofmt -l .); \
	if [ -n "$$ERRORS" ]; then \
		echo "The following files are not formatted properly according to gofmt:"; \
		echo "$$ERRORS"; \
		echo "Run 'gofmt -w .' to fix them."; \
		exit 1; \
	fi

.PHONY: vet
vet: ## Run go vet static analysis tool
	go vet ./...

.PHONY: test
test: ## Run all unit tests
	go test -v ./...

.PHONY: bench
bench: ## Run all performance benchmarks
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

.PHONY: cover
cover: ## Run tests and display statement coverage percentage in terminal
	go test -cover ./...

.PHONY: cover-html
cover-html: ## Generate and open interactive HTML test coverage report
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

.PHONY: build
build: fmt-check vet ## Compile the packages
	go build ./...

.PHONY: clean
clean: ## Remove generated test coverage and build files
	rm -f coverage.out

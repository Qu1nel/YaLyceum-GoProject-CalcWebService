path		:= .

.DEFAULT_GOAL	:= help
MAKEFLAGS	+= --silent --no-print-directory

RED		:= \\033[0;31m
GREEN	:= \\033[1;32m
YELLOW	:= \\033[1;33m
BLUE	:= \\033[0;36m
RESET	:= \\033[0m

export GOFLAGS="-buildvcs=false"

OS := $(shell uname)

ifeq ($(OS),Linux)
    BINARY := bin/CalcService
	TEST_PATH := ./internal/server ./pkg/calculator
else ifeq ($(OS),Darwin)
    BINARY := bin/CalcService
	TEST_PATH := ./internal/server ./pkg/calculator
else
    BINARY := bin/CalcService.exe
	TEST_PATH := .\internal\server .\pkg\calculator
endif

ENTRY_POINT=./cmd/main.go

# Main
##############################################################################


.PHONY: run
run: install-dep  ## Launch app (main.go)
	go run $(ENTRY_POINT)

.PHONY: build-run
build-run:  ## Launch build app (int ./bin)
	./$(BINARY)

.PHONY: install-dep
install-dep:  ## Install tools for dev
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0
	go install -v github.com/go-critic/go-critic/cmd/gocritic@latest


# Build
##############################################################################


.PHONY: build
build: go-mod-tidy ## Build the app
	go build -o $(BINARY) $(ENTRY_POINT)


# Lint
##############################################################################


.PHONY: lint-check lint lints check
lint-check: go-mod-tidy gofmt goimports gocyclo golangci-lint gocritic ## Complete code base check with linters and formatters


# Alias
lint: lint-check  ## Alias for 'lint-check'
lints: lint-check  ## Alias for 'lint-check'
check: lint-check  ## Alias for 'lint-check'


.PHONY: gofmt
gofmt:  ## Use 'go fmt' utility as linter
	@echo -e
	@echo -e "$(BLUE)Applying go fmt..."
	@echo -e "$(GREEN)================$(RESET)"
	@echo -e
	go fmt ./...
	@echo -e

.PHONY: goimports
goimports:  ## Use 'goimports' utility as linter
	@echo -e
	@echo -e "$(BLUE)Applying goimports..."
	@echo -e "$(GREEN)================$(RESET)"
	@echo -e
	goimports -w .
	@echo -e

.PHONY: gocyclo
gocyclo:  ## Use 'gocyclo' utility as linter
	@echo -e
	@echo -e "$(BLUE)Applying gocyclo..."
	@echo -e "$(GREEN)===================$(RESET)"
	@echo -e
	gocyclo -over 15 .
	@echo -e

.PHONY: golangci-lint
golangci-lint:  ## Use 'golangci-lint' utility as formatter
	@echo -e
	@echo -e "$(BLUE)Applying golangci-lint..."
	@echo -e "$(GREEN)=================$(RESET)"
	@echo -e
	golangci-lint run ./...
	@echo -e

.PHONY: gocritic
gocritic:  ## Use 'gocritic' utility as formatter
	@echo -e
	@echo -e "$(BLUE)Applying gocritic..."
	@echo -e "$(GREEN)=================$(RESET)"
	@echo -e
	gocritic check ./...
	@echo -e

.PHONY: go-mod-tidy
go-mod-tidy:  ## Use 'go mod tidy' utility as formatter
	@echo -e
	@echo -e "$(BLUE)Applying go mod tidy..."
	@echo -e "$(GREEN)=================$(RESET)"
	@echo -e
	go mod tidy -v
	@echo -e


# Tests
##############################################################################


.PHONY: test
test: tests  ## Alias for 'tests'

.PHONY: tests
tests:  ## Unit tests for Go
	@echo -e
	@echo -e "$(BLUE)Applying go test..."
	@echo -e "$(GREEN)=================$(RESET)"
	@echo -e
	go test -tags=unit -timeout 30s -short -v $(TEST_PATH)
	@echo -e


# Clean cache
##############################################################################


.PHONY: clean
clean:  ## Clear cache
	go clean -cache


# Other
##############################################################################


.PHONY: help
help:  ## Show this output, i.e. help to use the commands
	grep -E '^[a-zA-Z_-]+:.*?# .*$$' Makefile | awk 'BEGIN {FS = ":.*?# "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


info-%:
	@make --dry-run --always-make $* | grep -v "info"


print-%:
	@$(info '$*'='$($*)')


.SILENT:

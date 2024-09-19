## —— Cache Warmer —————————————————————————————————————————
help: ## Outputs this help screen
	@grep -E '(^[a-zA-Z0-9_-]+:.*?##.*$$)|(^##)' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}{printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}' | sed -e 's/\[32m##/[33m/'

## —— Project ———————————————————————————————————————————————————————————————
VERSION := $(shell git describe --tags --exact-match 2>/dev/null || git rev-parse HEAD)
PROJECT_NAME := cache-warmer

run: ## Run the main go file on a Symfony project
	go run $(PROJECT_NAME).go $(path)

build: ## Build the vcw executable for the Linux/macOS
	${MAKE} lint
	CGO_ENABLED=0 go build -ldflags="-X main.version=$(VERSION) -s -w" -o $(PROJECT_NAME)
	strip $(PROJECT_NAME)
	shasum -a 256 $(PROJECT_NAME)

build-win: ## Build the vcw executable for Windows
	go build -ldflags="-X main.version=$(VERSION) -s -w" -o $(PROJECT_NAME).exe
	shasum -a 256 $(PROJECT_NAME).exe

clean: ## Clean all executables
	rm -f $(PROJECT_NAME) $(PROJECT_NAME).exe

deps: clean ## Clean dependencies
	go mod tidy -e
	go get -d -v ./...

update: ## Update dependencies
	go get -u ./...

## —— Tests ✅ —————————————————————————————————————————————————————————————————
test: ## Run all tests
	go test -count=1 -v ./... -coverprofile cover.out

## —— Coding standards ✨ ——————————————————————————————————————————————————————
lint: ## Run gofmt, simplify and lint
	gofmt -s -l -w .

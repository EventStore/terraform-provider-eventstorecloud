default: help

SHELL := bash

.PHONY: build
build:  ## Builds the app
	go build

.PHONY: generate
generate:  ## Generates the docs
	go generate

.PHONY: ci
ci: ## Performs the same checks as ci
	go build
	go generate
	git diff --exit-code  || (echo 'missing commits, or code was not generated' && exit 1)

.PHONY: install
install: ## Install the binary to the default target path
	go install

.PHONY: help
help: ## Display this information. Default target.
	@echo "Valid targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

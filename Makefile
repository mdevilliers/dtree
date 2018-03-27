
PACKAGES := $(shell go list ./... | grep -v /vendor/ )

GO_TEST = go test -covermode=atomic
GO_VET = go vet 
GO_COVER = go tool cover
GO_BENCH = go test -bench=.

all: vet test

.PHONY: all

get-build-deps: ## install build dependencies
	# for checking licences
	go get github.com/chespinoza/goliscan
	# various static analysis tools
	go get honnef.co/go/tools/cmd/megacheck

.PHONY: get-build-deps

check-vendor-licenses: ## check if dependencies licenses meet project requirements
	@goliscan check --direct-only -strict
	@goliscan check --indirect-only -strict

.PHONY: check-vendor-licenses

megacheck: ## run megacheck on the codebase
	megacheck $(PACKAGES)

.PHONY: megacheck

test: ## run tests
	$(GO_TEST) $(PACKAGES)

.PHONY: test

vet: ## run go vet
	$(GO_VET) $(PACKAGES)

.PHONY: vet

clean: ## clean up
	rm -rf tmp/

.PHONY: clean


# 'help' parses the Makefile and displays the help text
help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: help

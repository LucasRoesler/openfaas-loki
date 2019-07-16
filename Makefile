.NAME := openfaas-loki
.BIN_NAME := loki-provider
.PKG := github.com/LucasRoesler/$(.NAME)
.IMAGE_PREFIX=theaxer/$(.NAME)

.GIT_COMMIT=$(shell git rev-parse HEAD)
.GIT_VERSION=$(shell git describe --tags 2>/dev/null || echo "$(.GIT_COMMIT)")
.GIT_UNTRACKEDCHANGES := $(shell git status --porcelain --untracked-files=no)
ifneq ($(.GIT_UNTRACKEDCHANGES),)
	.GIT_COMMIT := $(.GIT_COMMIT)-dirty
endif

################################
################################
.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: env
env: ## Print debug information about your local environment
	@echo git: $(shell git version)
	@echo go: $(shell go version)
	@echo golint: $(shell which golint)
	@echo gofmt: $(shell which gofmt)
	@echo staticcheck: $(shell which staticcheck)


.PHONY: changelog
changelog: ## Print git hitstory based changelog
	@git --no-pager log --no-merges --pretty=format:"%h : %s (by %an)" $(shell git describe --tags --abbrev=0)...HEAD
	@echo ""

################################
################################
.PHONY: lint
lint: ## Verifies `golint` passes
	@echo "+ $@"
	@golint -set_exit_status $(shell go list ./pkg/...)

.PHONY: fmt
fmt: $(shell find ./pkg) ## Verifies all files have been `gofmt`ed
	@echo "+ $@"
	@gofmt -s -l . | tee /dev/stderr

.PHONY: .test-ci
.test-ci:
	@echo "+ test"
	GO111MODULE=on go test -cover ./pkg/...

.PHONY: test
test: $(shell find ./pkg) lint  ## Runs the go tests
	-@$(MAKE) .test-ci

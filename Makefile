.NAME := openfaas-loki
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
	@golint -set_exit_status $(shell go list ./pkg/... ./cmd/...)

.PHONY: fmt
fmt: $(shell find ./pkg ./cmd) ## Verifies all files have been `gofmt`ed
	@echo "+ $@"
	@gofmt -s -l . | tee /dev/stderr

.PHONY: test
test: $(shell find ./pkg ./cmd)  ## Runs the go tests
	@echo "+ test"
	GO111MODULE=on go test -race -cover ./...

.PHONY: install
install: $(shell find ./pkg ./cmd) ## Build the project and store the binaries in the GOPATH
	@echo "+ install $(.NAME): $(.GIT_COMMIT)"
	@CGO_ENABLED=0 go install \
		-v -ldflags \
		"-X ${.PKG}/pkg.GitCommit=${.GIT_COMMIT} -X ${.PKG}/pkg.Version=${.GIT_VERSION}" \
		.


.PHONY: build
build:  ## Create a docker image using the binary from make build
	DOCKER_BUILDKIT=1 docker build \
		-t ${.IMAGE_PREFIX}:latest \
		-t ${.IMAGE_PREFIX}:dev \
		--build-arg GIT_COMMIT=${.GIT_COMMIT} \
		--build-arg GIT_VERSION=${.GIT_VERSION} \
		-f ./Dockerfile .



.PHONY: package
package:  ## Package the helm chart
	@echo "+ helm package"
	helm package chart/openfaas-loki -d docs/
	helm repo index ./docs

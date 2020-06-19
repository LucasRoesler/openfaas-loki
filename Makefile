.NAME := openfaas-loki
.PKG := github.com/LucasRoesler/$(.NAME)
.IMAGE_PREFIX=theaxer/$(.NAME)
.TAG=latest

.GIT_COMMIT=$(shell git rev-parse HEAD)
.GIT_VERSION=$(shell git describe --tags 2>/dev/null || echo "$(.GIT_COMMIT)")
.GIT_UNTRACKEDCHANGES := $(shell git status --porcelain --untracked-files=no)
ifneq ($(.GIT_UNTRACKEDCHANGES),)
	.GIT_COMMIT := $(.GIT_COMMIT)-dirty
endif

ARCHS=amd64 arm64 armhf ppc64le
BUILD_ARGS=--build-arg GIT_COMMIT=$(.GIT_COMMIT) --build-arg GIT_VERSION=$(.GIT_VERSION)

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
build: $(addprefix build-,$(ARCHS))  ## Build Docker images for all architectures

.PHONY: build-%
build-%:
	DOCKER_BUILDKIT=1 docker build $(BUILD_ARGS) --build-arg go_opts="GOARCH=$*" \
		-t ${.IMAGE_PREFIX}:${.TAG}-$* \
		-f ./Dockerfile .

build-armhf:
	DOCKER_BUILDKIT=1 docker build $(BUILD_ARGS) --build-arg go_opts="GOARCH=arm GOARM=6" \
		-t ${.IMAGE_PREFIX}:${.TAG}-armhf \
		-f ./Dockerfile .

.PHONY: push
push: $(addprefix push-,$(ARCHS)) ## Push Docker images for all architectures

.PHONY: push-%
push-%:
	docker push ${.IMAGE_PREFIX}:${.TAG}-$*

.PHONY: manifest
manifest: ## Create and push Docker manifest to combine all architectures in multi-arch Docker image
	docker manifest create --amend ${.IMAGE_PREFIX}:${.TAG} $(addprefix ${.IMAGE_PREFIX}:${.TAG}-,$(ARCHS))
	$(MAKE) $(addprefix manifest-annotate-,$(ARCHS))
	docker manifest push -p ${.IMAGE_PREFIX}:${.TAG}

.PHONY: manifest-annotate-%
manifest-annotate-%:
	docker manifest annotate ${.IMAGE_PREFIX}:${.TAG} ${.IMAGE_PREFIX}:${.TAG}-$* --os linux --arch $*

.PHONY: manifest-annotate-armhf
manifest-annotate-armhf:
	docker manifest annotate ${.IMAGE_PREFIX}:${.TAG} ${.IMAGE_PREFIX}:${.TAG}-armhf --os linux --arch arm --variant v6

.PHONY: package
package:  ## Package the helm chart
	@echo "+ helm package"
	helm package chart/openfaas-loki -d docs/
	helm repo index ./docs

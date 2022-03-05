
include .bingo/Variables.mk

.PROJECT_ROOT=$(shell pwd)
OWNER=LucasRoesler
OWNER_LOWER:=$(shell echo "${OWNER}" | tr '[:upper:]' '[:lower:]')
NAME=openfaas-loki
REGISTRY=


GIT_COMMIT=$(shell git rev-parse HEAD)
GIT_VERSION=$(shell git describe --tags --always --dirty 2>/dev/null)
.GIT_UNTRACKEDCHANGES := $(shell git status --porcelain --untracked-files=no)
ifneq ($(.GIT_UNTRACKEDCHANGES),)
	GIT_VERSION := $(GIT_VERSION)-$(shell date +"%s")
endif

GOBINS=$(shell go env GOBIN)

default: install

.PHONY: version
version:
	@echo $(GIT_VERSION)

## Build env
setup: $(GORELEASER) $(GOLANGCI_LINT) ## install the required build environment binaries


################################
################################
.PHONY: lint
lint: $(GOLANGCI_LINT) ## Verifies `golangci-lint` passes
	@echo "+ $@"
	@$(GOLANGCI_LINT) run ./...

.PHONY: lint-fix
lint-fix: $(GOLANGCI_LINT) ## Verifies `golangci-lint` passes
	@echo "+ $@"
	@$(GOLANGCI_LINT) run --fix ./...

.PHONY: fmt
fmt:
	@echo "+ $@"
	@gofmt -s -l . | tee /dev/stderr

test:
	@go test -cover ./...


${GOBINS}/${NAME}:  $(shell find . -name '*.go') go.*
	@echo "+ $@"
	go install -tags 'osusergo netgo' -v -ldflags "\
		-X github.com/${OWNER}/${NAME}/pkg/cmd.GitCommit=$(GIT_COMMIT) \
		-X github.com/${OWNER}/${NAME}/pkg/cmd.Version=$(GIT_VERSION)" \
		.

install: ${GOBINS}/${NAME}

dist: $(GORELEASER) $(shell find . -name '*.go') go.*
	$(GORELEASER) build --skip-validate --rm-dist


ARCH:=linux/amd64

image: dist
	docker buildx build \
		--load \
		--platform ${ARCH} \
		--tag $(REGISTRY)$(OWNER_LOWER)/$(NAME):$(GIT_VERSION) \
		.
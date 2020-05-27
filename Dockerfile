# syntax=docker/dockerfile:experimental
FROM golang:1.14 as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

ENV GO111MODULE=on
ENV CGO_ENABLED=0
ARG GOOS=linux
ARG GIT_COMMIT=0
ARG GIT_VERSION=dev
ARG REPO=github.com/LucasRoesler/openfaas-loki

WORKDIR /app
COPY cmd ./cmd
COPY pkg ./pkg
COPY main.go .

RUN --mount=type=cache,target=/root/.cache/go-build --mount=type=cache,target=/go/pkg/mod go build -o /bin/openfaas-loki \
    -v -ldflags "\
    -X $REPO/pkg.GitCommit=$GIT_COMMIT -X $REPO/pkg.Version=$GIT_VERSION \
    -extldflags \"-static\"" \
    .

FROM alpine:3.11 as image

RUN apk --no-cache --update add ca-certificates

COPY --from=builder /bin/openfaas-loki /bin/openfaas-loki

ENTRYPOINT ["openfaas-loki"]

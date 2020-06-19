# syntax=docker/dockerfile:experimental
FROM golang:1.14-alpine as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

ENV GO111MODULE=on
ENV CGO_ENABLED=0
ARG go_opts
ARG GOOS=linux
ARG GIT_COMMIT=0
ARG GIT_VERSION=dev
ARG REPO=github.com/LucasRoesler/openfaas-loki

WORKDIR /app
COPY cmd ./cmd
COPY pkg ./pkg
COPY main.go .

RUN --mount=type=cache,target=/root/.cache/go-build --mount=type=cache,target=/go/pkg/mod \
    env $go_opts go build -o /bin/openfaas-loki \
    -v -ldflags "\
    -X $REPO/pkg.GitCommit=$GIT_COMMIT -X $REPO/pkg.Version=$GIT_VERSION \
    -extldflags \"-static\"" \
    .

# we can't add user in next stage because it's from scratch
# ca-certificates and tmp folder are also missing in scratch
# so we add all of it here and copy files in next stage
RUN addgroup -S app \
    && adduser -S -g app app \
    && mkdir /scratch-tmp

FROM scratch

EXPOSE 9191

ENV http_proxy      ""
ENV https_proxy     ""
USER app

COPY --from=builder /etc/passwd /etc/group /etc/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder --chown=app:app /scratch-tmp /tmp
COPY --from=builder /bin/openfaas-loki /bin/openfaas-loki

ENTRYPOINT ["openfaas-loki"]

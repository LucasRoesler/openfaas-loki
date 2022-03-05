FROM --platform=${TARGETPLATFORM:-linux/amd64} gcr.io/distroless/static:nonroot

ARG TARGETPLATFORM
ARG TARGETOS
ARG TARGETARCH

COPY dist/openfaas-loki_${TARGETOS}_${TARGETARCH}/openfaas-loki /

ENTRYPOINT ["/openfaas-loki"]
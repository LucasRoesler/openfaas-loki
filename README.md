# OpenFaaS Loki Provider

[![CircleCI](https://img.shields.io/circleci/build/github/LucasRoesler/openfaas-loki/master.svg)](https://circleci.com/gh/LucasRoesler/openfaas-loki) [![Go Report Card](https://goreportcard.com/badge/github.com/LucasRoesler/openfaas-loki)](https://goreportcard.com/report/github.com/LucasRoesler/openfaas-loki) [![Docker](https://img.shields.io/docker/pulls/theaxer/openfaas-loki.svg)](https://cloud.docker.com/repository/docker/theaxer/openfaas-loki)

A [Loki](https://github.com/grafana/loki) powered log provider for [OpenFaaS](https://www.openfaas.com)

## Install with Helm

```sh
helm upgrade --install ofloki ./deployment/openfaas-loki \
    --namespace openfaas \
    --set lokiURL=http://loki.monitoring:3100 \
    --set logLevel=DEBUG
```

Then update the gateway with the environment variable described in the NOTES output of the helm install.

Test the installation using

```sh
faas-cli store deploy nodeinfo
echo "" | faas-cli invoke nodeinfo
faas-cli logs nodeinfo --tail=3
```

## Development flow
OpenFaaS Loki is built with go 1.12+ and uses go modules

### Run unit tets

```sh
make test
```


### Local install

You can install using
```sh
make install
```

Check the installation using

```sh
openfaas-loki --version
```


### Build Docker
To enable efficient builds, the Dockerfile uses the experimental RUN syntax to support build-time caches. This requires enabling the "experimental" features in your Docker installation.  This will enable using [buildkit as the build engine](https://github.com/moby/buildkit/blob/master/frontend/dockerfile/docs/experimental.md#run---mounttypecache)

Build the Docker image and verify the build using

```sh
make build
docker run theaxer/openfaas-loki:dev --version
```

